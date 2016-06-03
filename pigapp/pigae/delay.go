package pigae

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
)

var cursorRecursiveDelayer *delay.Function

func init() {
	cursorRecursiveDelayer = delay.Func("cursor-recurse", func(c context.Context, delayableName, cursorStr string) error {
		dq := delayables[delayableName]

		if cursorStr != "" {
			log.Infof(c, "Starting at cursor %v", cursorStr)
			startCursor, err := datastore.DecodeCursor(cursorStr)
			if err != nil {
				return err
			}
			dq.query = dq.query.Start(startCursor)
		}

		chunk := make([]interface{}, dq.chunkSize)
		for i := range chunk {
			chunk[i] = dq.entityFactory()
		}
		toSaveKeys := make([]*datastore.Key, 0, dq.chunkSize)
		toSaveEntities := make([]interface{}, 0, dq.chunkSize)
		finished := false

		iterator := dq.query.Run(c)
		for i := 0; i < dq.chunkSize; i++ {
			key, err := iterator.Next(&chunk[i])
			if err == datastore.Done {
				finished = true
				break
			} else if err != nil {
				// ouch :(
				return err
			}

			saveIt, err := dq.action(c, chunk[i])
			if err != nil {
				log.Errorf(c, "Action %v on entity %v : %v", dq.name, chunk[i], err)
				continue
				// TODO fail job or not?
			}
			if saveIt {
				toSaveKeys = append(toSaveKeys, key)
				toSaveEntities = append(toSaveEntities, chunk[i])
			}
		}

		if len(toSaveEntities) > 0 {
			_, err := datastore.PutMulti(c, toSaveKeys, toSaveEntities)
			if err != nil {
				log.Errorf(c, "Failed saving %d entities to the Datastore : %v", len(toSaveEntities), err)
				return err
			}
		}

		if finished {
			return nil
		}

		cursor, err := iterator.Cursor()
		if err != nil {
			// ouch :(
			return err
		}

		log.Infof(c, "Stopping at cursor %v", cursor.String())
		return cursorRecursiveDelayer.Call(c, delayableName, cursor.String())
	})
}

var delayables map[string]DelayableQuery

func DelayableList() []string {
	list := make([]string, 0, len(delayables))
	for _, dq := range delayables {
		list = append(list, dq.name)
	}
	return list
}

type Action func(context.Context, interface{}) (resave bool, err error)

type DelayableQuery struct {
	name          string
	query         *datastore.Query
	action        Action
	entityFactory func() interface{}
	chunkSize     int
}

func registerDelayable(name string, query *datastore.Query, action Action, entityFactory func() interface{}) DelayableQuery {
	dq := DelayableQuery{
		name:          name,
		query:         query,
		action:        action,
		entityFactory: entityFactory,
		chunkSize:     50,
	}
	delayables[name] = dq
	return dq
}
