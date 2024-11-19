package parser

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var lastColor ansiColor

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
	cUnknown    Class = "unknown"
)

func (seq *ansiSeq) getClass() Class {
	switch *seq {
	case ansiSeq{graphics: reset, color: cyan}:
		return cNumber
	case ansiSeq{graphics: reset, color: colorReset}:
		return cExpression
	case ansiSeq{graphics: italic, color: yellow}:
		return cVariable
	case ansiSeq{graphics: reset, color: green}:
		return cUnit
	case ansiSeq{graphics: reset, color: red}:
		return cError
	case ansiSeq{graphics: italic, color: colorReset}:
		return cErrorMsg
	case ansiSeq{graphics: reset, color: yellow}:
		return cBoolean
	default:
		return cUnknown
	}
}

// parse an ansi sequence out from a string.
//
// Will return the class, offset of the string to skip over the ansi sequence,
// or any error
func parseAnsiSeq(tok string) (ansiSeq, int, error) {
	part := 0
	var num strings.Builder
	var seq ansiSeq

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
			seqLength := i + 1
			return seq, seqLength, nil
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
	return ansiSeq{}, 0, errors.New("Could not find end of ansi seq")
}

func (seq *ansiSeq) addPart(num string, part int) error {

	digit, err := strconv.Atoi(num)
	if err != nil {
		return err
	}

	if part == 0 {
		seq.graphics = ansiGraphic(digit)
		if seq.graphics != reset {
			seq.color = lastColor
		}
	} else {
		seq.color = ansiColor(digit)
		lastColor = seq.color
	}

	return nil
}
