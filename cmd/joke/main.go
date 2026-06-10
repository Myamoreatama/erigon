// joke is a small command-line tool that fetches a random joke from a public
// API and prints it to stdout.
//
// Usage:
//
//	joke [-api <url>]
//
// Flags:
//
//	-api   Override the joke API URL (default: https://official-joke-api.appspot.com/random_joke)
package main

import (
	"flag"
	"fmt"
	"os"
)

var apiURL = flag.String("api", "", "joke API URL (leave empty to use the default endpoint)")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-api <url>]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Fetches a random joke from a public API and prints it.")
		fmt.Fprintln(os.Stderr, "\nFlags:")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	fetcher := NewFetcher(*apiURL)

	joke, err := fetcher.FetchJoke()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(joke)
}
