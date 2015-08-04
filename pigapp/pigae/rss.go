package pigae

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"text/template"

	. "github.com/Deleplace/programming-idioms/pig"

	"appengine"
)

// RssItem is a news.
type RssItem struct {
	XMLName     xml.Name `xml:"item"`
	Link        string   `xml:"link"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	GUID        GUID     `xml:"guid"`
}

// GUID wraps a URI that is probably not a Permalink.
type GUID struct {
	Value       string `xml:",innerxml"`
	IsPermaLink bool   `xml:"isPermaLink,attr"`
}

const rssTemplateString = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
<channel>
<title>{{.FeedTitle}}</title>
<link>{{.SiteLink}}</link>
<description>{{.FeedDescription}}</description>
<atom:link href="{{.FeedURL}}" rel="self" type="application/rss+xml" />
{{range .Items}}
{{.}}
{{end}}
</channel>
</rss>`

const rssPubDatelayout = "Mon, 2 Jan 2006 15:04:05 GMT"

var rssTemplate, _ = template.New("rss").Parse(rssTemplateString)

// RssFacade is the Facade for the RSS news feeds.
type RssFacade struct {
	FeedTitle       string
	SiteLink        string
	FeedDescription string
	Items           []string
	FeedURL         string
}

func rssRecentlyUpdated(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	_, idioms, _ := dao.getAllIdioms(c, 10, "-VersionDate")
	dateUpdate := func(idiom *Idiom) string { return idiom.VersionDate.Format(rssPubDatelayout) }
	idiomVersionGuidation := func(idiom *Idiom) string {
		return fmt.Sprintf("%v/guid/idiom/%v/version/%v", env.Host, idiom.Id, idiom.Version)
	}
	return rss(w, c, r, idioms, dateUpdate, idiomVersionGuidation, "/rss-recently-updated", "Programming Idioms recently updated idioms", "Idioms recently modified or having new implementations", "<br/><br/>Last updated in ")
}

func rssRecentlyCreated(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	_, idioms, _ := dao.getAllIdioms(c, 10, "-Id")
	dateCreation := func(idiom *Idiom) string { return idiom.CreationDate.Format(rssPubDatelayout) }
	idiomGuidation := func(idiom *Idiom) string { return fmt.Sprintf("%v/guid/idiom/%v", env.Host, idiom.Id) }
	return rss(w, c, r, idioms, dateCreation, idiomGuidation, "/rss-recently-created", "Programming Idioms recently created idioms", "Idioms recently created", "<br/><br/>Implemented in ")
}

func rss(w http.ResponseWriter, c appengine.Context, r *http.Request, idioms []*Idiom, datation func(*Idiom) string, guidation func(*Idiom) string, path string, feedTitle string, feedDescription string, listIntro string) error {

	itemsAsStrings := make([]string, len(idioms))
	for i, idiom := range idioms {
		desc := idiom.LeadParagraph
		sort.Sort(&implByVersionDateSorter{idiom.Implementations})
		if len(idiom.Implementations) > 0 {
			desc += listIntro + printNiceLang(idiom.Implementations[0].LanguageName)
			for _, impl := range idiom.Implementations[1:] {
				desc += ", " + impl.LanguageName
			}
			desc += "."
		}
		item := &RssItem{
			Link:        env.Host + NiceIdiomRelativeURL(idiom),
			Title:       idiom.Title,
			Description: desc,
			PubDate:     datation(idiom),
			GUID: GUID{
				Value:       guidation(idiom),
				IsPermaLink: false,
			},
		}
		buff, err := xml.MarshalIndent(item, "  ", "    ")
		if err != nil {
			return err
		}
		itemsAsStrings[i] = string(buff)
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	data := &RssFacade{
		FeedTitle:       feedTitle,
		SiteLink:        env.Host,
		FeedDescription: feedDescription,
		Items:           itemsAsStrings,
		FeedURL:         env.Host + path,
	}
	return rssTemplate.ExecuteTemplate(w, "rss", data)
}
