package parser

import (
	"strings"
)

type Token struct {
	Value string `json:"value"`
	Class Class  `json:"class"`
}

type Line []Token

// Gets the tokens out from a qalc expression
func parseTokens(input string) Line {

	input = strings.TrimLeft(input, " \t")
	input = strings.ReplaceAll(input, "\n", "")
	// Give input a default ansi seq to begin with
	input = "[0;0m" + input
	split := strings.Split(input, "\x1B")

	tokens := make([]Token, 0)

	for _, tok := range split {
		seq, offset, err := parseAnsiSeq(tok)
		if err != nil {
			// ignore any tokens that produce errors
			continue
		}
		tok = tok[offset:]

		split := splitEquals(tok)
		for _, token := range split {
			tokens = append(tokens, Token{
				Value: token,
				Class: seq.getClass(),
			})
		}

	}

	return tokens
}

// Split a token at either `=` or `≈` so that we can have the equals as its own
// token
//
// # Example
//
// splitEquals("foo = bar ≈ 10") -> ["foo ", "=", " bar ", "≈", " 10"]
// splitEquals("foobar") -> ["foobar"]
func splitEquals(input string) []string {
	res := make([]string, 0)
	var acc strings.Builder

	for _, rune := range input {
		if rune == '≈' || rune == '=' {
			res = append(res, acc.String())
			res = append(res, string(rune))
			acc = strings.Builder{}
		} else {
			acc.WriteRune(rune)
		}
	}

	res = append(res, acc.String())
	return res
}

func ParseLines(lines []string) []Line {
	res := make([]Line, 0)

	for _, line := range lines {
		res = append(res, parseTokens(line))
	}

	return res
}
