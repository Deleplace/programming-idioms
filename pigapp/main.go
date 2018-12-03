package main

import (
	"google.golang.org/appengine"
)

func main() {
	// "It is imperative to invoke flush before your main function exits"
	defer sd.Flush()

	appengine.Main()
}
