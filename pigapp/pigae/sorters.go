package pigae

import (
	"sort"

	. "github.com/Deleplace/programming-idioms/pig"
)

// Parameter favoriteLanguages contain user favorite languages, in decreasing order of
// interest (favoriteLanguages[0] is the most important, etc.).
// Puts favorites first, but does not care about which are "most favorite"
func favoritesFirst(languages []string, favoriteLanguages []string) {
	// TODO make an idiom of that algo :)
	n := len(languages)
	k := 0
	i := 1
	for i < n {
		if k < i && StringSliceContains(favoriteLanguages, languages[i]) {
			tmp := languages[k]
			languages[k] = languages[i]
			languages[i] = tmp
			k++
		} else {
			i++
		}
	}
}

// Parameter favoriteLanguages contain user favorite languages, in decreasing order of
// interest (favoriteLanguages[0] is the most important, etc.).
// Inspired by planetSorter
// in https://golang.org/pkg/sort/
type languageSorter struct {
	languages []string
	scores    map[string]int
	implCount map[string]int
}

func (s *languageSorter) Len() int {
	return len(s.languages)
}
func (s *languageSorter) Swap(i, j int) {
	s.languages[i], s.languages[j] = s.languages[j], s.languages[i]
}
func (s *languageSorter) Less(i, j int) bool {
	a, isfavA := s.scores[s.languages[i]]
	b, isfavB := s.scores[s.languages[j]]
	if isfavA {
		if isfavB {
			return a < b
		} else {
			return true
		}
	} else {
		if isfavB {
			return false
		} else {
			return s.implCount[s.languages[i]] >= s.implCount[s.languages[j]]
		}
	}
}

// Parameter favoriteLanguages contain user favorite languages, in decreasing order of
// interest (favoriteLanguages[0] is the most important, etc.).
func favoritesFirstWithOrder(languages []string, favoriteLanguages []string, implCount map[string]int) {
	// TODO unit test all this
	scores := map[string]int{}
	for i, lang := range favoriteLanguages {
		scores[lang] = i
	}
	ls := &languageSorter{
		languages: languages,
		scores:    scores,
		implCount: implCount,
	}
	sort.Sort(ls)
}

// Inspired by planetSorter
// in http://golang.org/pkg/sort/
type implByLanguageSorter struct {
	impl   []Impl
	scores map[string]int
}

func (s *implByLanguageSorter) Len() int {
	return len(s.impl)
}
func (s *implByLanguageSorter) Swap(i, j int) {
	s.impl[i], s.impl[j] = s.impl[j], s.impl[i]
}
func (s *implByLanguageSorter) Less(i, j int) bool {
	a, isfavA := s.scores[s.impl[i].LanguageName]
	b, isfavB := s.scores[s.impl[j].LanguageName]
	if isfavA == isfavB {
		if isfavA {
			// both fav: user-defined pref order
			return a < b
		} else {
			// both non-fav: alphabetical order of lang name
			return s.impl[i].LanguageName < s.impl[j].LanguageName
		}
	} else {
		// all fav before all non-fav
		return isfavA
	}
}

// Puts favorites first, but does not care about which are "most favorite"
// This alters the ordering of the slice.
// Use to display data, not to persist it.
func implFavoriteLanguagesFirst(implementations []Impl, favoriteLanguages []string) {
	// TODO make an idiom of that algo :)
	n := len(implementations)
	k := 0
	i := 1
	for i < n {
		if k < i && StringSliceContains(favoriteLanguages, implementations[i].LanguageName) {
			tmp := implementations[k]
			implementations[k] = implementations[i]
			implementations[i] = tmp
			k++
		} else {
			i++
		}
	}
}

// Parameter favoriteLanguages contain user favorite languages, in decreasing order of
// interest (favoriteLanguages[0] is the most important, etc.).
func implFavoriteLanguagesFirstWithOrder(idiom *Idiom, favoriteLanguages []string, selectedLang string, includeNonFav bool) {
	// TODO unit test all this
	implementations := idiom.Implementations
	scores := map[string]int{}
	for i, lang := range favoriteLanguages {
		scores[lang] = i
	}
	if selectedLang != "" {
		// Selected language will be first in the list
		scores[selectedLang] = -1
	}

	if !includeNonFav {
		// Let's remove all uninteresting impls
		n := len(implementations)
		for i := 0; i < n; {
			_, fav := scores[implementations[i].LanguageName]
			if fav {
				i++
			} else {
				implementations[i] = implementations[n-1]
				n--
			}
		}
		implementations = implementations[:n]
		idiom.Implementations = implementations
	}

	is := &implByLanguageSorter{
		impl:   implementations,
		scores: scores,
	}
	sort.Sort(is)
}

func sortIdiomsByRating(idioms []*Idiom) {
	sort.Sort(&idiomByRatingSorter{idioms})
}

type idiomByRatingSorter struct {
	idioms []*Idiom
}

func (s *idiomByRatingSorter) Len() int {
	return len(s.idioms)
}
func (s *idiomByRatingSorter) Swap(i, j int) {
	s.idioms[i], s.idioms[j] = s.idioms[j], s.idioms[i]
}
func (s *idiomByRatingSorter) Less(i, j int) bool {
	// Higher ratings first
	return s.idioms[i].Rating > s.idioms[j].Rating
}

func sortIdiomsByVersionDate(idioms []*Idiom) {
	sort.Sort(&idiomByVersionDateSorter{idioms})
}

type idiomByVersionDateSorter struct {
	idioms []*Idiom
}

func (s *idiomByVersionDateSorter) Len() int {
	return len(s.idioms)
}
func (s *idiomByVersionDateSorter) Swap(i, j int) {
	s.idioms[i], s.idioms[j] = s.idioms[j], s.idioms[i]
}
func (s *idiomByVersionDateSorter) Less(i, j int) bool {
	// Higher version dates first
	return s.idioms[i].VersionDate.After(s.idioms[j].VersionDate)
}

// Sort implementations : recently updated first (used for RSS)
type implByVersionDateSorter struct {
	impl []Impl
}

func (s *implByVersionDateSorter) Len() int {
	return len(s.impl)
}
func (s *implByVersionDateSorter) Swap(i, j int) {
	s.impl[i], s.impl[j] = s.impl[j], s.impl[i]
}
func (s *implByVersionDateSorter) Less(i, j int) bool {
	return s.impl[i].VersionDate.Unix() > s.impl[j].VersionDate.Unix()
}
