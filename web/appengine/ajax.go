package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	. "github.com/Deleplace/programming-idioms/idioms"
)

// Response is a generic container suitable to be directly converted into a JSON HTTP response.
// See http://nesv.blogspot.fr/2012/09/super-easy-json-http-responses-in-go.html
type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func ajaxIdiomVote(w http.ResponseWriter, r *http.Request) error {
	profile, err := mustUserProfile(r, w)
	if err != nil {
		return err
	}

	idiomIDStr := r.FormValue("idiomId")
	upOrDown := r.FormValue("choice")
	var incr int
	if upOrDown == "up" {
		incr = 1
		if err := togglesMissing(w, r, "idiomVotingUp"); err != nil {
			return err
		}
	} else if upOrDown == "down" {
		incr = -1
		if err := togglesMissing(w, r, "idiomVotingDown"); err != nil {
			return err
		}
	} else {
		return PiErrorf(http.StatusBadRequest, "Vote choice should be up or down")
	}
	ctx := r.Context()
	idiomID := String2Int(idiomIDStr)

	vote := IdiomVoteLog{
		IdiomId: idiomID,
		IpHash:  Sha1hash(r.RemoteAddr),
		Value:   incr,
		Date:    time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	var newRating int
	var myVote int
	if newRating, myVote, err = daoVotes.idiomVote(ctx, vote, profile.Nickname); err != nil {
		// w.WriteHeader(500)
		// fmt.Fprint(w, Response{"success": false, "message": err.Error()})
		return err
	}
	fmt.Fprint(w, Response{"success": true, "rating": newRating, "myVote": myVote})
	return nil
}

func ajaxImplVote(w http.ResponseWriter, r *http.Request) error {
	profile, err := mustUserProfile(r, w)
	if err != nil {
		return err
	}

	implIDStr := r.FormValue("implId")
	upOrDown := r.FormValue("choice")
	var incr int
	if upOrDown == "up" {
		incr = 1
		if err := togglesMissing(w, r, "implVotingUp"); err != nil {
			return err
		}
	} else if upOrDown == "down" {
		incr = -1
		if err := togglesMissing(w, r, "implVotingDown"); err != nil {
			return err
		}
	} else {
		return PiErrorf(http.StatusBadRequest, "Vote choice should be up or down")
	}
	ctx := r.Context()
	implID := String2Int(implIDStr)

	vote := ImplVoteLog{
		ImplId: implID,
		IpHash: Sha1hash(r.RemoteAddr),
		Value:  incr,
		Date:   time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	var newRating int
	var myVote int
	if newRating, myVote, err = daoVotes.implVote(ctx, vote, profile.Nickname); err != nil {
		// w.WriteHeader(500)
		// fmt.Fprint(w, Response{"success": false, "message": err.Error()})
		return err
	}
	fmt.Fprint(w, Response{"success": true, "rating": newRating, "myVote": myVote})
	return nil
}

func demoSiteSuggest(lang string) string {
	suggestion := ""
	switch strings.ToLower(lang) {
	case "js":
		suggestion = "https://jsfiddle.net/nick/..."
	case "go":
		suggestion = "https://play.golang.org/p/..."
	}
	if rand.Intn(10) < 4 {
		return "https://gist.github.com/..."
	}
	return suggestion
}

func ajaxDemoSiteSuggest(w http.ResponseWriter, r *http.Request) error {
	lang := r.FormValue("lang")
	suggestion := demoSiteSuggest(lang)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"suggestion": suggestion})
	return nil
}

func typeaheadLanguages(w http.ResponseWriter, r *http.Request) error {
	userInput := r.FormValue("userInput")
	suggestions := LanguageAutoComplete(userInput)
	w.Header().Set("Content-Type", "application/json")
	// TODO browser cache 2d
	// TODO server cache 2d
	// FIXME this prints {"options":null} for an empty result list, which is not the most frontend-friendly.
	fmt.Fprint(w, Response{"options": suggestions})
	return nil
}

func supportedLanguages(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"languages": AllNiceLangs})
	return nil
}
