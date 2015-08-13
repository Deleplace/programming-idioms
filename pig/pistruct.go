package pig

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Note : with the GAE datastore it is *not* possible
// to have a slice of slices inside a struct.

// Idiom is the main entity of programming-idioms.org .
// An Idiom contains its Implementations.
// It is in theory independent from any framework, but has been used only in
// Google App Engine so far.
type Idiom struct {
	// Auto-incremented 1, 2, 3...
	Id int
	// Reserved in case one idiom derived from another
	OrigId int
	// Title : the "idiom name"
	Title string
	// Idiom Description : 1 to 3 lines are fine
	LeadParagraph string
	// The name of the original creator of this idiom on this site
	Author string
	// The date of creation of this idiom on this site
	CreationDate time.Time
	// The name of the last person who modified this idiom
	LastEditor string
	// Please acknowledge sources (idiom statement, not snippet).
	OriginalAttributionURL string
	// Picture representing the concept, if necessary
	Picture  string
	ImageURL string
	// Autoincremented at each update 1, 2, 3...
	Version int
	// Date of last update
	VersionDate time.Time
	// List of implementations of this idiom in specific languages
	Implementations []Impl
	// (Denormalized) number of contained implementation, for datastore querying
	ImplCount int
	// How many votes for the idiom itself  (votes up - votes down)
	Rating int
	// Index-like array of important words : those from the title
	// DEPRECATED: use the new Text Search API instead.
	WordsTitle []string
	// Index-like array of words from title, description and implementation contents
	// DEPRECATED: use the new Text Search API instead.
	Words []string
	// Did the admin validate this idiom statement ?
	Checked bool
	// Extra calculated data like "Has this idiom been upvoted by this user?"
	// Ignored by the datastore.
	Deco IdiomRenderingDecoration `datastore:"-" json:"-"`
	// Related idioms ids "See also..."
	RelatedIdiomIds []int
	// NoSQL-style : store directly some data from other objects
	RelatedIdiomTitles []string
}

// Impl is a specific implementation of one Idiom in one programming language.
// It is in theory independent from any framework, but has been used only in
// Google App Engine so far.
type Impl struct {
	// Internal Id. Not displayed (but present in URL).
	Id int
	// Reserved in case one impl derived from another
	OrigId int
	// The name of the original creator of this implementation on this site.
	Author string
	// The date of creation of this implementation on this site
	CreationDate time.Time
	// The name of the last person who modified this impl.
	LastEditor string
	// The programming language of this impl.
	// It is used to visualy identify the impl inside the idiom.
	// But note that an idiom may have several implementations for same language.
	LanguageName string
	// The snippet.
	// Should contain only code, preferably no comments.
	CodeBlock string
	// Please acknowledge sources.
	OriginalAttributionURL string
	DemoURL                string
	// Editor comments
	AuthorComment string
	// Autoincremented at each update 1, 2, 3...
	Version int
	// Date of last update
	VersionDate time.Time
	// How many votes for this specific impl  (votes up - votes down)
	Rating int
	// Did the admin validate this implementation ?
	Checked bool
	// Extra calculated data like "Has this implementation been upvoted by this user?"
	// Ignored by the datastore.
	Deco ImplRenderingDecoration `datastore:"-" json:"-"`
	// Prerequisites : appart from main code section
	ImportsBlock string
}

// IdiomRenderingDecoration is the "current user" vote on this Idiom, if any.
// This struct does not contain the Idiom ID, so it must be part of a larger struct.
type IdiomRenderingDecoration struct {
	UpVoted   bool
	DownVoted bool
}

// ImplRenderingDecoration is the "current user" vote on this Impl, if any.
// This struct does not contain Impl ID nor Idiom ID, so it must be part of a larger struct.
type ImplRenderingDecoration struct {
	UpVoted   bool
	DownVoted bool
}

// IdiomVoteLog is a history trace of an Idiom vote, from a specific user.
// This struct does not contain the nickname of the voter.
// However it does contain the Idiom ID.
// Each vote will have a voting booth as ancestor, specific for the nickname.
type IdiomVoteLog struct {
	IdiomId int
	// Typicaly +1 or -1
	Value int
	// Stored only to prevent abusive multiple votes
	IpHash string
}

// ImplVoteLog is a history trace of an Impl vote, from a specific user.
// This structure does not contain the nickname of the voter.
// However it does contain the Idiom ID and Impl ID.
// Each vote will have a voting booth as ancestor, specific for the nickname.
type ImplVoteLog struct {
	IdiomId int
	ImplId  int
	// Typicaly +1 or -1
	Value int
	// Stored only to prevent abusive multiple votes
	IpHash string
}

// IdiomHistory stores all the history: old versions of Idioms.
type IdiomHistory struct {
	// Just embeds Idiom
	Idiom
	// If needed, add specific history fields
}

func (ih *IdiomHistory) AsIdiomPtr() *Idiom {
	return &(ih.Idiom)
}

/* ---- */

// FindImplInIdiom is a (unoptimized) iteration to retrieve an Impl by its ID,
// inside an Idiom.
func (idiom *Idiom) FindImplInIdiom(implId int) (int, Impl, bool) {
	impl := Impl{}
	for i, impl := range idiom.Implementations {
		if impl.Id == implId {
			return i, impl, true
		}
	}
	return -1, impl, false
}

// ExtractIndexableWords compute the list of words contained in an Idiom.
// First return value is the list of all matchable words.
// Second return value is the list of matchable words from title only.
func (idiom *Idiom) ExtractIndexableWords() ([]string, []string) {
	w := SplitForIndexing(idiom.Title, true)
	w = append(w, fmt.Sprintf("%d", idiom.Id))
	wTitle := w
	w = append(w, SplitForIndexing(idiom.LeadParagraph, true)...)
	for _, impl := range idiom.Implementations {
		w = append(w, SplitForIndexing(impl.CodeBlock, true)...)
		if len(impl.AuthorComment) >= 3 {
			w = append(w, SplitForIndexing(impl.AuthorComment, true)...)
		}
		w = append(w, strings.ToLower(impl.LanguageName))
		w = append(w, fmt.Sprintf("%d", impl.Id))
	}
	return w, wTitle
}

var regexpWhiteSpace = regexp.MustCompile("[ \\t\\n]")
var regexpWhiteSpaceDash = regexp.MustCompile("[ \\t\\n-]")
var regexpDigitsOnly = regexp.MustCompile("^\\d+$")

// SplitForIndexing cuts sentences or paragrahs into words.
// Words of 2 letters of less are discarded.
func SplitForIndexing(s string, normalize bool) []string {
	if normalize {
		s = NormalizeRunes(s)
	}
	chunks := regexpWhiteSpace.Split(s, -1)
	realChunks := make([]string, 0, len(chunks))

	for _, chunk := range chunks {
		// Accepted :
		// All words having at least 3 characters
		// All 1-digits words and 2-digits words
		if len(chunk) >= 3 || regexpDigitsOnly.MatchString(chunk) {
			realChunks = append(realChunks, NormalizeRunes(chunk))
		}
	}

	// Stategy for dash-compound words: all bits get indexed (in addition to the full compound)
	for _, chunk := range chunks {
		if strings.Contains(chunk, "-") {
			for _, bit := range strings.Split(chunk, "-") {
				if bit != "" {
					realChunks = append(realChunks, bit)
				}
			}
		}
	}

	return realChunks
}

// SplitForSearching cuts an input search string into a slice of search terms.
func SplitForSearching(s string, normalize bool) []string {
	if normalize {
		s = NormalizeRunes(s)
	}
	chunks := regexpWhiteSpaceDash.Split(s, -1)
	chunks = FilterOut(chunks, []string{""})
	// All typed chunk are considered acceptable search terms
	return chunks
}

// NormalizeRunes discard special characters from a string, for indexing and for searching.
// Some letters with diacritics are replaced by the same letter without diacritics.
func NormalizeRunes(str string) string {
	str = strings.ToLower(str)
	norm := func(r rune) rune {
		switch r {
		// TODO find a standard golang normalization ?
		case ' ', '\t', '(', ')', '"', '\'', ',', ';', ':', '?', '.', '/', '+':
			return ' '
		case '%', '^', '=', '`', '*', '&', '!', '°', '_':
			return ' '
		case 'à', 'ä':
			return 'a'
		case 'ç':
			return 'c'
		case 'é', 'è', 'ê', 'ë':
			return 'e'
		case 'ï', 'î':
			return 'i'
		case 'ô', 'ö':
			return 'o'
		case 'û', 'ü':
			return 'u'
		}
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == '-':
			return r
		}
		// Unknown characters should not be allowed in
		return -1
	}
	return strings.Map(norm, str)
}

func containsInt(a []int, x int) bool {
	for _, i := range a {
		if i == x {
			return true
		}
	}
	return false
}

// AddRelation creates a bidirectional link between 2 related Idioms.
func (idiom *Idiom) AddRelation(other *Idiom) {
	if !containsInt(idiom.RelatedIdiomIds, other.Id) {
		idiom.RelatedIdiomIds = append(idiom.RelatedIdiomIds, other.Id)
		idiom.RelatedIdiomTitles = append(idiom.RelatedIdiomTitles, other.Title)
	}
	if !containsInt(other.RelatedIdiomIds, idiom.Id) {
		other.RelatedIdiomIds = append(other.RelatedIdiomIds, idiom.Id)
		other.RelatedIdiomTitles = append(other.RelatedIdiomTitles, idiom.Title)
	}
}
