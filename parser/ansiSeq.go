package parser

import (
	"errors"
	"strconv"
	"unicode"
)

type ansiSeq struct {
	color    ansiColor
	graphics ansiGraphic
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
	cNumber     Class = "Number"
	cExpression Class = "Expression"
	cVariable   Class = "Variable"
	cUnit       Class = "Unit"
	cError      Class = "Error"
	cErrorMsg   Class = "ErrorMsg"
	cBoolean    Class = "Boolean"
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
func parseAnsiSeq(tok string) (Class, int, error) {
	// Skip the first token
	for tok[0] != '[' {
		if len(tok) > 1 {
			tok = tok[1:]
		} else {
			return "", 0, errors.New("Ansi seq not found in string segment")
		}
	}
	part := 0
	num := make([]rune, 0)
	var seq ansiSeq
	var seqLength int

	for i, char := range tok {
		switch char {
		case '[':
			// ignore [
			continue
		case 'm':
			// 'm' is the ending character in an ansi seq
			seqLength = i
			break
		case ';':
			seq.addPart(string(num), part)
			part += 1
			continue
		}
		if unicode.IsDigit(char) {
			num = append(num, char)
		}
	}
	return seq.getClass(), seqLength, nil
}

func (seq *ansiSeq) addPart(num string, part int) {

	digit, _ := strconv.Atoi(num)

	if part == 0 {
		seq.graphics = ansiGraphic(digit)
	} else {
		seq.color = ansiColor(digit)
	}
}
