package main

import (
	"context"
	"fmt"
	"reflect"

	stdlog "log"

	"google.golang.org/appengine"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
)

// 2021-12: delay funcs are no longer executed in local dev env with dev_appserver.py
//
// Workaround: IF local dev env THEN execute immediately (not in a delayed task)

type callable interface {
	Call(c context.Context, args ...interface{}) error
	// if needed:
	// Task(args ...interface{}) (*taskqueue.Task, error)
}

func delayFunc(name string, f interface{}) callable {
	if appengine.IsDevAppServer() {
		stdlog.Printf("Registering not-delayed func: %q\n", name)
		return notDelayedFunc{name: name, f: f}
	} else {
		return delay.Func(name, f)
	}
}

type notDelayedFunc struct {
	name string
	f    interface{}
}

func (ndf notDelayedFunc) Call(c context.Context, args ...interface{}) error {
	log.Infof(c, "Executing immediatly (not delayed): %q", ndf.name)

	fv := reflect.ValueOf(ndf.f)
	t := fv.Type()
	if t.Kind() != reflect.Func {
		return fmt.Errorf("%q: type %T kind %v, not a proper function", ndf.name, ndf.f, t.Kind())
	}

	argsv := make([]reflect.Value, len(args))
	for i, arg := range args {
		argsv[i] = reflect.ValueOf(arg)
	}
	in := append(
		[]reflect.Value{reflect.ValueOf(c)},
		argsv...,
	)

	resultv := fv.Call(in)
	result0 := resultv[0].Interface()
	if result0 == nil {
		return nil
	}
	return result0.(error)
}
