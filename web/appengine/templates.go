package main

import (
	"fmt"
	"html/template"
	"reflect"
	"strings"

	. "github.com/Deleplace/programming-idioms/idioms"
)

var templates, _ = initTemplates()

func initTemplates() (*template.Template, error) {
	var err error
	t := template.New("programming-idioms")

	funcMap := template.FuncMap{
		"markup2CSS":            markup2CSS,
		"mainStreamLanguages":   MainStreamLanguages,
		"moreLanguages":         MoreLanguages,
		"allToggleNames":        allToggleNames,
		"printNiceLang":         PrintNiceLang,
		"printNiceLangs":        PrintNiceLangs,
		"prettifyCSSClass":      prettifyCSSClass,
		"prettifyExtension":     prettifyExtension,
		"normLang":              NormLang,
		"langBadgeClass":        langBadgeClass,
		"isInStringList":        isInStringList,
		"idEqual":               idEqual,
		"decorate":              decorate,
		"langCoverageClass":     langCoverageClass,
		"niceIdiomURL":          NiceIdiomURL,
		"niceIdiomIDTitleURL":   NiceIdiomIDTitleURL,
		"niceImplURL":           NiceImplURL,
		"themeDir":              themeDirectory,
		"hostPrefix":            hostPrefix,
		"host":                  host,
		"filterOut":             FilterOut,
		"toggled":               isToggled,
		"join":                  strings.Join,
		"dict":                  dict,
		"hiddenizeExtraColumns": hiddenizeExtraColumns,
		"empty":                 empty,
		"notEmpty":              notEmpty,
		"blank":                 isBlank,
		"implementationsFor":    implementationsFor,
		"shorten":               Shorten,
		"atom2string":           atom2string,
		"atom2int":              atom2int,
		"trim":                  strings.TrimSpace,
		"ifval":                 ifval,
		"diffClass":             diffClass,
		"plus":                  plus,
		"minus":                 minus,
		"hasSuffix":             strings.HasSuffix,
		"replace":               strings.Replace,
		"logo":                  languageLogo,
	}
	t = t.Funcs(funcMap)
	folders := []string{
		"template",
		"template/page",
		"template/page/about",
		"template/page/admin",
		"template/page/hybrid",
		"template/ajax-block",
		"template/content/widget",
		"template/content/block",
		"template/header",
		"template/footer",
	}
	for _, f := range folders {
		t, err = t.ParseGlob(f + "/*.html")
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func isInStringList(lang string, favlangs []string) bool {
	return StringSliceContains(favlangs, lang)
}

func idEqual(a int, b int) bool {
	return a == b
}

func empty(x interface{}) bool {
	v := reflect.ValueOf(x)
	return v.Len() == 0
}

func notEmpty(x interface{}) bool {
	v := reflect.ValueOf(x)
	return v.Len() > 0
}

func isBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Directory containing CSS, JS, and pictures.
// Warning : contains a leading slash
// Warning : does not contain a trailing slash
func themeDirectory() string {
	if toggles["themeVirtualVersioning"] {
		return "/" + ThemeVersion + "_" + ThemeDate
	}
	return "/" + ThemeVersion
}

func hostPrefix() string {
	if toggles["useAbsoluteUrls"] {
		return env.Host
	}
	return ""
}

func host() string {
	return env.Host
}

// Poor workaround for proper PageMeta.Toggles
//
// Note this is reading global toggles, it doesn't
// work at all with page custom toggles.
func isToggled(name string) bool {
	return toggles[name]
}

// If no favorites, all badges are blue
// If favorite, badge is green
// If there are favorites and this one is not, badge is grey
func langBadgeClass(lang string, favlangs []string) string {
	if len(favlangs) > 0 {
		if StringSliceContains(favlangs, lang) {
			return "badge-fav-lang"
		} else {
			return "badge-non-fav-lang"
		}
	} else {
		return "badge-lang"
	}
}

// If favorite, circle is green
func langCoverageClass(lang string, favlangs []string) string {
	if StringSliceContains(favlangs, lang) {
		return "coverage-fav-lang"
	}
	return ""
}

// A "functional if": depending of first (bool) argument,
// return second or third argument.
// Note that all arguments are always evaluated, even if discarded.
func ifval(cond bool, value1, value2 interface{}) interface{} {
	if cond {
		return value1
	}
	return value2
}

// diffClass compare arguments to determine some CSS class to apply:
// - empty when equals
// - "touched" when both values are different and non-empty
// - "created" when left is empty and right is non-empty
// - "deleted" when left is non-empty and right is empty
func diffClass(leftArg, rightArg interface{}) string {
	switch left := leftArg.(type) {
	case string:
		right, ok := rightArg.(string)
		if !ok {
			panic(fmt.Errorf("unexpected right type %T\n", rightArg))
		}

		// \r\n are pita
		left = strings.Replace(left, "\r\n", "\n", -1)
		right = strings.Replace(right, "\r\n", "\n", -1)

		if left == right {
			return ""
		}
		if left == "" {
			return "created"
		}
		if right == "" {
			return "deleted"
		}
		return "touched"
	// TODO other useful types?
	default:
		panic(fmt.Errorf("unexpected left type %T\n", leftArg))
	}
}

func plus(a, b int) int {
	return a + b
}

func minus(a, b int) int {
	return a - b
}
