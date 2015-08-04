package pig

import (
	"crypto/sha1"
	"fmt"
	"io"
	"strconv"
)

// StringSliceContains determines whether this slice contain this string?
func StringSliceContains(hay []string, needle string) bool {
	for _, straw := range hay {
		if straw == needle {
			return true
		}
	}
	return false
}

// StringSliceEquals determines whether two string slices are the same.
// Arbitrary regards nil and []string{} as NOT equal.
func StringSliceEquals(a []string, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	n := len(a)
	if n != len(b) {
		return false
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Last returns the last element of a slice,
// or empty string if the slice is empty.
func Last(array []string) (result string) {
	if n := len(array); n >= 1 {
		result = array[n-1]
	}
	return
}

// FilterOut creates a new slice from s, removing any element contained in forbidden.
// Accept nils and empty lists.
// Never yields nil, but may yield an empty list.
func FilterOut(s []string, forbidden []string) []string {
	p := make([]string, 0, len(s))
	for _, v := range s {
		if !StringSliceContains(forbidden, v) {
			p = append(p, v)
		}
	}
	return p
}

// FilterStrings creates a new slice from s, removing any element that doesn't match the predicate fn.
//
// From http://blog.golang.org/go-slices-usage-and-internals
//
// Accept nils and empty lists.
// Never yields nil, but may yield an empty list.
//
// Items that match the predicate are kept (not thrown away, do not confuse!)
func FilterStrings(s []string, fn func(string) bool) []string {
	p := make([]string, 0, len(s))
	for _, v := range s {
		if fn(v) {
			p = append(p, v)
		}
	}
	return p
}

// RemoveEmptyStrings creates a new slice from s, removing all empty elements.
func RemoveEmptyStrings(s []string) []string {
	p := make([]string, 0, len(s))
	for _, v := range s {
		if v != "" {
			p = append(p, v)
		}
	}
	return p
}

// MapStrings applies f to each element of str, and returns a new slice of the
// same size containing all the results.
func MapStrings(str []string, f func(string) string) []string {
	result := make([]string, len(str))
	for i, s := range str {
		result[i] = f(s)
	}
	return result
}

// Sha1hash is a shorthand to get a sha1 hash as a string.
func Sha1hash(s string) string {
	h := sha1.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// String2Int is a shorthand to convert a string to an int.
// It doesn't panic and does'nt return any error, but it may return -1.
// Thus it should be used only to parse positive integer that are extermely
// likely to be well-formed.
func String2Int(str string) int {
	if number64, err := strconv.ParseInt(str, 10, 32); err == nil {
		return int(number64)
	} else {
		return -1
	}
}

// Min returns the minimum of a and b.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of a and b.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
