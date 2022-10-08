package main

import (
	"html/template"
	"regexp"
	"strings"
)

// Markup in idioms comments is extremely lightweight: it only
// knows one syntax _x for identifiers.

// It may be interpreted client-side (in a Preview),
// or server-side in HTML pages and RSS feed.

func markup2HTML(paragraph string) string {
	return emphasize(paragraph)
}

func markup2CSS(paragraph string) template.HTML {
	untrusted := paragraph
	sanitized := template.HTMLEscapeString(untrusted)
	emphasized := emphasizeCSS(sanitized)
	linebroken := linebreak(emphasized)
	return template.HTML(linebroken)
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

// emphasize the "underscored" identifiers
//
// _x -> <em>x</em>
func emphasizeCSS(sentence string) string {
	// After a word break,
	// an underscore char,
	// and then a group of valid identifier chars.
	re := regexp.MustCompile("\\b_([\\w$]+)")
	return re.ReplaceAllString(sentence, "<em>${1}</em>")
}

// "a\nb" -> "a<br/>b"
func linebreak(sentence string) string {
	sentence = strings.Replace(sentence, "\r\n", "<br/>", -1)
	sentence = strings.Replace(sentence, "\n\r", "<br/>", -1)
	sentence = strings.Replace(sentence, "\n", "<br/>", -1)
	return sentence
}
