package main

import (
	"fmt"
	"net/http"
	"strings"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func lookForNickname(r *http.Request) string {
	if cookie, errkie := r.Cookie("Nickname"); errkie == nil {
		return cookie.Value
	}
	return ""
}

func lookForFavoriteLanguages(r *http.Request) []string {
	if cookie, errkie := r.Cookie("my-languages"); errkie == nil {
		langs := strings.Split(cookie.Value, "_")
		langs = RemoveEmptyStrings(langs)
		langs = MapStrings(langs, NormLang)
		return langs
	}
	return nil
}

func seeNonFavorite(r *http.Request) bool {
	if cookie, errkie := r.Cookie("my-languages"); errkie != nil || cookie.Value == "" {
		// No favorite langs? Then you really need to see the other langs
		return true
	}
	if cookie, errkie := r.Cookie("see-non-favorite"); errkie == nil {
		if cookie.Value == "0" {
			return false
		}
	}
	// By default one should see all languages
	return true
}

func readUserProfile(r *http.Request) UserProfile {
	u := UserProfile{
		Nickname:          lookForNickname(r),
		FavoriteLanguages: lookForFavoriteLanguages(r),
		SeeNonFavorite:    seeNonFavorite(r),
		IsAdmin:           IsAdmin(r),
	}
	if u.Nickname != "" || len(u.FavoriteLanguages) > 0 {
		c := r.Context()
		log.Infof(c, "%v", u)
	}
	return u
}

func mustUserProfile(r *http.Request, w http.ResponseWriter) (UserProfile, error) {
	profile := readUserProfile(r)
	if profile.Nickname == "" {
		return profile, PiError{"You must already have a nickname.", http.StatusBadRequest}
	}
	return profile, nil
}

func setNicknameCookie(w http.ResponseWriter, nickname string) http.Cookie {
	newCookie := http.Cookie{
		Name:  "Nickname",
		Value: nickname,
		Path:  "/",
	}
	http.SetCookie(w, &newCookie)
	return newCookie
}

// langs should be underscore-separated, and end with an underscore
func setLanguagesCookie(w http.ResponseWriter, langs string) http.Cookie {
	newCookie := http.Cookie{
		Name:  "my-languages",
		Value: langs,
		Path:  "/",
	}
	http.SetCookie(w, &newCookie)
	return newCookie
}

//
// This URL will display homepage, and set soft profile cookies.
// That way users may transfer preferences to another browser,
// by emailing themselves or otherwise copy-pasting the URL.
//
// "nickname" is optional.
//
// "langs" parameter must be an underscore-separated list.
//
func bookmarkableUserURL(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userProfile := readUserProfile(r)
	c := r.Context()

	// Todo : encode/decode nickname with special chars
	nickname := vars["nickname"]
	langs := vars["langs"]
	if langs[0] == '_' {
		langs = langs[1:]
	}

	if nickname != "" {
		setNicknameCookie(w, nickname)
		userProfile.Nickname = nickname
	}
	setLanguagesCookie(w, langs)

	langsArray := strings.Split(langs, "_")
	if langsArray[len(langsArray)-1] == "" {
		langsArray = langsArray[:len(langsArray)-1]
	}
	langsArray = MapStrings(langsArray, NormLang)
	userProfile.FavoriteLanguages = langsArray

	// Display homepage, with updated profile
	return homeView(w, c, userProfile)
}

// Hard profiles?
//
// Not used yet.
// TODO To be adapted to : Handle optional user strong auth
func handleAuth(w http.ResponseWriter, r *http.Request) error {
	// Cf https://developers.google.com/appengine/docs/go/users/
	c := r.Context()
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, "/")
		if err != nil {
			return err
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return nil
	}
	fmt.Fprintf(w, "Hello, %v!", u)
	return nil
}
