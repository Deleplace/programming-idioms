package pigae

import (
	"html/template"
	"reflect"
	"strings"

	. "github.com/Deleplace/programming-idioms/pig"
)

var templates, _ = initTemplates()

func initTemplates() (*template.Template, error) {
	var err error
	t := template.New("programming-idioms")

	funcMap := template.FuncMap{
		"mainStreamLanguages":   mainStreamLanguages,
		"moreLanguages":         moreLanguages,
		"allToggleNames":        allToggleNames,
		"printNiceLang":         printNiceLang,
		"printNiceLangs":        printNiceLangs,
		"prettifyCSSClass":      prettifyCSSClass,
		"normLang":              normLang,
		"langBadgeClass":        langBadgeClass,
		"isInStringList":        isInStringList,
		"idEqual":               idEqual,
		"decorate":              decorate,
		"langCoverageClass":     langCoverageClass,
		"NiceIdiomURL":          NiceIdiomURL,
		"NiceIdiomIDTitleURL":   NiceIdiomIDTitleURL,
		"NiceImplURL":           NiceImplURL,
		"themeDir":              themeDirectory,
		"hostPrefix":            hostPrefix,
		"host":                  host,
		"filterOut":             FilterOut,
		"toggled":               isToggled,
		"prettifyExtension":     prettifyExtension,
		"join":                  strings.Join,
		"dict":                  dict,
		"hiddenizeExtraColumns": hiddenizeExtraColumns,
		"empty":                 empty,
		"notEmpty":              notEmpty,
	}
	t = t.Funcs(funcMap)
	folders := []string{
		"template",
		"template/page",
		"template/page/admin",
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

// Workaround for proper PageMeta.Toggles
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
