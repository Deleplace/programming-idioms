package pig

// Those facades are intended to feed HTML templates with renderable data.
// Slices of slices are allowed, but they must be programmaticaly constructed
// out of the GAE Datastore structures.

// See individual controllers for specific pages facades

// PageMeta is basic metadata useful for any web page.
type PageMeta struct {
	PageTitle   string
	Toggles     Toggles
	SearchQuery string
	ExtraCss    []string
	ExtraJs     []string
}

// UserProfile is a soft (non-secure) user profile
type UserProfile struct {
	Nickname          string
	FavoriteLanguages []string
	SeeNonFavorite    bool
	// Field IsAdmin will never be set by user himself
	IsAdmin bool
}

// LanguageSingleSelector is used to specify the prefilled value
// of a programming language selection widget.
type LanguageSingleSelector struct {
	// Name of the HTML element
	FieldName string
	// Value of the widget: standardized name of a programming language
	Selected string
}
