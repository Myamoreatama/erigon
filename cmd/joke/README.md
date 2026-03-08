# joke — Random Joke Generator

`joke` is a tiny command-line tool bundled with this repository that fetches a
random joke from a public joke API and prints it to standard output.

## Usage

```
joke [-api <url>]
```

### Flags

| Flag   | Default                                                            | Description                              |
|--------|--------------------------------------------------------------------|------------------------------------------|
| `-api` | `https://official-joke-api.appspot.com/random_joke` | Override the joke API endpoint URL |

### Examples

Fetch a joke using the default API endpoint:

```bash
./build/bin/joke
```

Sample output:

```
Why did the Go developer quit?

  ... Because they didn't get arrays.
```

Override the API endpoint (useful for testing):

```bash
./build/bin/joke -api http://localhost:8080/joke
```

## Build

From the repository root:

```bash
cd cmd/joke && go build -o ../../build/bin/joke .
```

Or use the Makefile wildcard rule:

```bash
make joke.cmd
```

## Tests

```bash
go test ./cmd/joke/...
```

The tests use Go's built-in `net/http/httptest` package to mock the external
API — no network access is required.

## Notes

- HTTP requests time out after **10 seconds**.
- Non-200 responses from the API are treated as errors.
- The tool is entirely self-contained in `cmd/joke/` and does **not** modify
  any existing Erigon core packages.
