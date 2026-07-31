package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mpraes/tabyte/internal/runtime"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "serve" {
		fmt.Println("usage: tabyte serve [--no-open] [--persist]")
		os.Exit(1)
	}

	openBrowser := true
	persist := false
	for _, a := range os.Args[2:] {
		switch a {
		case "--no-open":
			openBrowser = false
		case "--persist":
			persist = true
		}
	}
	if os.Getenv("TABYTE_NO_OPEN") == "1" {
		openBrowser = false
	}
	if v := strings.TrimSpace(os.Getenv("TABYTE_PERSIST")); v == "1" || strings.EqualFold(v, "true") {
		persist = true
	}

	if err := runtime.Serve(runtime.ServeOptions{
		Addr:        "127.0.0.1:8787",
		OpenBrowser: openBrowser,
		Persist:     persist,
		DBPath:      os.Getenv("TABYTE_DB_PATH"),
	}); err != nil {
		fmt.Fprintf(os.Stderr, "serve error: %v\n", err)
		os.Exit(1)
	}
}
