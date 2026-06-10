package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newMockServer starts a test HTTP server that responds with the given fixture.
func newMockServer(t *testing.T, statusCode int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if body != nil {
			if err := json.NewEncoder(w).Encode(body); err != nil {
				t.Errorf("mock server: encode response: %v", err)
			}
		}
	}))
}

func TestFetchJoke_Success(t *testing.T) {
	want := Joke{ID: 1, Type: "general", Setup: "Why did the Go developer quit?", Punchline: "Because they didn't get arrays."}

	srv := newMockServer(t, http.StatusOK, want)
	defer srv.Close()

	fetcher := NewFetcher(srv.URL)
	got, err := fetcher.FetchJoke()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestFetchJoke_HTTPError(t *testing.T) {
	srv := newMockServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	fetcher := NewFetcher(srv.URL)
	_, err := fetcher.FetchJoke()
	if err == nil {
		t.Fatal("expected an error for non-200 status, got nil")
	}
}

func TestFetchJoke_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("not-json"))
	}))
	defer srv.Close()

	fetcher := NewFetcher(srv.URL)
	_, err := fetcher.FetchJoke()
	if err == nil {
		t.Fatal("expected a JSON parse error, got nil")
	}
}

func TestFetchJoke_IncompleteJoke(t *testing.T) {
	// A valid JSON object but with empty setup/punchline.
	incomplete := Joke{ID: 5, Type: "general", Setup: "", Punchline: ""}

	srv := newMockServer(t, http.StatusOK, incomplete)
	defer srv.Close()

	fetcher := NewFetcher(srv.URL)
	_, err := fetcher.FetchJoke()
	if err == nil {
		t.Fatal("expected an error for incomplete joke, got nil")
	}
}

func TestNewFetcher_DefaultURL(t *testing.T) {
	f := NewFetcher("")
	if f.apiURL != DefaultJokeAPIURL {
		t.Errorf("expected default URL %q, got %q", DefaultJokeAPIURL, f.apiURL)
	}
}

func TestNewFetcher_CustomURL(t *testing.T) {
	custom := "http://example.com/joke"
	f := NewFetcher(custom)
	if f.apiURL != custom {
		t.Errorf("expected custom URL %q, got %q", custom, f.apiURL)
	}
}

func TestJoke_String(t *testing.T) {
	j := Joke{Setup: "Why did the chicken cross the road?", Punchline: "To get to the other side."}
	s := j.String()
	if s == "" {
		t.Error("String() returned empty string")
	}
}
