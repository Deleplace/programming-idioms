package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func userMessageBoxAjax(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userProfile := readUserProfile(r)
	username := userProfile.Nickname
	keys, messages, err := dao.getMessagesForUser(ctx, username)
	if err != nil {
		return err
	}

	// Transform (keys + entities) into json objects
	jsonMessages := make([]map[string]interface{}, len(messages))
	for i, key := range keys {
		message := messages[i]
		jsonMessage := make(map[string]interface{}, 3)
		jsonMessage["key"] = key
		jsonMessage["message"] = message.Message
		jsonMessage["creationDate"] = message.CreationDate
		jsonMessages[i] = jsonMessage
	}
	jsonResponse := Response{"messages": jsonMessages}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, jsonResponse)
	// TODO add a 3-mn browser cache header?
	return nil
}

func dismissUserMessage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	keyStr := r.FormValue("key")
	log.Infof(ctx, "Dismissing user message for key %v", keyStr)
	key, err := datastore.DecodeKey(keyStr)
	if err != nil {
		return err
	}
	_, err = dao.dismissMessage(ctx, key)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
