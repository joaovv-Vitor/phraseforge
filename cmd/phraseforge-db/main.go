package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	configuration, err := parseConfig(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := run(context.Background(), os.Stdout, configuration); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
