# Testing Strategy — `go-anilist`

This document defines testing standards for `go-anilist`, ensuring deterministic,
fast, isolated, and reproducible test suites without external network access.

---

## 1. Core Testing Principles

1. **Zero Real Network Traffic**: No unit test may connect to `https://graphql.anilist.co`. All network traffic is mocked locally.
2. **Determinism and Speed**: The full suite must complete in less than 2 seconds.
3. **Concurrency Safety**: All tests must pass cleanly under the Go race detector (`go test -v -race ./...`).
4. **Bidirectional Verification**: Tests must verify both what the client *sends* (query, variables, HTTP headers) and what it *parses* (typed domain models, quota extraction, errors).

---

## 2. Local Test Server (`httptest.Server`)

Because AniList GraphQL dispatches all operations via HTTP `POST` to the root `/`,
the mock server inspects the JSON payload:

```go
type testMuxHandler struct {
    handlers map[string]http.HandlerFunc
}

func setupTestServer(t *testing.T) (*Client, *testMuxHandler, func()) {
    t.Helper()
    muxHandler := &testMuxHandler{
        handlers: make(map[string]http.HandlerFunc),
    }

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }
        if r.Header.Get("Content-Type") != "application/json" {
            http.Error(w, "bad content type", http.StatusBadRequest)
            return
        }

        var reqBody struct {
            Query         string                 `json:"query"`
            OperationName string                 `json:"operationName"`
            Variables     map[string]interface{} `json:"variables"`
        }
        if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
            http.Error(w, "invalid json body", http.StatusBadRequest)
            return
        }

        opKey := reqBody.OperationName
        if handler, ok := muxHandler.handlers[opKey]; ok {
            handler(w, r)
            return
        }

        http.Error(w, fmt.Sprintf("unhandled test operation: %s", opKey), http.StatusNotFound)
    }))

    client := NewClient(nil, WithBaseURL(server.URL))
    return client, muxHandler, server.Close
}
```

---

## 3. Fixtures and Test Data (`testdata/`)

Mock JSON responses are kept in `testdata/`:

```
testdata/
  viewer_success.json            ← Viewer query response
  media_search_success.json      ← Search query with media list and pageInfo
  medialist_get_success.json     ← Existing media list entry
  medialist_save_success.json    ← SaveMediaListEntry mutation response
```

Loaded in unit tests via `os.ReadFile("testdata/...")`.

---

## 4. Context Cancellation Testing

Each service must verify that a cancelled `context.Context` (`context.Canceled` or
`context.DeadlineExceeded`) immediately aborts the HTTP request without goroutine leaks.