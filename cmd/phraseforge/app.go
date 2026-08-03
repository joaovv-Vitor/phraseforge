package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/joaovv-Vitor/phraseforge/internal/phrase"
)

func run(args []string, out io.Writer, categories []phrase.Category) error {
	if len(args) == 0 {
		return fmt.Errorf("command is required: use generate or categories")
	}

	switch args[0] {
	case "generate":
		return runGenerate(args[1:], out, categories)
	case "categories":
		for _, name := range phrase.CategoryNames(categories) {
			if _, err := fmt.Fprintln(out, name); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown command %q: use generate or categories", args[0])
	}
}

func runGenerate(args []string, out io.Writer, categories []phrase.Category) error {
	flags := flag.NewFlagSet("generate", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	categoryName := flags.String("category", "programming", "category used to generate phrases")
	count := flags.Int("count", 1, "number of phrases to generate")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("parse generate flags: %w", err)
	}
	if len(flags.Args()) > 0 {
		return fmt.Errorf("generate command: unexpected argument %q", flags.Args()[0])
	}
	if *count < 1 {
		return fmt.Errorf("generate command: count must be greater than zero")
	}

	category, err := phrase.FindCategory(categories, *categoryName)
	if err != nil {
		return fmt.Errorf("run generate command: %w", err)
	}

	generatedPhrases := make([]string, 0, *count)
	for range *count {
		generated, err := phrase.Generate(category.Parts)
		if err != nil {
			return fmt.Errorf("run generate command: %w", err)
		}
		generatedPhrases = append(generatedPhrases, generated)
	}

	for index, generated := range generatedPhrases {
		if *count == 1 {
			_, err = fmt.Fprintln(out, generated)
		} else {
			_, err = fmt.Fprintf(out, "%d. %s\n", index+1, generated)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
