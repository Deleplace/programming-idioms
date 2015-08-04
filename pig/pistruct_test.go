package pig

import (
	"testing"
)

func sampleIdiom() *Idiom {
	idiom := Idiom{
		Id:            1313,
		OrigId:        1313,
		Title:         "Sample idiom title",
		LeadParagraph: "Sample idiom lead paragraph",
		Version:       1,
		Author:        "Bob",
		Implementations: []Impl{
			Impl{
				Id:            17,
				OrigId:        17,
				Author:        "John",
				LanguageName:  "Java",
				CodeBlock:     `System.println("Sample");`,
				Version:       1,
				AuthorComment: "From class java.lang.System",
			},
			Impl{
				Id:            18,
				OrigId:        18,
				Author:        "Jack",
				LanguageName:  "Python",
				CodeBlock:     `print "Sample"`,
				Version:       1,
				AuthorComment: "Put that in file sample.py",
			},
		},
		ImplCount: 2,
		// etc.
	}
	return &idiom
}

// ---

var splitForSearchingTests = []struct {
	in        string
	normalize bool
	out       []string
}{
	// Empty cases
	{"", false, []string{}},
	{"", true, []string{}},
	// Single word cases
	{"map", false, []string{"map"}},
	{" map", false, []string{"map"}},
	{"map ", false, []string{"map"}},
	{" map  ", false, []string{"map"}},
	{"map", true, []string{"map"}},
	{" map", true, []string{"map"}},
	{"map ", true, []string{"map"}},
	{" map  ", true, []string{"map"}},
	// Multi word cases
	{"map function", false, []string{"map", "function"}},
	{" map  function ", false, []string{"map", "function"}},
	{"map-function", false, []string{"map", "function"}},
	{"map\t\tfunction", false, []string{"map", "function"}},
	{"map function", true, []string{"map", "function"}},
	{" map  function ", true, []string{"map", "function"}},
	{"map-function", true, []string{"map", "function"}},
	{"map\t\tfunction", true, []string{"map", "function"}},
	// Normalization
	{"o'hara", false, []string{"o'hara"}},
	{"o'hara", true, []string{"o", "hara"}},
	{" café ", false, []string{"café"}},
	{" café ", true, []string{"cafe"}},
}

func TestSplitForSearching(t *testing.T) {
	for i, tt := range splitForSearchingTests {
		out := SplitForSearching(tt.in, tt.normalize)
		if !StringSliceEquals(out, tt.out) {
			t.Errorf("%d. SplitForSearching(%v, %v) => %v, want %v", i, tt.in, tt.normalize, out, tt.out)
		}
	}
}

// ---

// Title, Lead, Code must be indexed
var containedWords = []string{"sample", "title", "lead", "java"}

// Authors are not indexed
var notContainedWords = []string{"", "john", "jack"}

func TestIndexWords(t *testing.T) {
	// TODO testwith full text api
}
