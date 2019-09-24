package main

import (
	"fmt"
	"strings"

	. "github.com/Deleplace/programming-idioms/pig"
)

func uriNormalize(s string) string {
	s = strings.Trim(NormalizeRunes(s), " -/")
	s = strings.Replace(s, " ", "-", -1)
	s = strings.Replace(s, "[", "-", -1)
	s = strings.Replace(s, "]", "-", -1)
	s = strings.Replace(s, ",", "-", -1)
	s = strings.Replace(s, ";", "-", -1)
	s = strings.Replace(s, "--", "-", -1)
	s = strings.Replace(s, "--", "-", -1) // again ;)
	s = strings.Trim(s, " -/")
	return s
}

// NiceIdiomURL produces a canonical URL for an Idiom.
func NiceIdiomURL(idiom *Idiom) string {
	// Note : cannot define new methods on non-local type pig.Idiom
	// so idiom is first argument
	return fmt.Sprintf("%v/idiom/%v/%v", hostPrefix(), idiom.Id, uriNormalize(idiom.Title))
}

// NiceImplURL produces a canonical URL for an Impl.
func NiceImplURL(idiom *Idiom, implID int, implLang string) string {
	// Note : cannot define new methods on non-local type pig.Idiom
	// so idiom is first argument
	return fmt.Sprintf("%v/%v/%v", NiceIdiomURL(idiom), implID, uriNormalize(implLang))
}

// NiceImplRelativeURL produces a relative canonical URL for an Impl.
func NiceImplRelativeURL(idiom *Idiom, implID int, implLang string) string {
	// Note : cannot define new methods on non-local type pig.Idiom
	// so idiom is first argument
	return fmt.Sprintf("%v/%v/%v", NiceIdiomRelativeURL(idiom), implID, uriNormalize(implLang))
}

// NiceIdiomIDTitleURL produces an URL for an Idiom, with specified title.
func NiceIdiomIDTitleURL(idiomID int, title string) string {
	return fmt.Sprintf("%v/idiom/%v/%v", hostPrefix(), idiomID, uriNormalize(title))
}

// NiceIdiomRelativeURL produces a relative canonical URL for an Idiom.
func NiceIdiomRelativeURL(idiom *Idiom) string {
	return fmt.Sprintf("/idiom/%v/%v", idiom.Id, uriNormalize(idiom.Title))
}

// UglyIdiomIDURL produces an URL for an Idiom ID.
func UglyIdiomIDURL(idiomID int) string {
	return fmt.Sprintf("%v/idiom/%v", hostPrefix(), idiomID)
}

// UglyIdiomIDURL produces an URL for an idiom ID and an impl ID.
func UglyIdiomIDImplIDURL(idiomID, implID int) string {
	return fmt.Sprintf("%v/idiom/%v/impl/%v", hostPrefix(), idiomID, implID)
}
