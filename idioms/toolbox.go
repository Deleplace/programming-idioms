package idioms

import (
	"crypto/sha1"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

// StringSliceContains determines whether this slice contains this string
func StringSliceContains(hay []string, needle string) bool {
	for _, straw := range hay {
		if straw == needle {
			return true
		}
	}
	return false
}

// StringSliceContainsCaseInsensitive determines whether this slice contains this string, regardless the case
func StringSliceContainsCaseInsensitive(hay []string, needle string) bool {
	needle = strings.ToLower(needle)
	for _, straw := range hay {
		straw = strings.ToLower(straw)
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
// From https://blog.golang.org/go-slices-usage-and-internals
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
	}
	return -1
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

// Shorten keeps only the n first bytes, add appends "..." if needed.
func Shorten(s string, n int) string {
	// TODO better handle UTF-8
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}

// Flatten replaces newlines by spaces
func Flatten(s string) string {
	s = strings.Replace(s, "\r\n", " ", -1)
	s = strings.Replace(s, "\n\r", " ", -1)
	s = strings.Replace(s, "\n", " ", -1)
	return s
}

// CRLF -> LF
func NoCR(s string) string {
	return strings.Replace(s, "\r\n", "\n", -1)
}

func Truncate(s string, maxChars int) string {
	runes := []rune(s)
	if len(runes) <= maxChars {
		return s
	}
	return string(runes[:maxChars])
}

func TruncateBytes(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		// Common case: return fast
		return s
	}

	buf := []byte(s)
	runes := make([]rune, 0, len(s))
	written := 0
	for len(buf) > 0 {
		rune, n := utf8.DecodeRune(buf)
		if written+n <= maxBytes {
			runes = append(runes, rune)
			written += n
			buf = buf[n:]
		} else {
			return string(runes)
		}
	}
	return string(runes)
}

// Concurrent launches provided funcs, and waits for their completion.
func Concurrent(funcs ...func()) {
	var wg sync.WaitGroup
	wg.Add(len(funcs))
	for _, f := range funcs {
		f := f
		go func() {
			f()
			wg.Done()
		}()
	}
	wg.Wait()
}

// Concurrent launches provided funcs, and returns a channel to notify completion.
func ConcurrentPromise(funcs ...func()) chan bool {
	ch := make(chan bool)
	go func() {
		Concurrent(funcs...)
		ch <- true
		close(ch)
	}()
	return ch
}

// ConcurrentWithAllErrors launches provided funcs, and gathers errors.
// If no errors, ok is true and the returned slice contains all nil values.
func ConcurrentWithAllErrors(funcs ...func() error) (ok bool, errs []error) {
	errs = make([]error, len(funcs))
	var wg sync.WaitGroup
	wg.Add(len(funcs))
	for i, f := range funcs {
		i := i
		f := f
		go func() {
			errs[i] = f()
			wg.Done()
		}()
	}
	wg.Wait()
	ok = true
	for _, err := range errs {
		if err != nil {
			ok = false
		}
	}
	return ok, errs
}

// ConcurrentWithAnyErrors launches provided funcs, and returns
// 1 error if at least 1 error occurred, nil otherwise.
func ConcurrentWithAnyError(funcs ...func() error) error {
	ok, errs := ConcurrentWithAllErrors(funcs...)
	if !ok {
		for _, err := range errs {
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CloneStringSlice makes a defensive copy of
func CloneStringSlice(a []string) []string {
	b := make([]string, len(a))
	copy(b, a)
	return b
}
