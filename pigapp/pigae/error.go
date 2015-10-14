package pigae

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"appengine"
)

// PiError is a custom error type, which embeds a HTTP error code.
type PiError struct {
	ErrorText string
	Code      int
}

func (p PiError) Error() string {
	return p.ErrorText
}

// ErrorFacade is the Facade for the Error page (BSOD).
type ErrorFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	Error       string
	ErrorCode   int
}

func errorPage(w http.ResponseWriter, r *http.Request, err error) {
	userProfile := readUserProfile(r)
	var text string
	var code int

	switch err.(type) {
	case PiError:
		pierr := err.(PiError)
		text = pierr.ErrorText
		code = pierr.Code
	default:
		text = err.Error()
		code = http.StatusInternalServerError
	}

	c := appengine.NewContext(r)
	c.Errorf(text)

	data := &ErrorFacade{
		PageMeta: PageMeta{
			PageTitle: "Oops",
			Toggles:   toggles,
		},
		UserProfile: userProfile,
		Error:       text,
		ErrorCode:   code,
	}

	w.WriteHeader(code)
	errt := templates.ExecuteTemplate(w, "page-error", data)
	if errt != nil {
		c.Errorf("Problem rendering error page: %v", errt.Error())
	}
}

func errorJSON(w http.ResponseWriter, r *http.Request, err error) {
	var text string
	var code int

	switch err.(type) {
	case PiError:
		pierr := err.(PiError)
		text = pierr.ErrorText
		code = pierr.Code
	default:
		text = err.Error()
		code = http.StatusInternalServerError
	}

	c := appengine.NewContext(r)
	c.Errorf(text)

	w.WriteHeader(code)
	fmt.Fprint(w, Response{
		"success": false,
		"message": text,
	})
}
