package main

import . "github.com/Deleplace/programming-idioms/pig"

// MakeQAStructuredData is a structured data constructor for Q&A-style page.
func MakeQAStructuredData(idiom *Idiom, selectedImplID int, selectedImplLang string) (qa QAStructuredData) {
	if selectedImplID <= 0 || selectedImplLang == "" {
		return
	}
	qa.Question = idiom.Title + ", in " + PrintNiceLang(selectedImplLang)
	qa.Text = deemphasize(idiom.LeadParagraph)
	qa.Author = idiom.Author
	qa.DateCreated = idiom.CreationDate.Format(iso8601)
	qa.ImageURL = idiom.ImageURL
	for _, impl := range idiom.Implementations {
		if impl.LanguageName == selectedImplLang {
			qa.Answers = append(qa.Answers, QAStructuredDataAnswer{
				Text:        impl.CodeBlock,
				Author:      impl.Author,
				DateCreated: impl.CreationDate.Format(iso8601),
				URL:         env.Host + NiceImplRelativeURL(idiom, impl.Id, impl.LanguageName),
			})
		}
	}

	return qa
}

// YYYY-MM-DD
const iso8601 = "2006-01-02"
