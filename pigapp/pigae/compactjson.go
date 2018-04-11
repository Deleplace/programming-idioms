package pigae

import (
	. "github.com/Deleplace/programming-idioms/pig"
)

// Sometimes we don't want to marshall all of the idioms data,
// we just want a relevant subset.
//
// This is especially useful when marshalling the whole list of
// idioms and implementations, for PWA purpose.

// See pistruct.go for field documentation.
type CompactIdiom struct {
	Id              int
	Title           string
	LeadParagraph   string
	ExtraKeywords   string `json:",omitempty"`
	ImageURL        string `json:",omitempty"`
	Version         int
	Implementations []CompactImpl
	RelatedIdiomIds []int `json:",omitempty"`
}

// See pistruct.go for field documentation.
type CompactImpl struct {
	Id                     int
	LanguageName           string
	ImportsBlock           string `json:",omitempty"`
	CodeBlock              string
	OriginalAttributionURL string `json:",omitempty"`
	DemoURL                string `json:",omitempty"`
	DocumentationURL       string `json:",omitempty"`
	AuthorComment          string `json:",omitempty"`
	PictureURL             string `json:",omitempty"`
}

func compactIdiom(idiom *Idiom) CompactIdiom {
	cidiom := CompactIdiom{
		Id:              idiom.Id,
		Title:           idiom.Title,
		LeadParagraph:   idiom.LeadParagraph,
		ExtraKeywords:   idiom.ExtraKeywords,
		ImageURL:        idiom.ImageURL,
		Version:         idiom.Version,
		RelatedIdiomIds: idiom.RelatedIdiomIds,
	}
	cidiom.Implementations = make([]CompactImpl, len(idiom.Implementations))
	for i, impl := range idiom.Implementations {
		cidiom.Implementations[i] = CompactImpl{
			Id:                     impl.Id,
			LanguageName:           impl.LanguageName,
			ImportsBlock:           impl.ImportsBlock,
			CodeBlock:              impl.CodeBlock,
			OriginalAttributionURL: impl.OriginalAttributionURL,
			DemoURL:                impl.DemoURL,
			DocumentationURL:       impl.DocumentationURL,
			AuthorComment:          impl.AuthorComment,
			PictureURL:             impl.PictureURL,
		}
	}
	return cidiom
}
