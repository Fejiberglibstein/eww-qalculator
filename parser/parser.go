package parser

import (
	"errors"
	"strings"
)

type Class string

const (
	Number Class = "number"
)

type Token struct {
	Value string `json:"value"`
	Class Class  `json:"class"`
}

type AnsiSeq struct {

}

// Gets the tokens out from a qalc expression
func ParseTokens(input string) {

	split := strings.Split(input, "\x1B")

	for _, tok := range split {
		parseAnsiSeq(tok)
	}

}

func parseAnsiSeq(tok string) (Class, int, error) {
	// Skip the first token
	for tok[0] != '[' {
		if len(tok) > 1 {
			tok = tok[1:]
		} else {
			return "", 0, errors.New("Ansi seq not found in string segment")
		}
	}
	for i, char := range tok {

	}
}
