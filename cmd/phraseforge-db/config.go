package main

import (
	"flag"
	"fmt"
	"io"
)

const (
	defaultDataFile     = "data/phrases.json"
	defaultDatabaseFile = "data/phraseforge.db"
)

type config struct {
	dataFile     string
	databaseFile string
}

func parseConfig(args []string) (config, error) {
	flags := flag.NewFlagSet("phraseforge-db", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	dataFile := flags.String("data-file", defaultDataFile, "path to the categories JSON file")
	databaseFile := flags.String("database-file", defaultDatabaseFile, "path to the SQLite database file")
	if err := flags.Parse(args); err != nil {
		return config{}, fmt.Errorf("parse setup flags: %w", err)
	}
	if len(flags.Args()) > 0 {
		return config{}, fmt.Errorf("unexpected argument %q", flags.Args()[0])
	}

	return config{
		dataFile:     *dataFile,
		databaseFile: *databaseFile,
	}, nil
}
