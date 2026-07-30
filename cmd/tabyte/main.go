package main

import (
	"fmt"
	"os"

	"github.com/mpraes/tabyte/internal/runtime"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "serve" {
		fmt.Println("usage: tabyte serve [--no-open]")
		os.Exit(1)
	}

	openBrowser := true
	for _, a := range os.Args[2:] {
		if a == "--no-open" {
			openBrowser = false
		}
	}
	// also honor env for smoke/CI
	if os.Getenv("TABYTE_NO_OPEN") == "1" {
		openBrowser = false
	}

	if err := runtime.Serve("127.0.0.1:8787", openBrowser); err != nil {
		fmt.Fprintf(os.Stderr, "serve error: %v\n", err)
		os.Exit(1)
	}
}