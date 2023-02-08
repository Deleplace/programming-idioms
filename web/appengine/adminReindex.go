package main

import (
	"fmt"
	"net/http"

	. "github.com/Deleplace/programming-idioms/idioms"

	"context"

	"google.golang.org/appengine/v2/datastore"
)

func adminReindexAjax(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	err := dao.deleteCache(ctx)
	if err != nil {
		errf(ctx, "Problem deleting cache: %v", err.Error())
	}
	err = dao.unindexAll(ctx)
	if err != nil {
		errf(ctx, "Problem deleting cache: %v", err.Error())
	}

	err = reindexDelayer.Call(ctx, "")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"message": "Reindexing launched in delayed tasks"})
	return nil
}

// Number of idioms being process by each single delayed task
const reindexBatchSize = 5

var reindexDelayer callable

func init() {
	reindexDelayer = delayFunc("reindex-idioms", func(ctx context.Context, cursorStr string) error {
		q := datastore.NewQuery("Idiom")
		if cursorStr != "" {
			logf(ctx, "Starting at cursor %v", cursorStr)
			cursor, err := datastore.DecodeCursor(cursorStr)
			if err != nil {
				return err
			}
			q = q.Start(cursor)
		}
		iterator := q.Run(ctx)

		reindexedIDs := make([]int, 0, reindexBatchSize)
		defer func() {
			logf(ctx, "Reindexed idioms %v", reindexedIDs)
		}()

		for i := 0; i < reindexBatchSize; i++ {
			var idiom Idiom
			key, err := iterator.Next(&idiom)
			if err == datastore.Done {
				logf(ctx, "Reindexing completed.")
				return nil
			} else if err != nil {
				// ouch :(
				return err
			}

			err = indexIdiomFullText(ctx, &idiom, key)
			if err != nil {
				errf(ctx, "Reindexing full text idiom %d : %v", idiom.Id, err)
			}
			err = indexIdiomCheatsheets(ctx, &idiom)
			if err != nil {
				errf(ctx, "Reindexing cheatsheet of idiom %d : %v", idiom.Id, err)
			}

			reindexedIDs = append(reindexedIDs, idiom.Id)
		}

		cursor, err := iterator.Cursor()
		if err != nil {
			// ouch :(
			return err
		}
		logf(ctx, "Stopping at cursor %v", cursor.String())
		reindexDelayer.Call(ctx, cursor.String())
		return nil
	})
}
