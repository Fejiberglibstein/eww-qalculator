package parser

import (
	"slices"
	"strings"
	"testing"
)

func TestAnsiSeq(t *testing.T) {
	inp := "[10;6mfoo"

	seq, offset, err := parseAnsiSeq(inp)
	if err != nil {
		t.Error("Got error ", err)
	}

	if offset != 6 {
		t.Error("expected offset 6, got ", offset)
	}

	expected := ansiSeq{
		graphics: 10,
		color:    6,
	}

	if seq != expected {
		t.Error("seq not expected, got", seq)
	}

	inp = "[10mfoo"

	seq, offset, err = parseAnsiSeq(inp)
	if err != nil {
		t.Error("Got error ", err)
	}

	if offset != 4 {
		t.Error("expected offset 6, got ", offset)
	}

	expected = ansiSeq{
		graphics: 10,
		color:    6,
	}

	if seq != expected {
		t.Error("seq not expected, got", seq)
	}

}

func TestAnsi2m(t *testing.T) {
	inp := "[3mEmpty expression"
	lastColor = colorReset // stupid fix
	seq, offset, err := parseAnsiSeq(inp)
	if err != nil {
		t.Error("Got error ", err)
	}

	if offset != 3 {
		t.Error("expected offset 3, got ", offset)
	}

	expected := ansiSeq{
		graphics: 3,
		color:    0,
	}

	if seq != expected {
		t.Error("seq not expected, got", seq)
	}
}

func TestSplitEquals(t *testing.T) {
	res := splitEquals("foo = bar = le")
	if slices.Compare(res, []string{"foo ", "=", " bar ", "=", " le"}) != 0 {
		t.Error("Not equal, got ", strings.Join(res, ","))
	}

	res = splitEquals("foobar")
	if slices.Compare(res, []string{"foobar"}) != 0 {
		t.Error("Not equal, got ", strings.Join(res, ","))
	}
}

func newToken(value string) Token {
	return Token{
		Class: "expression",
		Value: value,
	}
}

func compareResult(r1, r2 Result) bool {
	if len(r1.Actual) != len(r2.Actual) || len(r1.Approximate) != len(r2.Approximate) {
		return false
	}

	for i, _ := range r1.Actual {
		if r1.Actual[i] != r2.Actual[i] {
			return false
		}
	}

	for i, _ := range r1.Approximate {
		if r1.Approximate[i] != r2.Approximate[i] {
			return false
		}
	}
	return true
}

func TestParseEqual(t *testing.T) {
	_, results := EvaluateEquation([]Line{
		[]Token{
			newToken("(10 - 3) "),
			newToken("="),

			newToken(" 7 "),
			newToken("347"),

			newToken("≈"),

			newToken("res2"),
		},
	})

	if !(len(results) != 0 && compareResult(results[0], Result{
		Actual: []Token{
			newToken(" 7 "),
			newToken("347"),
		},
		Approximate: []Token{
			newToken("res2"),
		},
	})) {
		t.Error("Not equal, got ", results[0])
	}
}
