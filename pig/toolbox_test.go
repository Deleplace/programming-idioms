package pig

import (
	"testing"
)

//
// Inspired by http://code.google.com/p/go-wiki/wiki/TableDrivenTests
//

// ---

var stringSliceContainsTests = []struct {
	inHay    []string
	inNeedle string
	out      bool
}{
	// Empty cases
	{nil, "foobar", false},
	{[]string{}, "foobar", false},
	// Negative cases
	{[]string{"other"}, "foobar", false},
	{[]string{"other", "other2"}, "foobar", false},
	// Positive cases
	{[]string{"foobar", "other2"}, "foobar", true},
	{[]string{"other2", "foobar"}, "foobar", true},
}

func TestStringSliceContains(t *testing.T) {
	for i, tt := range stringSliceContainsTests {
		if found := StringSliceContains(tt.inHay, tt.inNeedle); found != tt.out {
			t.Errorf("%d. StringSliceContains(%v, %v) => %v, want %v", i, tt.inHay, tt.inNeedle, found, tt.out)
		}
	}
}

// ---

var stringSliceEqualsTests = []struct {
	inA []string
	inB []string
	out bool
}{
	// Empty cases
	{nil, nil, true},
	{[]string{}, []string{}, true},
	// Arbitrary regards nil and []string{} as NOT equal.
	{nil, []string{}, false},
	{[]string{}, nil, false},
	// Positive cases
	{[]string{""}, []string{""}, true},
	{[]string{" "}, []string{" "}, true},
	{[]string{"a"}, []string{"a"}, true},
	{[]string{"a", "b"}, []string{"a", "b"}, true},
	{[]string{"a", "a"}, []string{"a", "a"}, true},
	// Negative cases
	{[]string{""}, nil, false},
	{[]string{" "}, nil, false},
	{[]string{"a"}, nil, false},
	{[]string{""}, []string{}, false},
	{[]string{" "}, []string{}, false},
	{[]string{"a"}, []string{}, false},
	{[]string{""}, []string{" "}, false},
	{[]string{" "}, []string{""}, false},
	{[]string{"a"}, []string{""}, false},
}

func TestStringSliceEquals(t *testing.T) {
	for i, tt := range stringSliceEqualsTests {
		out := StringSliceEquals(tt.inA, tt.inB)
		if out != tt.out {
			t.Errorf("%d. StringSliceEqual(%v, %v) => %v, want %v", i, tt.inA, tt.inB, out, tt.out)
		}
	}
}

// ---

var lastTests = []struct {
	in  []string
	out string
}{
	// Empty cases
	{nil, ""},
	{[]string{}, ""},
	// Normal cases
	{[]string{"aaa"}, "aaa"},
	{[]string{"aaa", "bbb"}, "bbb"},
}

func TestLast(t *testing.T) {
	for i, tt := range lastTests {
		if y := Last(tt.in); y != tt.out {
			t.Errorf("%d. Last(%v) => %v, want %v", i, tt.in, y, tt.out)
		}
	}
}

// ---

var filterOutTests = []struct {
	inS         []string
	inForbidden []string
	out         []string
}{
	// Empty cases
	{nil, nil, []string{}},
	{[]string{"aaa", "bbb"}, nil, []string{"aaa", "bbb"}},
	{nil, []string{"aaa", "bbb"}, []string{}},
	// Normal cases
	{[]string{"aaa", "bbb"}, []string{"ccc", "ddd"}, []string{"aaa", "bbb"}},
	{[]string{"aaa", "bbb"}, []string{"aaa", "ddd"}, []string{"bbb"}},
	{[]string{"aaa", "bbb"}, []string{"ccc", "aaa"}, []string{"bbb"}},
	{[]string{"aaa", "bbb"}, []string{"aaa", "bbb"}, []string{}},
	{[]string{"aaa", "bbb"}, []string{"bbb", "aaa"}, []string{}},
	{[]string{"aaa", "bbb", "ccc"}, []string{"ccc", "aaa"}, []string{"bbb"}},
}

func TestFilterOut(t *testing.T) {
	for i, tt := range filterOutTests {
		if filtered := FilterOut(tt.inS, tt.inForbidden); !StringSliceEquals(filtered, tt.out) {
			t.Errorf("%d. StringSliceContains(%v, %v) => %v, want %v", i, tt.inS, tt.inForbidden, filtered, tt.out)
		}
	}
}

// ---

var filterStringsTests = []struct {
	in  []string
	inF func(string) bool
	out []string
}{
	// Empty cases
	{nil, func(x string) bool { return len(x) >= 1 && x[0] == 't' }, []string{}},
	{[]string{}, func(x string) bool { return len(x) >= 1 && x[0] == 't' }, []string{}},
	// Normal cases
	{[]string{"banana", "toy", "monster", "", "truck"}, func(x string) bool { return len(x) >= 1 && x[0] == 't' }, []string{"toy", "truck"}},
	{[]string{"banana", "monster", ""}, func(x string) bool { return len(x) >= 1 && x[0] == 't' }, []string{}},
}

func TestFilterStrings(t *testing.T) {
	for i, tt := range filterStringsTests {
		filtered := FilterStrings(tt.in, tt.inF)
		if !StringSliceEquals(filtered, tt.out) {
			t.Errorf("%d. FilterStrings(%v, %v) => %v, want %v", i, tt.in, tt.inF, filtered, tt.out)
		}
	}
}

// ---

var mapStringsTests = []struct {
	in  []string
	inF func(string) string
	out []string
}{
	// Empty cases
	{nil, func(x string) string { return "**" + x + "**" }, []string{}},
	{[]string{}, func(x string) string { return "**" + x + "**" }, []string{}},
	// Normal cases
	{[]string{"banana"}, func(x string) string { return "**" + x + "**" }, []string{"**banana**"}},
	{[]string{"banana", "toy"}, func(x string) string { return "**" + x + "**" }, []string{"**banana**", "**toy**"}},
}

func TestMapStrings(t *testing.T) {
	for i, tt := range mapStringsTests {
		if transformed := MapStrings(tt.in, tt.inF); !StringSliceEquals(transformed, tt.out) {
			t.Errorf("%d. MapStrings(%v, %v) => %v, want %v", i, tt.in, tt.inF, transformed, tt.out)
		}
	}
}

// ---

var sha1hashTests = []struct {
	in  string
	out string
}{
	// Empty case
	{"", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
	// Normal cases
	{"aaa", "7e240de74fb1ed08fa08d38063f6a6a91462a815"},
}

func TestSha1hash(t *testing.T) {
	for i, tt := range sha1hashTests {
		if y := Sha1hash(tt.in); y != tt.out {
			t.Errorf("%d. Sha1hash(%v) => %v, want %v", i, tt.in, y, tt.out)
		}
	}
}

// ---
