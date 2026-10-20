# Reference Analysis: `go-myanimelist`

This document analyzes the architecture, design patterns, and engineering decisions of
[`go-myanimelist`](https://github.com/nstratos/go-myanimelist) (author: Stratos Neiros),
used as the primary reference design for `go-anilist`.

---

## 1. Attribution and License

- **License**: MIT License **[DOC-GO]**.
  ```
  Copyright (c) 2015-present Stratos Neiros <nstratos@gmail.com>
  ```
- **Attribution**: The MIT license allows reuse of patterns, client structures, and testing
  methodologies provided that the copyright and permission notices are preserved.

---

## 2. Package Structure

`go-myanimelist` organizes its entire core within a single root package (`mal`):
- `mal.go`: Defines `Client`, `Response`, `ErrorResponse`, `Do()`, `NewRequest()`, and core helpers.
- `anime.go`, `manga.go`, `user.go`, `forum.go`: Services and domain models.
- `useranimelist.go`, `usermangalist.go`: Sub-resources for user tracking lists.
- `*_test.go`: Unit tests and integration suite.

**Evaluation**: A single root package prevents unnecessary fragmentation, import cycles,
and redundant package hierarchies for moderate-sized Go client libraries.

---

## 3. Client & Service Delegation Pattern

The client acts as the central coordinator and delegates to typed sub-services:

```go
type Client struct {
    client  *http.Client
    BaseURL *url.URL

    Anime *AnimeService
    Manga *MangaService
    User  *UserService
    Forum *ForumService
}

func NewClient(httpClient *http.Client) *Client {
    if httpClient == nil {
        httpClient = &http.Client{}
    }
    baseURL, _ := url.Parse(defaultBaseURL)
    c := &Client{client: httpClient, BaseURL: baseURL}
    c.Anime = &AnimeService{client: c}
    c.Manga = &MangaService{client: c}
    c.User  = &UserService{client: c}
    return c
}
```

Each service stores an unexported back-reference to `client *Client`, allowing it to access
`c.Do(ctx, req, v)` and shared HTTP transport without duplicating network logic.

---

## 4. Typed Functional Options

Rather than long parameter lists or config structs cluttered with nil pointers,
`go-myanimelist` utilizes interface-based functional options:

```go
type Option interface {
    apply(v *url.Values)
}

type Limit int
func (l Limit) apply(v *url.Values) { v.Set("limit", strconv.Itoa(int(l))) }
```

Specialized interfaces (`AnimeListOption`, `UpdateMyAnimeListStatusOption`) provide
compile-time safety, preventing invalid options from being passed to specific methods.

---

## 5. Authentication

- Delegates OAuth2 token management to `golang.org/x/oauth2`.
- Consumers initialize `Client` by passing an authenticated `*http.Client` produced by
  `oauth2.Config.Client(ctx, token)`.
- This cleanly separates token refresh and persistence from the API SDK logic.

---

## 6. Pagination

- As a REST API, MyAnimeList returns absolute links in a `Paging` object with `next` and `previous`.
- `go-myanimelist` parses the `offset` parameter from these URLs and exposes `NextOffset`
  and `PrevOffset` on the `Response` struct.

---

## 7. Error Handling

- Defines `ErrorResponse` implementing `error`:
  ```go
  type ErrorResponse struct {
      Response *http.Response
      Message  string `json:"message"`
      Err      string `json:"error"`
  }
  ```
- Evaluates HTTP status codes and parses JSON error bodies into structured types.
- Supports inspection via `errors.As`.

---

## 8. Testing Strategy

- Uses `net/http/httptest.Server` for isolated, offline unit tests:
  ```go
  func setup() (client *Client, mux *http.ServeMux, teardown func()) {
      mux = http.NewServeMux()
      server := httptest.NewServer(mux)
      client = NewClient(nil)
      client.BaseURL, _ = url.Parse(server.URL + "/")
      return client, mux, server.Close
  }
  ```
- Helper functions: `testMethod`, `testURLValues`, `testBody`, `testErrorResponse`.

---

## 9. Comparison: What to Adopt vs What NOT to Copy for `go-anilist`

| Pattern in `go-myanimelist` | Adopt in `go-anilist`? | Technical Justification |
|-----------------------------|:----------------------:|-------------------------|
| **Client / Services Pattern** | **YES** | Idiomatic Go structure; groups methods by entity (`Media`, `MediaList`, `Viewer`). |
| **`*http.Client` Injection** | **YES** | Allows injecting OAuth2 tokens, custom proxies, and `httptest` servers. |
| **Typed Functional Options** | **YES (adapted)** | Great ergonomics; instead of mutating `url.Values`, options set GraphQL variables. |
| **Enriched `Response` Struct** | **YES** | Critical for exposing `X-RateLimit-*` quota headers to callers. |
| **`httptest.Server` Strategy** | **YES (adapted)** | Fully deterministic testing without internet access. |
| **REST URLs & `url.Values`** | **NO** | AniList is 100% GraphQL via HTTP POST to a single endpoint. |
| **Offset URL Pagination** | **NO** | AniList uses `Page(page: $page, perPage: $perPage)` and `pageInfo`. |
| **Errors Solely on HTTP Status** | **NO** | GraphQL can return HTTP 200 with top-level `errors` array in JSON. |
| **URL Path Muxing in Tests** | **NO** | All requests hit `/`; test dispatchers must inspect JSON `query` or `operationName`. |