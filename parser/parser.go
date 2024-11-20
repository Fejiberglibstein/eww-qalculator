package parser

import (
	"strings"
)

type Token struct {
	Value string `json:"value"`
	Class Class  `json:"class"`
}

type Line []Token

func ParseLines(lines []string) []Line {
	res := make([]Line, 0)

	for _, line := range lines {
		res = append(res, parseTokens(line))
	}

	return res
}

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

// Represents the result of a qalc calculation:
//
// 10/3 + 4 = 22/3 ≈ 7.33333 ->
//
//	Result {
//	   Approximate: 7.33333
//	   Actual: 22/3
//	}
type Result struct {
	// The actual result, always included.
	Actual []Token
	// The approximate result, used when the expression has a ≈ in it.
	//
	// If there is no approximate result, then this will be empty
	Approximate []Token
}

func resultAcc(addTo string, result *Result, acc []Token) {
	if addTo == "approximate" {
		result.Approximate = acc
	}
	if addTo == "actual" {
		result.Actual = acc
	}
}

func GetResults(lines []Line) []Result {
	results := make([]Result, 0)
	for _, tokens := range lines {

		var result Result
		addTo := ""
		acc := make([]Token, 0)

		for _, token := range tokens {
			switch token.Value {
			case "=":
				resultAcc(addTo, &result, acc)
				if len(result.Actual) == 0 {
					addTo = "actual"
				}
				acc = make([]Token, 0)
			case "≈":
				resultAcc(addTo, &result, acc)
				addTo = "approximate"
				acc = make([]Token, 0)
			default:
				acc = append(acc, token)
			}
		}
		resultAcc(addTo, &result, acc)
		results = append(results, result)
	}

	return results
}
