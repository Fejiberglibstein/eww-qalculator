#include "./vector.h"
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>

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

// List of all colors qalc uses and what they represent
#define NUMBER ((RESET << 8) + CYAN)
#define EXPRESSION ((RESET << 8) + COLOR_RESET)
#define VARIABLE ((ITALIC << 8) + YELLOW)
#define UNIT ((RESET << 8) + GREEN)
#define ERROR ((RESET << 8) + RED)
#define ERROR_MSG ((ITALIC << 8) + COLOR_RESET)
#define BOOLEAN ((RESET << 8) + YELLOW)

#define QALC_HIST "/home/austint/.local/state/qalc_hist"
#define EXPR_LEN 4096

int qalc_send(char *expr);
int open_qalc();
void print_line(Line tokens);
void print_equation(Equation equation);
int parse_ansi_seq(AnsiSeq *seq, char c);
Token parse_token(char *tok);
void parse_line(char *line, Equation *equation);

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
    FILE *tail = popen("tail -F " QALC_HIST " | qalc --terse", "r");
    if (tail == NULL) {
        fprintf(stderr, "Could not tail qalc history\n");
        return 1;
    }

    char *lines;

    Equation equation = {.lines = {0}, .valid = true};
    for (;;) {
        fgets(result_buf, sizeof(result_buf), tail);

        /*for (int i = 0; result_buf[i] != '\0'; i++) {*/
        /*    printf("(%c)", result_buf[i]);*/
        /*    fflush(NULL);*/
        /*}*/

        parse_line(result_buf, &equation);
    }
    return 0;
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
    case UNIT:
        asprintf(&ret, "unit");
        break;
    case ERROR:
        asprintf(&ret, "error");
        break;
    case ERROR_MSG:
        asprintf(&ret, "error-msg");
        break;
    case BOOLEAN:
        asprintf(&ret, "boolean");
        break;
    default:
        asprintf(&ret, "unknown %d;%d", seq.reset_sequence, seq.color);
        break;
    }
    return ret;
}

void print_line(Line tokens) {
    printf("[");

    for (int i = 0; i < tokens.length; i++) {
        Token token = tokens.items[i];
        char *token_class = get_token_class(token.ansi);
        printf("{\"class\":\"%s\",\"value\":\"%s\"}", token_class, token.value);
        free(token_class);
        if (i < tokens.length - 1) {
            putchar(',');
        }
    }
    printf("]");
}

void print_equation(Equation equation) {
    if (equation.lines.length <= 2) {
        printf("[]");
        return;
    }

    printf("[");
    // We want to skip the first and last lines since they're empty
    for (int i = 1; i < equation.lines.length - 1; i++) {
        print_line(equation.lines.items[i]);
        if (i < equation.lines.length - 2) {
            putchar(',');
        }
    }
    printf("]");
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
            for (int j = 0; j < tokens.length; j++) {
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

    // If the line has indentation, replace the indentation with a empty ansi
    // code. This fixes some issues if the first character is part of an ansi
    // sequence
    for (int i = 0; line[i] == ' '; i ++) {
        if (i == 0) {
            line[i] = '[';
        } else {
            line[i] = 'm';
        }
    }

    char *tok = strtok(line, "\e\n");
    Line tokens = {0};

    while (tok != NULL) {
        Token token = parse_token(tok);

        if (token.value != NULL) {
            vec_append(&tokens, token);
            // If the line has a warning, we should ignore this token.
            if (strcmp(token.value, "warning: ") == 0) {
                /*equation->valid = false;*/
            }
        }

        tok = strtok(NULL, "\e\n");
    }
    vec_append(&(equation->lines), tokens);
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

    char *val = NULL;
    int len = strlen(&(tok[i]));
    if (len != 0) {
        val = calloc(len + 1, sizeof(char));
        strcpy(val, &(tok[i]));
    }

    // Replace all double quotes with single quotes
    for (int i = 0; i < len; i++) {
        if (val[i] == '"') {
            val[i] = '\'';
        }
    }
    return (Token){.value = val, .ansi = seq};
}

// Parse a character into an AnsiSeq.
//
// This function will be called for each character in a token until 'm' is
// reached, indicating the end of the sequence.
//
// This function will return `1` when the sequence is over, and `2` if there is
// no sequence. In all other cases it will return `0`;
int parse_ansi_seq(AnsiSeq *seq, char c) {
    static AnsiColor current_color = COLOR_RESET;

    switch (c) {
    case '[':
        return 0;
    case ';':
        if (seq->reset_sequence == RESET && seq->part == 0) {
            current_color = COLOR_RESET;
        }
        seq->part++;
        return 0;
    case 'm':
        if (seq->reset_sequence == RESET && seq->part == 0) {
            current_color = COLOR_RESET;
        }
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
