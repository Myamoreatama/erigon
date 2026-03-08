// Package main implements a random joke generator that fetches jokes from a
// public joke API.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// DefaultJokeAPIURL is the default endpoint for fetching random jokes.
	DefaultJokeAPIURL = "https://official-joke-api.appspot.com/random_joke"

	// defaultTimeout is the HTTP client request timeout.
	defaultTimeout = 10 * time.Second
)

// Joke represents a single joke returned by the API.
type Joke struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

// String returns a human-readable representation of the joke.
func (j Joke) String() string {
	return fmt.Sprintf("%s\n\n  ... %s", j.Setup, j.Punchline)
}

// Fetcher retrieves random jokes from an HTTP endpoint.
type Fetcher struct {
	apiURL string
	client *http.Client
}

// NewFetcher creates a Fetcher that uses the given API URL.
// If apiURL is empty, DefaultJokeAPIURL is used.
func NewFetcher(apiURL string) *Fetcher {
	if apiURL == "" {
		apiURL = DefaultJokeAPIURL
	}
	return &Fetcher{
		apiURL: apiURL,
		client: &http.Client{Timeout: defaultTimeout},
	}
}

// FetchJoke calls the joke API and returns a Joke or an error.
func (f *Fetcher) FetchJoke() (Joke, error) {
	resp, err := f.client.Get(f.apiURL) //nolint:noctx
	if err != nil {
		return Joke{}, fmt.Errorf("fetching joke: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Joke{}, fmt.Errorf("joke API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Joke{}, fmt.Errorf("reading response body: %w", err)
	}

	var joke Joke
	if err := json.Unmarshal(body, &joke); err != nil {
		return Joke{}, fmt.Errorf("parsing joke response: %w", err)
	}

	if joke.Setup == "" || joke.Punchline == "" {
		return Joke{}, fmt.Errorf("joke API returned incomplete joke (id=%d)", joke.ID)
	}

	return joke, nil
}
