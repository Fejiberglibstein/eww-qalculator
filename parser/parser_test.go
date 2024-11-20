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

func TestEqual(t *testing.T) {
	res := splitEquals("foo = bar = le")
	if slices.Compare(res, []string{"foo ", "=", " bar ", "=", " le"}) != 0 {
		t.Error("Not equal, got ", strings.Join(res, ","))
	}

	res = splitEquals("foobar")
	if slices.Compare(res, []string{"foobar"}) != 0 {
		t.Error("Not equal, got ", strings.Join(res, ","))
	}
}
