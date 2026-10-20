# go-anilist

[![Go Reference](https://pkg.go.dev/badge/github.com/hichamll04/go-anilist.svg)](https://pkg.go.dev/github.com/hichamll04/go-anilist)
[![CI](https://github.com/hichamll04/go-anilist/actions/workflows/ci.yml/badge.svg)](https://github.com/hichamll04/go-anilist/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A lightweight, standalone Go client library for the [AniList](https://anilist.co) GraphQL API (v2).
Designed with **zero external dependencies** (standard library only), direct streaming deserialization, and strict memory efficiency for background sync daemons and automated media services.

---

## Features

- **Zero External Dependencies**: Built entirely upon the Go standard library (`net/http`, `encoding/json`, `context`, `time`, `net/url`).
- **Streaming JSON**: Direct streaming decoding with `json.NewDecoder` without buffering redundant responses in memory.
- **Rate Limit Aware**: Automatic inspection of AniList quota headers (`X-RateLimit-*`) and native handling of HTTP `429 Too Many Requests` (`Retry-After`).
- **Modular Services**:
  - `User`: Authenticated profile (`Viewer`) and list preferences (`ScoreFormat`).
  - `Media`: Anime search and lookup by AniList ID and MyAnimeList ID (`idMal`).
  - `MediaList`: Watch status tracking, user collections, and progress mutations (`SaveEntry`).
  - `Auth`: Full OAuth2 authorization helpers (`AuthorizeURL`, `ExchangeCode`, `RefreshToken`).
- **100% Deterministic Testing**: Complete test suite using `net/http/httptest` with zero real-network requests.

---

## Installation

```bash
go get github.com/hichamll04/go-anilist
```

---

## Usage Examples

### 1. Client Initialization

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/hichamll04/go-anilist"
)

func main() {
    // Initialize client with a static OAuth2 Bearer token
    client := anilist.NewClient(nil, anilist.WithToken("your_access_token_here"))

    // Fetch authenticated user profile
    user, resp, err := client.User.CurrentUser(context.Background())
    if err != nil {
        log.Fatalf("Error fetching viewer: %v", err)
    }

    fmt.Printf("User: %s (ID: %d)\n", user.Name, user.ID)
    fmt.Printf("Remaining API quota: %d requests\n", resp.RateLimit.Remaining)
}
```

### 2. Anime Search and Lookup

```go
results, _, err := client.Media.Search(
    ctx,
    "Sousou no Frieren",
    anilist.WithFormat(anilist.MediaFormatTV),
    anilist.WithSeasonYear(2023),
    anilist.WithLimit(5),
)
if err != nil {
    log.Fatal(err)
}

for _, anime := range results {
    fmt.Printf("[%d] %s (%d episodes)\n", anime.ID, anime.Title.Romaji, *anime.Episodes)
}
```

### 3. Update Watch Progress

```go
progress := 16
status := anilist.MediaListStatusCurrent
score := 9.5

entry, _, err := client.MediaList.SaveEntry(ctx, anilist.SaveMediaListEntryInput{
    MediaID:  154587,
    Progress: &progress,
    Status:   &status,
    Score:    &score,
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated media %d: Episode %d (%s)\n", entry.MediaID, entry.Progress, entry.Status)
```

### 4. Structured Error Handling

```go
if errors.Is(err, anilist.ErrNotFound) {
    fmt.Println("Requested resource not found on AniList")
} else if errors.Is(err, anilist.ErrRateLimited) {
    fmt.Println("Rate limit exceeded (HTTP 429)")
} else if errors.Is(err, anilist.ErrUnauthorized) {
    fmt.Println("Invalid or expired OAuth2 token")
}
```

---

## Acknowledgements & Attribution

The client design (sub-service segregation via `Client.Service`, clean `*http.Client` injection, functional options, and deterministic `httptest.Server` test methodology) is inspired by the [`go-myanimelist`](https://github.com/nstratos/go-myanimelist) library created by **Stratos Neiros** (distributed under the MIT License).

`go-anilist` adapts these patterns for the AniList GraphQL v2 API, maintaining idiomatic Go conventions and zero external dependencies.

---

## License

Distributed under the [MIT License](LICENSE).