package main

import (
	"fmt"
	"os"

	"github.com/mpraes/tabyte/internal/runtime"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "serve" {
		fmt.Println("usage: tabyte serve")
		os.Exit(1)
	}
	if err := runtime.Serve("127.0.0.1:8787"); err != nil {
		fmt.Fprintf(os.Stderr, "serve error: %v\n", err)
		os.Exit(1)
	}
}