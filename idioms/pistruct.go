package idioms

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
	// Id is auto-incremented 1, 2, 3...
	// TODO name ID instead
	Id int

	// Reserved in case one idiom derived from another
	// TODO name OrigID instead
	OrigId int

	// Title is like the "idiom name"
	Title string

	// LeadParagraph is the idiom Description : 1 to 3 lines are fine
	LeadParagraph string

	// ExtraKeywords for indexation+search
	ExtraKeywords string

	// Author is the name of the original creator of this idiom on this website
	Author string

	// CreationDate is the date of creation of this idiom on this site
	CreationDate time.Time

	// LastEditor is the name of the last person who modified this idiom
	LastEditor string

	// EditSummary is the comment explaining why LastEditor made the edit.
	// It is not displayed except in history views.
	EditSummary string

	// LastEditedImplID is the ID of the only impl modified by
	// last edit. LastEditedImplID should be 0 if last
	// edit was on the idiom statement, not on an impl.
	LastEditedImplID int

	// Please acknowledge sources (idiom statement, not snippet).
	OriginalAttributionURL string

	// Picture representing the concept, if necessary
	// DEPRECATED
	Picture string

	// ImageURL to illustrate this idiom.
	ImageURL string

	// ImageWidth, ImageHeight hints for web rendering: avoids FOUC
	ImageWidth, ImageHeight int

	// ImageAlt is a description of a picture for accessibility purpose
	ImageAlt string

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

	// Protected when "only admin can edit"
	Protected bool

	// Variables from the lead paragraph, that we'd like every
	// impl snippet to contain
	Variables []string

	// RelatedURLs as extra pieces of documentation for the idiom statement.
	RelatedURLs []string

	// RelatedURLLabels are nice text for Related URLs hyperlinks.
	RelatedURLLabels []string
}

// Impl is a specific implementation of one Idiom in one programming language.
// It is in theory independent from any framework, but has been used only in
// Google App Engine so far.
type Impl struct {
	// Id is Internal. Not displayed on screen (but present in URL).
	// TODO name ID instead
	Id int

	// OrigId is reserved in case one impl derived from another
	// TODO name OrigID instead
	OrigId int

	// Author is the name of the original creator of this implementation on this site.
	Author string

	// CreationDate of this implementation on this website
	CreationDate time.Time

	// LastEditor is the name of the last person who modified this impl.
	LastEditor string

	// LanguageName is the programming language of this impl.
	// It is used to visualy identify the impl inside the idiom.
	// But note that an idiom may have several implementations for same language.
	LanguageName string

	// CodeBlock contains the snippet.
	// It should contain only instructions code, not comments.
	CodeBlock string

	// OriginalAttributionURL: please acknowledge sources.
	OriginalAttributionURL string

	// DemoURL is an optional link to an online demo
	DemoURL string

	// DocumentationURL is an optional link to official doc
	DocumentationURL string

	// AuthorComment comments about the CodeBlock.
	// This comment is always displayed on the right of the code.
	// TODO rename this to CodeBlockComment.
	AuthorComment string

	// Version is incremented at each update 1, 2, 3...
	Version int

	// VersionDate of last update
	VersionDate time.Time

	// Rating is the votes count for this specific impl  (votes up - votes down)
	Rating int

	// Checked is true if an admin has validated this implementation.
	Checked bool

	// ImplRenderingDecoration is some extra calculated data like "Has this implementation been upvoted by this user?"
	// Ignored by the datastore.
	Deco ImplRenderingDecoration `datastore:"-" json:"-"`

	// ImportsBlock contains the import directives, appart from main code section.
	ImportsBlock string

	// PictureURL to illustrate this impl.
	PictureURL string

	// Protected when "only admin can edit"
	Protected bool
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
	// Matching is set to true if current impl matches user text search query.
	Matching bool
	// SearchedLang is set to true if current impl lang is the user typed lang.
	SearchedLang bool
}

// IdiomVoteLog is a history trace of an Idiom vote, from a specific user.
// This struct does not contain the nickname of the voter.
// However it does contain the Idiom ID.
// Each vote will have a voting booth as ancestor, specific for the nickname.
type IdiomVoteLog struct {
	IdiomId int
	// Typicaly +1 or -1
	Value int
	// IpHash stored only to prevent abusive multiple votes
	IpHash string
	Date   time.Time
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
	// IpHash stored only to prevent abusive multiple votes
	IpHash string
	Date   time.Time
}

// IdiomHistory stores all the history: old versions of Idioms.
type IdiomHistory struct {
	// Just embeds Idiom
	Idiom
	// If needed, add specific history fields
	UpdatedImplId int
	// TODO: how to get rid properly?
	// Got `datastore: cannot load field "EditorSummary" into a "pig.IdiomHistory": no such struct field`
	EditorSummary string `deprecated`
	// IdiomOrImplLastEditor is redundant storage of most recent impl update's editor,
	// to be directly indexed and displayed in history list.
	IdiomOrImplLastEditor string
}

func (ih *IdiomHistory) AsIdiomPtr() *Idiom {
	return &(ih.Idiom)
}

type MessageForUser struct {
	CreationDate,
	FirstViewDate,
	LastViewDate,
	DismissalDate,
	ExpirationDate time.Time
	Message  string
	Username string
}

/* ---- */

// FindImplInIdiom is a (unoptimized) iteration to retrieve an Impl by its ID,
// inside an Idiom.
//
// It returns a pointer to the Impl, not a copy.
func (idiom *Idiom) FindImplInIdiom(implId int) (int, *Impl, bool) {
	for i := range idiom.Implementations {
		impl := &idiom.Implementations[i]
		if impl.Id == implId {
			return i, impl, true
		}
	}
	return -1, nil, false
}

// FindRecentlyUpdatedImpl is a (unoptimized) iteration to retrieve
// the most recently updated impl, inside an Idiom.
//
// It returns a pointer to the Impl, not a copy.
func (idiom *Idiom) FindRecentlyUpdatedImpl() *Impl {
	var recentImpl *Impl
	var d time.Time
	for i := range idiom.Implementations {
		impl := &idiom.Implementations[i]
		if impl.VersionDate.After(d) {
			recentImpl = impl
			d = impl.VersionDate
		}
	}
	return recentImpl
}

// ExtractIndexableWords compute the list of words contained in an Idiom.
// First return value is the list of all matchable words.
// Second return value is the list of matchable words from title only.
func (idiom *Idiom) ExtractIndexableWords() (w []string, wTitle []string, wLead []string) {
	w = SplitForIndexing(idiom.Title, true)
	w = append(w, fmt.Sprintf("%d", idiom.Id))
	wTitle = w

	wLead = SplitForIndexing(idiom.LeadParagraph, true)
	w = append(w, wLead...)
	// ExtraKeywords as not as important as Title, rather as important as Lead
	wKeywords := SplitForIndexing(idiom.ExtraKeywords, true)
	wLead = append(wLead, wKeywords...)
	w = append(w, wKeywords...)

	for i := range idiom.Implementations {
		impl := &idiom.Implementations[i]
		wImpl := impl.ExtractIndexableWords()
		w = append(w, wImpl...)
	}

	return w, wTitle, wLead
}

// ExtractIndexableWords compute the list of words contained in an Impl.
func (impl *Impl) ExtractIndexableWords() []string {
	w := make([]string, 0, 20)
	w = append(w, fmt.Sprintf("%d", impl.Id))
	w = append(w, strings.ToLower(impl.LanguageName))
	w = append(w, SplitForIndexing(impl.ImportsBlock, true)...)
	w = append(w, SplitForIndexing(impl.CodeBlock, true)...)
	if len(impl.AuthorComment) >= 3 {
		w = append(w, SplitForIndexing(impl.AuthorComment, true)...)
	}
	if langExtras, ok := langsExtraKeywords[impl.LanguageName]; ok {
		w = append(w, langExtras...)
	}
	// Note: we don't index external URLs.
	return w
}

// FindIdiomOrImplLastEditor returns last user who touched something
func (idiom *Idiom) FindIdiomOrImplLastEditor() string {
	last := idiom.LastEditor
	mostRecentDate := idiom.VersionDate
	for _, impl := range idiom.Implementations {
		lame := -time.Second // lame delta because idiom.VersionDate is always greater
		if !impl.VersionDate.Before(mostRecentDate.Add(lame)) {
			last = impl.LastEditor
			mostRecentDate = impl.VersionDate
		}
	}
	return last
}

// ComputeIdiomOrImplLastEditor computes IdiomOrImplLastEditor
func (hist *IdiomHistory) ComputeIdiomOrImplLastEditor() {
	hist.IdiomOrImplLastEditor = hist.LastEditor
	mostRecentDate := hist.VersionDate
	for _, impl := range hist.Implementations {
		lame := -time.Second // lame delta because hist.VersionDate is always greater
		if !impl.VersionDate.Before(mostRecentDate.Add(lame)) {
			hist.IdiomOrImplLastEditor = impl.LastEditor
			mostRecentDate = impl.VersionDate
		}
	}
}

var regexpWhiteSpace = regexp.MustCompile("[ \\t\\n]")
var regexpWhiteSpaceDash = regexp.MustCompile("[ \\t\\n-]")
var RegexpDigitsOnly = regexp.MustCompile("^\\d+$")

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
		if len(chunk) >= 3 || RegexpDigitsOnly.MatchString(chunk) {
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
	chunks = FilterStrings(chunks, func(chunk string) bool {
		// Empty strings are no good.
		return len(chunk) >= 1
	})

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
		idiom.EditSummary = fmt.Sprintf("Linked to idiom #%d [%v]", other.Id, other.Title)
	}
	if !containsInt(other.RelatedIdiomIds, idiom.Id) {
		other.RelatedIdiomIds = append(other.RelatedIdiomIds, idiom.Id)
		other.RelatedIdiomTitles = append(other.RelatedIdiomTitles, idiom.Title)
		other.EditSummary = fmt.Sprintf("Linked to idiom #%d [%v]", idiom.Id, idiom.Title)
	}
}

// VariablesComma e.g. ["x", "result"] -> "x,result"
func (idiom *Idiom) VariablesComma() string {
	return strings.Join(idiom.Variables, ",")
}
