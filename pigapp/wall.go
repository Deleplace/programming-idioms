package main

import (
	"html/template"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"
)

// WallFacade is the Facade for the Wall page.
type WallFacade struct {
	PageMeta    PageMeta
	HtmlMessage template.HTML
}

// From a message, creates a handler
func makeWall(msg string) betterHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return wall(w, r, msg)
	}
}

func wall(w http.ResponseWriter, r *http.Request, htmlMessage string) error {
	customToggles := Toggles{}
	/* Anyway in the wall, ALL features are deactivated
		for k, v := range toggles {
	    	customToggles[k] = v
		}
		customToggles["writable"] = false
		customToggles["greetings"] = false
		customToggles["searchable"] = false
	*/
	data := WallFacade{
		PageMeta: PageMeta{
			PageTitle: "Programming-Idioms is currently not available",
			Toggles:   customToggles,
		},
		HtmlMessage: template.HTML(htmlMessage),
	}
	return templates.ExecuteTemplate(w, "page-wall", data)
}
