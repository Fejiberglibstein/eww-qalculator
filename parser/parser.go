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
	split := strings.Split(input, "\x1B")

	tokens := make([]Token, 0)

	for _, tok := range split {
		seq, offset, err := parseAnsiSeq(tok)
		if err != nil {
			// ignore any tokens that produce errors
			continue
		}
		tok = tok[offset:]

		tokens = append(tokens, Token{
			Value: tok,
			Class: seq.getClass(),
		})
	}

	return tokens
}

func ParseLines(lines []string) []Line {
	res := make([]Line, 0)

	for _, line := range lines {
		res = append(res, parseTokens(line))
	}

	return res
}
