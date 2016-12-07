package pig

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

//
// Language names exist in 3 forms : nice, standard, lowercase
// Ex : "C++", "Cpp", "cpp"
//

var mainStreamLangs = [...]string{"C", "Cpp", "Csharp", "Go", "Java", "JS", "Obj-C", "PHP", "Python", "Ruby"}

// Return alpha codes for each language (no encoding problems).
// See PrintNiceLang to display them more fancy.
func MainStreamLanguages() []string {
	return mainStreamLangs[:]
}

var moreLangs = [...]string{"Ada", "Caml", "Clojure", "Cobol", "D", "Dart", "Elixir", "Erlang", "Fortran", "Haskell", "Lua", "Lisp", "Pascal", "Perl", "Prolog", "Rust", "Scala", "Scheme", "VB"}

func MoreLanguages() []string {
	// These do *not* include the MainStreamLanguages()
	return moreLangs[:]
}

var synonymLangs = map[string]string{
	"Javascript":   "JS",
	"Objective C":  "Obj-C",
	"Visual Basic": "VB",
}

var AllLangs []string
var AllNiceLangs []string

func AllLanguages() []string {
	if AllLangs == nil {
		mainstream := MainStreamLanguages()
		more := MoreLanguages()
		AllLangs = make([]string, len(mainstream)+len(more))
		copy(AllLangs, mainstream)
		copy(AllLangs[len(mainstream):], more)
		sort.Strings(AllLangs)
		AllNiceLangs = make([]string, len(AllLangs))
		for i, lg := range AllLangs {
			AllNiceLangs[i] = PrintNiceLang(lg)
		}
	}
	return AllLangs
}

// autocompletions is a map[string][]string
var autocompletions = precomputeAutocompletions()

func LanguageAutoComplete(fragment string) []string {
	fragment = strings.ToLower(fragment)

	// Dynamic search (slow)
	// options := []string{}
	// for _, lg := range AllLanguages() {
	// 	if strings.Contains(strings.ToLower(lg), fragment) || strings.Contains(strings.ToLower(PrintNiceLang(lg)), fragment) {
	// 		options = append(options, PrintNiceLang(lg))
	// 	}
	// }
	// return options

	// Precomputed search (fast)
	return autocompletions[fragment]
}

func PrintNiceLang(lang string) string {
	switch strings.TrimSpace(strings.ToLower(lang)) {
	case "cpp":
		return "C++"
	case "csharp":
		return "C#"
	default:
		return lang
	}
}

func PrintNiceLangs(langs []string) []string {
	nice := make([]string, len(langs))
	for i, lang := range langs {
		nice[i] = PrintNiceLang(lang)
	}
	return nice
}

func PrintShortLang(lang string) string {
	switch strings.TrimSpace(strings.ToLower(lang)) {
	case "clojure":
		return "Clj"
	case "cobol":
		return "Co bol"
	case "cpp":
		return "C++"
	case "csharp":
		return "C#"
	case "erlang":
		return "Er lang"
	case "elixir":
		return "Eli xir"
	case "fortran":
		return "For tran"
	case "haskell":
		return "Has kell"
	case "obj-c":
		return "Obj C"
	case "pascal":
		return "Pas"
	case "python":
		return "Py"
	case "scheme":
		return "scm"
	case "prolog":
		return "Pro log"
	default:
		return lang
	}
}

func indexByLowerCase(langs []string) map[string]string {
	m := map[string]string{}
	for _, lg := range langs {
		m[strings.ToLower(lg)] = lg
	}
	return m
}

var langLowerCaseIndex = indexByLowerCase(AllLanguages())

func NormLang(lang string) string {
	lg := strings.TrimSpace(strings.ToLower(lang))
	switch lg {
	case "c++":
		return "Cpp"
	case "c#":
		return "Csharp"
	case "javascript":
		return "JS"
	case "golang":
		return "Go"
	case "objective c":
		return "Obj-C"
	default:
		return langLowerCaseIndex[lg]
	}
}

func precomputeAutocompletions() map[string][]string {
	m := make(map[string][]string, 100)

	// Crazy shadowing of variable "lg" is allowed in go...
	for _, trueLg := range AllLanguages() {
		niceLg := PrintNiceLang(trueLg)
		for _, lg := range []string{trueLg, niceLg} {
			lg := strings.ToLower(lg)
			fragments := substrings(lg)
			for _, frag := range fragments {
				if !StringSliceContains(m[frag], niceLg) {
					m[frag] = append(m[frag], niceLg)
				}
			}
		}
	}

	for lg, trueLg := range synonymLangs {
		niceLg := PrintNiceLang(trueLg)
		lg := strings.ToLower(lg)
		fragments := substrings(lg)
		for _, frag := range fragments {
			if !StringSliceContains(m[frag], niceLg) {
				m[frag] = append(m[frag], niceLg)
			}
		}
	}

	fmt.Fprintf(os.Stderr, "---\n")
	return m
}

func substrings(s string) []string {
	L := len(s)
	seen := make(map[string]bool, L*L)
	fragments := make([]string, L*L)
	// This works well for language names with only 1-byte chars, not for any string
	for i := 0; i < L; i++ {
		for j := i + 1; j <= L; j++ {
			frag := s[i:j]
			if seen[frag] {
				continue
			}
			seen[frag] = true
			fragments = append(fragments, frag)
		}
	}
	return fragments
}

var langsExtraKeywords = map[string][]string{
	"Clojure": []string{"clj", "cljs", "cljc", "edn"},
	"Csharp":  []string{"cs"},
	"D":       []string{"dlang"},
	"Elixir":  []string{"ex", "exs"},
	"Go":      []string{"golang"},
	"JS":      []string{"javascript"},
	"Obj-C":   []string{"Objective", "Objective-C", "mm"},
	"Python":  []string{"py"},
	"Ruby":    []string{"rb"},
	"VB":      []string{"visual", "basic"},
	"Erlang":  []string{"erl", "hrl"},
	"Fortran": []string{"for", "f90", "f95"},
	"Haskell": []string{"hs", "lhs"},
	"Pascal":  []string{"pp", "pas", "inc"},
	"Perl":    []string{"pl"},
	"Rust":    []string{"rs"},
	"Scala":   []string{"sc"},
	"Scheme":  []string{"scs", "ss"},
}
