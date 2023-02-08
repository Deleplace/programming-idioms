package main

import (
	"context"
	"fmt"
	"os"
)

func logf(ctx context.Context, message string, params ...interface{}) {
	fmt.Printf(message+"\n", params...)
}

func errf(ctx context.Context, message string, params ...interface{}) {
	fmt.Fprintf(os.Stderr, message+"\n", params...)
}
