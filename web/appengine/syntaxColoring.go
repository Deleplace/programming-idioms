package main

import (
	"strings"

	. "github.com/Deleplace/programming-idioms/idioms"
)

// PrettyAdaptor contains a CSS class and a syntax coloring JS script, for some language.
// OBSOLETE (using hljs now)
type PrettyAdaptor struct {
	CssClass    string
	JsExtension string
}

// Pretty maps some languages to a syntax coloring JS script.
// OBSOLETE (using hljs now)
var Pretty = map[string]PrettyAdaptor{
	"python":  {"py", ""},
	"ruby":    {"rb", ""},
	"csharp":  {"cs", ""},
	"go":      {"go", "lang-go.js"},
	"basic":   {"basic", "lang-basic.js"},
	"dart":    {"dart", "lang-dart.js"},
	"erlang":  {"erlang", "lang-erlang.js"},
	"lisp":    {"lisp", "lang-lisp.js"},
	"lua":     {"lua", "lang-lua.js"},
	"matlab":  {"matlab", "lang-matlab.js"},
	"pascal":  {"pascal", "lang-pascal.js"},
	"r":       {"r", "lang-r.js"},
	"scala":   {"scala", "lang-scala.js"},
	"scheme":  {"scm", "lang-lisp.js"},
	"haskell": {"hs", "lang-hs.js"},
	"clojure": {"clj", "lang-clj.js"},
}

// OBSOLETE (using hljs now)
func prettifyCSSClass(lang string) string {
	lg := strings.TrimSpace(strings.ToLower(lang))
	suff := lg
	var p PrettyAdaptor
	var ok bool
	// See https://github.com/google/code-prettify/blob/master/README.md
	if p, ok = Pretty[strings.ToLower(NormLang(lang))]; !ok {
		return "lang-" + suff
	}
	return "lang-" + p.CssClass
}

// Just returns "" for no extension
// Not used anymore, see prettify-extra-languages.min.js
// OBSOLETE (using hljs now)
func prettifyExtension(lang string) string {
	p := Pretty[strings.ToLower(NormLang(lang))]
	return p.JsExtension
}

// Highlight.js maps some languages under a slightly different name
var hljsMapping = map[string]string{
	"js":    "language-javascript",
	"caml":  "language-ocaml",
	"obj-c": "language-objectivec",
	// Some are not supported by Highlight.js
	"cobol":  "language-plaintext",
	"pascal": "language-plaintext",
}

func hljsCSSClass(lang string) string {
	lg := strings.TrimSpace(strings.ToLower(NormLang(lang)))
	if hlClass, ok := hljsMapping[lg]; ok {
		// E.g. "js" -> "language-javascript"
		return hlClass
	}
	// Default, e.g. "python" -> "language-python"
	return "language-" + lg
}
