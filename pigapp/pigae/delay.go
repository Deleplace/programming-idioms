package pigae

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
)

type Action func(context.Context, interface{})

var cursorRecursiveDelayer *delay.Function

func init() {
	cursorRecursiveDelayer = delay.Func("cursor-recurse", func(c context.Context, query *datastore.Query, cursorStr string, f Action) error {
		log.Infof(c, "Starting at cursor %v", cursorStr)
		cursor, err := datastore.DecodeCursor(cursorStr)
		if err != nil {
			return err
		}
		query = query.Start(cursor)

		// TODO

		log.Infof(c, "Stopping at cursor %v", cursor.String())
		reindexDelayer.Call(c, query, cursor.String(), f)
		return nil
	})
}
