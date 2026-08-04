package main

import (
	"flag"
	"fmt"
	"io"
)

const defaultDataFile = "data/phrases.json"

type config struct {
	dataFile string
}

func parseConfig(args []string) (config, []string, error) {
	flags := flag.NewFlagSet("phraseforge", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	dataFile := flags.String("data-file", defaultDataFile, "path to the categories JSON file")
	if err := flags.Parse(args); err != nil {
		return config{}, nil, fmt.Errorf("parse application flags: %w", err)
	}

	return config{dataFile: *dataFile}, flags.Args(), nil
}
