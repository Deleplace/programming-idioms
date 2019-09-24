package main

import (
	"testing"
)

var emphasizeTests = []struct {
	text     string
	expected string
}{
	{"", ""},

	// Match
	{"_x", "<b><i>x</i></b>"},
	{"__id", "<b><i>_id</i></b>"},
	{"_$php", "<b><i>$php</i></b>"},
	{"Variables _a, _b, _c are declared.", "Variables <b><i>a</i></b>, <b><i>b</i></b>, <b><i>c</i></b> are declared."},

	// Don't match
	{"x", "x"},
	{"a_b", "a_b"},
	{"ab_", "ab_"},
	{"_", "_"},
}

func TestEmphasize(t *testing.T) {
	for i, tt := range emphasizeTests {
		if processed := emphasize(tt.text); processed != tt.expected {
			t.Errorf("%d. emphasize(%v) => %v, want %v", i, tt.text, processed, tt.expected)
		}
	}
}
