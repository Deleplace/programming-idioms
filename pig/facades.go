package pig

import "fmt"

// Those facades are intended to feed HTML templates with renderable data.
// Slices of slices are allowed, but they must be programmaticaly constructed
// out of the GAE Datastore structures.

// See individual controllers for specific pages facades

// PageMeta is basic metadata useful for any web page.
type PageMeta struct {
	// PageTitle is the title of this page
	PageTitle string
	// CanonicalURL helps deduplicating pages for SEO.
	// See https://support.google.com/webmasters/answer/139066
	CanonicalURL string
	// PageKeywords for meta tag "keyword"
	PageKeywords string
	// Toggles (global or custom) used to tune the output
	Toggles Toggles
	// SearchQuery is printed in search field, in case of new similar search
	SearchQuery string
	// ExtraCss after programming-idioms.css
	ExtraCss []string
	// ExtraJs after programming-idioms.js
	ExtraJs []string
	// PreventIndexingRobots for edit pages, etc.
	PreventIndexingRobots bool
	// SEO structured data
	QA QAStructuredData
}

// UserProfile is a soft (non-secure) user profile
type UserProfile struct {
	Nickname          string
	FavoriteLanguages []string
	SeeNonFavorite    bool
	// IsAdmin will never be set by user himself
	IsAdmin bool
}

func (u UserProfile) String() string {
	nonFav := ""
	if !u.SeeNonFavorite {
		nonFav = " (SeeNonFavorite OFF)"
	}
	return fmt.Sprintf("UserProfile[%s %v%s]", u.Nickname, u.FavoriteLanguages, nonFav)
}

// Empty means the visitor doesn't have any meaningful profile data in her cookies.
func (u UserProfile) Empty() bool {
	return u.Nickname == "" &&
		len(u.FavoriteLanguages) == 0 &&
		u.SeeNonFavorite == true &&
		u.IsAdmin == false
}

func EmptyUserProfile() UserProfile {
	return UserProfile{
		Nickname:          "",
		FavoriteLanguages: nil,
		SeeNonFavorite:    true,
		IsAdmin:           false,
	}
}

// LanguageSingleSelector is used to specify the prefilled value
// of a programming language selection widget.
type LanguageSingleSelector struct {
	// Name of the HTML element
	FieldName string
	// Value of the widget: standardized name of a programming language
	Selected string
}

// QAStructuredData is SEO structured data to display Rich results in
// search engines. It maps to a schema.org Q&A page.
// See https://developers.google.com/search/docs/data-types/qapage
type QAStructuredData struct {
	// Question is a language-specific specialization of the Idiom title,
	// e.g. "Depth-first traversing of a binary tree, in Python".
	// It may have the same value as the HTML page title.
	// When Question is empty, it is safe to assume the QAStructuredData is
	// empty/invalid and should not be written to the HTML page.
	Question string

	// Text is the long form of the question. It maps to the Idiom's lead paragraph.
	Text string

	Author string

	DateCreated string

	ImageURL string

	Answers []QAStructuredDataAnswer
}

// QAStructuredDataAnswer is structured data for a single Answer, i.e. a
// single Implementation.
// See https://developers.google.com/search/docs/data-types/qapage
type QAStructuredDataAnswer struct {
	// Text is mapped to the Implementation's CodeBlock (snippet).
	Text string

	Author string

	DateCreated string

	// URL is the absolute nice URL to this specific implementation.
	URL string
}
