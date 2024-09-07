#include "./vector.h"
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>

int open_qalc();
int qalc_send(char *);

typedef enum : uint8_t {
    RESET = 0,
    BOLD,
    DIM,
    ITALIC,
    UNDERLINE,
} AnsiGraphics;

typedef enum : uint8_t {
    COLOR_RESET = 0,
    DEFAULT = 39,
    BLACK = 30,
    RED = 31,
    GREEN = 32,
    YELLOW = 33,
    BLUE = 34,
    MAGENTA = 35,
    CYAN = 36,
    WHITE = 37,
} AnsiColor;

typedef struct AnsiSeq {
    // The reset sequence of the ansi sequence
    AnsiGraphics reset_sequence;
    // The color of the ansi sequence
    AnsiColor color;
    // The part of the ansi sequence that has currently been reached. Each
    // 'part' is delimited by semicolons. For example, the ansi sequence
    // '\x[1;31m' has 2 parts, '1' and '31'.
    int part;
} AnsiSeq;

typedef struct {
    AnsiSeq ansi;
    char *value;
} Token;

typedef Vector(Token) Line;
typedef struct {
    Vector(Line) lines;
    bool valid;
} Equation;

#define NUMBER ((RESET << 8) + CYAN)
#define EXPRESSION ((RESET << 8) + COLOR_RESET)
#define VARIABLE ((ITALIC << 8) + YELLOW)

#define QALC_HIST "/home/austint/.local/state/qalc_hist"
#define EXPR_LEN 4096

int main(int argc, char **argv) {
    if (argc < 2) {
        fprintf(stderr, "No arguments\n");
        return -1;
    }
    if (strcmp(argv[1], "open") == 0) {
        return open_qalc();
    } else if (strcmp(argv[1], "expr") == 0) {
        return qalc_send(argv[2]);
    }
    return 1;
}

char *get_token_class(AnsiSeq seq) {
    int color = (seq.reset_sequence << 8) + seq.color;
    char *ret;

    switch (color) {
    case NUMBER:
        asprintf(&ret, "number");
        break;
    case EXPRESSION:
        asprintf(&ret, "expression");
        break;
    case VARIABLE:
        asprintf(&ret, "variable");
        break;
    default:
        asprintf(&ret, "%d;%d", seq.reset_sequence, seq.color);
        break;
    }
    return ret;
}

void print_line(Line tokens) {
    printf("[");

    for (int i = 0; i < tokens.length; i++) {
        Token token = tokens.items[i];
        char *token_class = get_token_class(token.ansi);
        printf("{\"class\":\"%s\",\"value\":\"%s\"}", token_class,
               token.value);
        free(token_class);
		if (i < tokens.length - 1) {
			putchar(',');
		}
    }
    printf("]");
}

void print_equation(Equation equation) {
    printf("[");
    for (int i = 0; i < equation.lines.length; i++) {
        print_line(equation.lines.items[i]);
		if (i < equation.lines.length - 1) {
			putchar(',');
		}
    }
    printf("]");
}

// Parse a character into an AnsiSeq.
//
// This function will be called for each character in a token until 'm' is
// reached, indicating the end of the sequence.
//
// This function will return `1` when the sequence is over, and `2` if there is
// no sequence. In all other cases it will return `0`;
int parse_ansi_seq(AnsiSeq *seq, char c) {
    static AnsiColor current_color = WHITE;

    switch (c) {
    case '[':
        return 0;
    case ';':
        seq->part++;
        return 0;
    case 'm':
        if (seq->color == 0) {
            seq->color = current_color;
        }
        return 1;
    }

    if (c <= '9' && c >= '0') {
        uint8_t p = ((uint8_t *)seq)[seq->part];
        p *= 10;
        p += c - '0';
        ((uint8_t *)seq)[seq->part] = p;
        // If we're currently parsing the color
        if (seq->part == 1) {
            current_color = p;
        }

        return 0;
    }

    // If nothing else matches, there is no ansi sequence.
    return 2;
}

Token parse_token(char *tok) {
    AnsiSeq seq = {0};
    int i;
    for (i = 0; tok[i] != '\0'; i++) {
        int parse = parse_ansi_seq(&seq, tok[i]);
        if (parse == 1) {
            i++;
            break;
        } else if (parse == 2) {
            break;
        }
    }
    char *val = malloc(strlen(&(tok[i])));
	strcpy(val, &(tok[i]));

    return (Token){.value = val, .ansi = seq};
}

void parse_line(char *line, Equation *equation) {
    // If we reach this, then the equation is done
    if (strncmp(line, "> \n", 4) == 0) {
        if (equation->valid) {
            print_equation(*equation);
            putchar('\n');
			fflush(stdout);
        }

        // Clean up the equation, freeing all the tokens in the
        // equation
        for (int i = 0; i < equation->lines.length; i++) {
			Line tokens = equation->lines.items[i];
			for (int j = 0; j < tokens.length; j ++) {
				free(tokens.items[j].value);
			}
            vec_free(&tokens);
        }
        // Set length to 0 so we can override it later.
        equation->lines.length = 0;
        equation->valid = true;
    }

    if (line[0] == '>') {
        return;
    }

    // Remove the indentation at the beginning of a line
    while (line[0] == ' ') {
        line = &(line[1]);
    }

    char *tok = strtok(line, "\e\n");
    Line tokens = {0};

    while (tok != NULL) {
        Token token = parse_token(tok);

        // If the line has a warning, we should ignore this token.
        if (strcmp(token.value, "warning: ") == 0) {
            equation->valid = false;
        }
        if (strlen(token.value) != 0) {
            vec_append(&tokens, token);
        }

        tok = strtok(NULL, "\e\n");
    }
    vec_append(&(equation->lines), tokens);
}

int open_qalc() {
    char result_buf[EXPR_LEN];
    // Make qalc history file
    FILE *qalc = fopen(QALC_HIST, "w");
    if (qalc == NULL) {
        fprintf(stderr, "Could not make qalc history file\n");
        return 1;
    }
    fclose(qalc);

    // Read the last line of qalc_hist
    FILE *tail = popen("tail -F " QALC_HIST " | qalc", "r");
    if (tail == NULL) {
        fprintf(stderr, "Could not tail qalc history\n");
        return 1;
    }

    char *lines;

    Equation equation = {.lines = {0}, .valid = true};
    for (;;) {
        fgets(result_buf, sizeof(result_buf), tail);

        parse_line(result_buf, &equation);
    }
    return 0;
}

int qalc_send(char *expr) {
    FILE *qalc = fopen(QALC_HIST, "a");
    if (qalc == NULL) {
        fprintf(stderr, "Could not open qalc history file\n");
        return 1;
    }

    fprintf(qalc, "%s\n\n", expr);
    fclose(qalc);

    return 0;
}
