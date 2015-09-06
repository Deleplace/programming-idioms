package pigae

import (
	"strings"
)

// PrettyAdaptor contains a CSS class and a syntax coloring JS script, for some language.
type PrettyAdaptor struct {
	CssClass    string
	JsExtension string
}

// Pretty maps some languages to a syntax coloring JS script.
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

func prettifyCSSClass(lang string) string {
	lg := strings.TrimSpace(strings.ToLower(lang))
	suff := lg
	var p PrettyAdaptor
	var ok bool
	// See http://google-code-prettify.googlecode.com/svn/trunk/README.html
	if p, ok = Pretty[strings.ToLower(normLang(lang))]; !ok {
		return "lang-" + suff
	}
	return "lang-" + p.CssClass
}

// Just returns "" for no extension
// Not used anymore, see prettify-extra-languages.min.js
func prettifyExtension(lang string) string {
	p := Pretty[strings.ToLower(normLang(lang))]
	return p.JsExtension
}
