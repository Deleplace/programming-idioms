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

var mainStreamLangs = [...]string{"C", "Cpp", "Csharp", "Go", "Java", "JS", "Obj-C", "PHP", "Python", "Ruby", "Rust"}

// Return alpha codes for each language (no encoding problems).
// See PrintNiceLang to display them more fancy.
func MainStreamLanguages() []string {
	return mainStreamLangs[:]
}

var moreLangs = [...]string{"Ada", "Caml", "Clojure", "Cobol", "D", "Dart", "Elixir", "Erlang", "Fortran", "Groovy", "Haskell", "Kotlin", "Lua", "Lisp", "Pascal", "Perl", "Prolog", "Scala", "Scheme", "Smalltalk", "VB"}

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
		// Warning: this lazy init is a data race. Consider sync.Once instead.
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
	// Defensive copy (see issue #164)
	return CloneStringSlice(AllLangs)
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
	case "csharp", "cs":
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
	case "c++", "cc":
		return "Cpp"
	case "c#", "cs":
		return "Csharp"
	case "javascript":
		return "JS"
	case "golang":
		return "Go"
	case "objective c":
		return "Obj-C"
	case "py":
		return "Python"
	case "rs":
		return "Rust"
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

	for lg, words := range langsExtraKeywords {
		for _, word := range words {
			fragments := substrings(word)
			for _, frag := range fragments {
				if !StringSliceContains(m[frag], lg) {
					m[frag] = append(m[frag], lg)
				}
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
	"Erlang":  []string{"erl", "hrl"},
	"Fortran": []string{"for", "f90", "f95"},
	"Go":      []string{"golang"},
	"Haskell": []string{"hs", "lhs"},
	"JS":      []string{"javascript"},
	"Obj-C":   []string{"Objective", "Objective-C", "mm"},
	"Pascal":  []string{"pp", "pas", "inc", "turbopascal"},
	"Perl":    []string{"pl"},
	"Python":  []string{"py"},
	"Ruby":    []string{"rb"},
	"Rust":    []string{"rs"},
	"Scala":   []string{"sc"},
	"Scheme":  []string{"scs", "ss"},
	"VB":      []string{"visual", "basic", "vba", "vb6"},
}

func LanguageExtraKeywords(lg string) []string {
	// Careful, no defensive copy here!
	return langsExtraKeywords[lg]
}
