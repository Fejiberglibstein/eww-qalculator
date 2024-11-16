package parser

import (
	"log"
	"strings"
)

type Token struct {
	Value string `json:"value"`
	Class Class  `json:"class"`
}

// Gets the tokens out from a qalc expression
func ParseTokens(input string) {
	split := strings.Split(input, "\x1B")

	tokens := make([]Token, 0)

	for _, tok := range split {
		class, offset, err := parseAnsiSeq(tok)
		if err != nil {
			log.Panic("Error parsing ansi seq: ", err)
		}
		tokens = append(tokens, Token{
			Value: tok[offset:],
			Class: class,
		})
	}

}
