package pigae

import "regexp"

// Markup in idioms comments is extremely lightweight: it only
// knows one syntax _x for identifiers.

// It is almost always interpreted client-side,
// except in the RSS feed which is generated server-side.

func markup2HTML(paragraph string) string {
	return emphasize(paragraph)
}

// emphasize the "underscored" identifiers
//
// _x -> <b><i>x</i></b>
func emphasize(sentence string) string {
	// After a word break,
	// an underscore char,
	// and then a group of valid identifier chars.
	re := regexp.MustCompile("\\b_([\\w$]+)")
	return re.ReplaceAllString(sentence, "<b><i>${1}</i></b>")
}
