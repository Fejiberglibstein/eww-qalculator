package parser

import (
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
		color:    0,
	}

	if seq != expected {
		t.Error("seq not expected, got", seq)
	}

}
