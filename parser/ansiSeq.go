package parser

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"unicode"
)

type ansiSeq struct {
	graphics ansiGraphic
	color    ansiColor
}

type ansiColor uint8
type ansiGraphic uint8
type Class string

// all ansi graphics
const (
	reset     ansiGraphic = 0
	bold      ansiGraphic = 1
	dim       ansiGraphic = 2
	italic    ansiGraphic = 3
	underline ansiGraphic = 4
)

// All ansi colors
const (
	colorReset ansiColor = 0
	black      ansiColor = 30
	red        ansiColor = 31
	green      ansiColor = 32
	yellow     ansiColor = 33
	blue       ansiColor = 34
	magenta    ansiColor = 35
	cyan       ansiColor = 36
	white      ansiColor = 37
)

// All qalc ansi sequence class names
const (
	cNumber     Class = "number"
	cExpression Class = "expression"
	cVariable   Class = "variable"
	cUnit       Class = "unit"
	cError      Class = "error"
	cErrorMsg   Class = "errorMsg"
	cBoolean    Class = "boolean"
)

func (seq *ansiSeq) getClass() Class {
	var res Class
	switch *seq {
	case ansiSeq{graphics: reset, color: cyan}:
		res = cNumber
	case ansiSeq{graphics: reset, color: colorReset}:
		res = cExpression
	case ansiSeq{graphics: italic, color: yellow}:
		res = cVariable
	case ansiSeq{graphics: reset, color: green}:
		res = cUnit
	case ansiSeq{graphics: reset, color: red}:
		res = cError
	case ansiSeq{graphics: italic, color: colorReset}:
		res = cErrorMsg
	case ansiSeq{graphics: reset, color: yellow}:
		res = cBoolean
	}
	return res
}

// parse an ansi sequence out from a string.
//
// Will return the class, offset of the string to skip over the ansi sequence,
// or any error
func parseAnsiSeq(tok string) (ansiSeq, int, error) {
	// Skip the first token
	for tok[0] != '[' {
		if len(tok) > 1 {
			tok = tok[1:]
		} else {
			return ansiSeq{}, 0, errors.New("Ansi seq not found in string segment")
		}
	}
	part := 0
	var num strings.Builder
	var seq ansiSeq
	var seqLength int

	for i, char := range tok {
		switch char {
		case '[':
			// ignore [
			continue
		case 'm':
			// 'm' is the ending character in an ansi seq
			if err := seq.addPart(num.String(), part); err != nil {
				return ansiSeq{}, 0, err
			}
			seqLength = i + 1
			break
		case ';':
			if err := seq.addPart(num.String(), part); err != nil {
				return ansiSeq{}, 0, err
			}
			// Make a new string for num
			num = strings.Builder{}
			part += 1
			continue
		}
		if unicode.IsDigit(char) {
			num.WriteRune(char)
		}
	}
	return seq, seqLength, nil
}

func (seq *ansiSeq) addPart(num string, part int) error {

	digit, err := strconv.Atoi(num)
	if err != nil {
		log.Print(err)
		return err
	}

	if part == 0 {
		seq.graphics = ansiGraphic(digit)
	} else {
		seq.color = ansiColor(digit)
	}

	return nil
}
