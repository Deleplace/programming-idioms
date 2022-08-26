package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

const nbChanges = 50

func rssRecentChanges(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	_, changes, err := dao.getGlobalHistoryList(ctx, nbChanges)
	if err != nil {
		return err
	}
	path := "/rss-recent-changes"
	feedTitle := "Programming Idioms recent changes"
	feedDescription := "All recent edit actions on all idioms"

	itemsAsStrings := make([]string, len(changes))
	for i, change := range changes {
		title := "Change in [" + change.Idiom.Title + "]"
		if change.Idiom.Version == 1 {
			title = "Creation of [" + change.Idiom.Title + "]"
		}
		prev := change.Idiom.Version - 1 // always...?
		itemLink := fmt.Sprintf("%s/idiom/%d/diff/%d/%d", env.Host, change.Idiom.Id, prev, change.Idiom.Version)
		desc := "<br/>Edit: " + change.EditSummary +
			"<br/>Contributor: " + change.IdiomOrImplLastEditor + "."
		// TODO generate a short summary of the modified fields
		changeDate := change.VersionDate.Format(rssPubDatelayout)
		item := &RssItem{
			Link:        itemLink,
			Title:       markup2HTML(title),
			Description: markup2HTML(desc),
			PubDate:     changeDate,
			GUID: GUID{
				Value:       itemLink,
				IsPermaLink: true,
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
