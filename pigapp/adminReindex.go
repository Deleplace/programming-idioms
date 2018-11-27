package main

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/pig"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/delay"
)

func adminReindexAjax(w http.ResponseWriter, r *http.Request) error {
	c := r.Context()
	err := dao.deleteCache(c)
	if err != nil {
		warningf(c, "Problem deleting cache: %v", err.Error())
	}
	err = dao.unindexAll(c)
	if err != nil {
		warningf(c, "Problem deleting cache: %v", err.Error())
	}

	err = reindexDelayer.Call(c, "")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"message": "Reindexing launched in delayed tasks"})
	return nil
}

// Number of idioms being process by each single delayed task
const reindexBatchSize = 5

var reindexDelayer *delay.Function

func init() {
	reindexDelayer = delay.Func("reindex-idioms", func(c context.Context, cursorStr string) error {
		q := datastore.NewQuery("Idiom")
		if cursorStr != "" {
			infof(c, "Starting at cursor %v", cursorStr)
			cursor, err := datastore.DecodeCursor(cursorStr)
			if err != nil {
				return err
			}
			q = q.Start(cursor)
		}
		iterator := q.Run(c)

		reindexedIDs := make([]int, 0, reindexBatchSize)
		defer func() {
			infof(c, "Reindexed idioms %v", reindexedIDs)
		}()

		for i := 0; i < reindexBatchSize; i++ {
			var idiom Idiom
			key, err := iterator.Next(&idiom)
			if err == datastore.Done {
				infof(c, "Reindexing completed.")
				return nil
			} else if err != nil {
				// ouch :(
				return err
			}

			err = indexIdiomFullText(c, &idiom, key)
			if err != nil {
				errorf(c, "Reindexing full text idiom %d : %v", idiom.Id, err)
			}
			err = indexIdiomCheatsheets(c, &idiom)
			if err != nil {
				errorf(c, "Reindexing cheatsheet of idiom %d : %v", idiom.Id, err)
			}

			reindexedIDs = append(reindexedIDs, idiom.Id)
		}

		cursor, err := iterator.Cursor()
		if err != nil {
			// ouch :(
			return err
		}
		infof(c, "Stopping at cursor %v", cursor.String())
		reindexDelayer.Call(c, cursor.String())
		return nil
	})
}
