package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	. "github.com/Deleplace/programming-idioms/idioms"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/user"
)

// IsAdmin determines whether the current user is regarded as Admin by the Google auth provider.
func IsAdmin(r *http.Request) bool {
	ctx := r.Context() // TODO check if NewContext is expensive
	u := user.Current(ctx)
	return u != nil && u.Admin
}

// AdminFacade is the Facade for the Admin page.
type AdminFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
}

func admin(w http.ResponseWriter, r *http.Request) error {
	data := &AdminFacade{
		PageMeta: PageMeta{
			ExtraJs: []string{hostPrefix() + themeDirectory() + "/js/programming-idioms-admin.js"},
			Toggles: toggles,
		},
	}

	return templates.ExecuteTemplate(w, "page-admin", data)
}

func ajaxRefreshToggles(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	err := dao.deleteCache(ctx)
	if err != nil {
		return err
	}
	return refreshToggles(ctx)
}

func ajaxSetToggle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	name := r.FormValue("toggle")
	valueAsString := r.FormValue("value")

	value, err := strconv.ParseBool(valueAsString)
	if err != nil {
		return err
	}
	log.Infof(ctx, "Setting toggle %q to %v", name, value)
	toggles[name] = value

	// Save config in distributed Datastore and Memcached
	err = dao.saveAppConfigProperty(ctx, AppConfigProperty{
		AppConfigId: 0, // TODO meaningful AppConfigId
		Name:        name,
		Value:       value,
	})
	if err != nil {
		return err
	}

	// Flush the cache, as some toggles have an big impact on page rendering.
	err = memcache.Flush(ctx)
	if err == nil {
		log.Infof(ctx, "Memcached flushed by toggle set")
	} else {
		log.Errorf(ctx, "flushing Memcached: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"success": true})
	return nil
}

// For related idioms (i.e. linked idioms)
func ajaxCreateRelation(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	idiomAIdStr := r.FormValue("idiomAId")
	idiomAId := String2Int(idiomAIdStr)
	idiomBIdStr := r.FormValue("idiomBId")
	idiomBId := String2Int(idiomBIdStr)

	keyA, idiomA, err := dao.getIdiom(ctx, idiomAId)
	if err != nil {
		return PiErrorf(http.StatusNotFound, "%v", err)
	}

	keyB, idiomB, err := dao.getIdiom(ctx, idiomBId)
	if err != nil {
		return PiErrorf(http.StatusNotFound, "%v", err)
	}

	idiomA.AddRelation(idiomB)
	if err := dao.saveExistingIdiom(ctx, keyA, idiomA); err != nil {
		return PiErrorf(http.StatusNotFound, "%v", err)
	}
	if err := dao.saveExistingIdiom(ctx, keyB, idiomB); err != nil {
		return PiErrorf(http.StatusNotFound, "%v", err)
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func sendMessageForUserAjax(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	msg := MessageForUser{
		Username:     r.FormValue("username"),
		Message:      r.FormValue("message"),
		CreationDate: time.Now(),
	}
	log.Infof(ctx, "Saving message for user [%v]: [%v].", msg.Username, Flatten(Shorten(msg.Message, 30)))
	_, err := dao.saveNewMessage(ctx, &msg)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func ajaxAdminMemcacheFlush(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	err := memcache.Flush(ctx)
	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		fmt.Fprint(w, Response{
			"success": true,
			"message": "Memcache flushed :)",
		})
	}
	log.Infof(ctx, "Memcached flushed by admin")
	return err
}
