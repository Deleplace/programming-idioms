package main

import (
	"net/http"
	"strings"
	"time"

	. "github.com/Deleplace/programming-idioms/pig"
	"google.golang.org/appengine/log"

	"github.com/gorilla/mux"
)

// Language Backlog: community contribution nudge.
// Shows a few stats about the Coverage for this lang,
// and links to idioms/impls to be improved.

// 4 sections will show sampleSize lines each: Missing impl, Missing doc, Missing demo, Curation.
const sampleSize = 3

func backlogForLanguage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	vars := mux.Vars(r)
	rawLang := vars["lang"]
	lang := NormLang(rawLang)
	log.Infof(ctx, "Computing backlog for %s", lang)

	data := &BacklogLanguageFacade{
		PageMeta: PageMeta{
			PageTitle: "Idioms missing data for the " + lang + " implementation",
			Toggles:   toggles,
			ExtraCss:  []string{hostPrefix() + themeDirectory() + "/css/pages/backlog.css"},
			ExtraJs:   []string{hostPrefix() + themeDirectory() + "/js/pages/backlog.js"},
		},
		UserProfile:         readUserProfile(r),
		Lang:                lang,
		LanguageLogo:        languageLogo(lang),
		RecommendedDemoSite: recommendedDemoSite(rawLang),
	}

	log.Infof(ctx, "searchRandomImplsForLang(%q)...", rawLang)
	tip := time.Now()
	var err error
	data.CurationSuggestions, err = searchRandomImplsForLang(ctx, rawLang, sampleSize)
	log.Infof(ctx, "got %d curation suggestions for %s in %dms", len(data.CurationSuggestions), rawLang, time.Since(tip)/time.Millisecond)
	if err != nil {
		log.Errorf(ctx, "%v", err)
	}
	tip = time.Now()
	data.MissingDocDemo, err = searchMissingDocDemoForLang(ctx, rawLang, sampleSize)
	log.Infof(ctx, "got %d missingDoc and %d missingDemo for %s in %dms", len(data.MissingDocDemo.MissingDoc), len(data.MissingDocDemo.MissingDemo), rawLang, time.Since(tip)/time.Millisecond)
	if err != nil {
		log.Errorf(ctx, "%v", err)
	}

	tip = time.Now()
	data.MissingImpl, err = searchMissingImplForLang(ctx, lang, sampleSize)
	log.Infof(ctx, "got %d missingImpl idioms %s in %dms", len(data.MissingImpl.Stubs), rawLang, time.Since(tip)/time.Millisecond)
	if err != nil {
		log.Errorf(ctx, "%v", err)
	}

	log.Infof(ctx, "Done computing backlog for %s", lang)
	log.Infof(ctx, "Executing backlog template for %s", lang)
	tip = time.Now()
	err = templates.ExecuteTemplate(w, "page-backlog-language", data)
	log.Infof(ctx, "Done executing backlog template for %s in %dms", lang, time.Since(tip)/time.Millisecond)
	return err
}

// SearchResultsFacade is the facade for the Search Results page.
type BacklogLanguageFacade struct {
	PageMeta            PageMeta
	UserProfile         UserProfile
	Lang                string
	LanguageLogo        string
	RecommendedDemoSite DemoSite
	CurationSuggestions []*IdiomSingleton
	MissingDocDemo      backlogMissingDocDemo
	MissingImpl         backlogMissingImpl
}

// IdiomSingleton is an idiom that contains a single Implementation.
// It is intended for display logic only.
// It is NOT intended to be saved to the database.
type IdiomSingleton = Idiom

// IdiomStub is an idiom that contains zero Implementation.
// It is intended for display logic only.
// It is NOT intended to be saved to the database.
type IdiomStub = Idiom

type DemoSite struct {
	Name string
	URL  string
}

func recommendedDemoSite(lang string) DemoSite {
	var ds DemoSite
	switch strings.TrimSpace(strings.ToLower(lang)) {
	case "csharp", "cs", "fsharp", "fs":
		ds.Name = "SharpLab"
		ds.URL = "https://sharplab.io/"
	case "dart":
		ds.Name = "DartPad"
		ds.URL = "https://dartpad.dev/"
	case "go":
		ds.Name = "the Go Playground"
		ds.URL = "https://play.golang.org/"
	case "rust":
		ds.Name = "the Rust Playground"
		ds.URL = "https://play.rust-lang.org/"
	default:
		// No recommended demo site for this lang
	}
	return ds
}

// E.g. "python.svg"  for the file at static/default/img/logos/python.svg
func languageLogo(lang string) string {
	switch lg := strings.TrimSpace(strings.ToLower(lang)); lg {
	case "csharp", "cs":
		return "csharp.svg"
	case "c", "clojure", "dart", "elixir", "groovy", "haskell", "java", "kotlin", "lua", "php":
		return lg + ".svg"
	case "cpp", "c++":
		return "cpp.svg"
	case "d", "dlang":
		return "d.svg"
	case "go", "golang":
		return "go.svg"
	case "js", "javascript":
		return "js.svg"
	case "python", "py":
		return "python.svg"
	case "ruby", "rb":
		return "ruby.svg"
	case "rust", "rs":
		return "rust.svg"
	default:
		// No logo for this language
		return ""
	}
}

// Block "Curation"
func backlogBlockCuration(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	vars := mux.Vars(r)
	rawLang := vars["lang"]
	lang := NormLang(rawLang)
	log.Infof(ctx, "Computing backlog for %s", lang)

	data := &BacklogLanguageFacade{
		Lang: lang,
		// Mostly empty because we don't care about the full page metadata
	}

	log.Infof(ctx, "searchRandomImplsForLang(%q)...", rawLang)
	tip := time.Now()
	var err error
	data.CurationSuggestions, err = searchRandomImplsForLang(ctx, rawLang, sampleSize)
	log.Infof(ctx, "got %d curation suggestions for %s in %dms", len(data.CurationSuggestions), rawLang, time.Since(tip)/time.Millisecond)
	if err != nil {
		log.Errorf(ctx, "%v", err)
	}

	log.Infof(ctx, "Done computing backlog block Curation for %s", lang)
	log.Infof(ctx, "Executing backlog block Curation template for %s", lang)
	tip = time.Now()
	err = templates.ExecuteTemplate(w, "backlog-block-curation", data)
	log.Infof(ctx, "Done executing backlog block Curation template for %s in %dms", lang, time.Since(tip)/time.Millisecond)
	return err
}

//
// The Block handlers below are intended to allow refreshing only a portion of the Backlog page.
//

// Block "Docs & Demos"
func backlogBlockDocsDemos(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	vars := mux.Vars(r)
	rawLang := vars["lang"]
	lang := NormLang(rawLang)
	log.Infof(ctx, "Computing backlog for %s", lang)

	data := &BacklogLanguageFacade{
		Lang:                lang,
		RecommendedDemoSite: recommendedDemoSite(rawLang),
		// Mostly empty because we don't care about the full page metadata
	}

	tip := time.Now()
	var err error
	data.MissingDocDemo, err = searchMissingDocDemoForLang(ctx, rawLang, sampleSize)
	log.Infof(ctx, "got %d missingDoc and %d missingDemo for %s in %dms", len(data.MissingDocDemo.MissingDoc), len(data.MissingDocDemo.MissingDemo), rawLang, time.Since(tip)/time.Millisecond)
	if err != nil {
		log.Errorf(ctx, "%v", err)
	}

	log.Infof(ctx, "Done computing backlog block Docs-Demos for %s", lang)
	log.Infof(ctx, "Executing backlog block Docs-Demos templates for %s", lang)
	tip = time.Now()
	err = templates.ExecuteTemplate(w, "backlog-block-missing-doc", data)
	if err != nil {
		return err
	}
	err = templates.ExecuteTemplate(w, "backlog-block-missing-demo", data)
	log.Infof(ctx, "Done executing backlog block Docs-Demos templates for %s in %dms", lang, time.Since(tip)/time.Millisecond)
	return err
}

// Block "Missing implementations"
func backlogBlockMissingImpl(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	vars := mux.Vars(r)
	rawLang := vars["lang"]
	lang := NormLang(rawLang)
	log.Infof(ctx, "Computing backlog for %s", lang)

	data := &BacklogLanguageFacade{
		Lang: lang,
		// Mostly empty because we don't care about the full page metadata
	}

	log.Infof(ctx, "searchRandomImplsForLang(%q)...", rawLang)
	tip := time.Now()
	var err error
	data.MissingImpl, err = searchMissingImplForLang(ctx, lang, sampleSize)
	log.Infof(ctx, "got %d missingImpl idioms %s in %dms", len(data.MissingImpl.Stubs), rawLang, time.Since(tip)/time.Millisecond)
	if err != nil {
		log.Errorf(ctx, "%v", err)
	}

	log.Infof(ctx, "Done computing backlog block Missing-Impl for %s", lang)
	log.Infof(ctx, "Executing backlog block Missing-Impl template for %s", lang)
	tip = time.Now()
	err = templates.ExecuteTemplate(w, "backlog-block-missing-impl", data)
	log.Infof(ctx, "Done executing backlog block Missing-Impl template for %s in %dms", lang, time.Since(tip)/time.Millisecond)
	return err
}
