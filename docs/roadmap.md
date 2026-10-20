# Implementation Roadmap — `go-anilist`

This document details the development phases for the `go-anilist` client library.

---

## Phase 0: Architecture and Specification
- [x] AI agent guidelines and repository rules (`AGENTS.md`, `docs/ai/*`).
- [x] Architectural analysis of `go-myanimelist` (`docs/reference/go-myanimelist.md`).
- [x] AniList GraphQL API technical specifications (`docs/reference/anilist-api.md`).
- [x] Strict bounded scope definition (`docs/design/alcance.md`).
- [x] Evaluation of GraphQL integration approaches (`docs/design/graphql.md`).
- [x] Public API contracts and Go signatures (`docs/design/api-publica.md`).
- [x] Deterministic unit testing strategy with `httptest` (`docs/design/testing.md`).
- [x] SemVer policy and Conventional Commits relationship (`docs/design/versionado.md`).

---

## Phase 1: Core Package and Transport Layer
- [x] Module initialization (`go.mod`) with Go 1.27 and `CGO_ENABLED=0`.
- [x] `Client` struct, constructor `NewClient`, and functional options (`WithBaseURL`, `WithToken`, `WithUserAgent`, etc.).
- [x] Internal GraphQL transport layer (`doGraphQL`):
  - [x] JSON query and variables encoding.
  - [x] Streaming deserialization via `json.NewDecoder`.
  - [x] Quota header parsing (`X-RateLimit-*`, `Retry-After`).
  - [x] Structured error handling (`ErrorResponse`, `GraphQLError`, `ErrRateLimited`).
- [x] Unit test suite with `httptest.Server` for transport, context cancellation, and HTTP 429.

---

## Phase 2: User and Media Services (`User` & `Media`)
- [x] Implement `UserService`:
  - [x] `CurrentUser` query (`Viewer`) for user ID, username, and `ScoreFormat`.
  - [x] Unit tests with mock fixture `testdata/viewer_success.json`.
- [x] Implement `MediaService`:
  - [x] Text search (`Search`) with format, season year, and limit filters.
  - [x] Lookup by AniList ID (`GetByID`) and MyAnimeList ID (`GetByIDMal`).
  - [x] Unit tests with mock fixture `testdata/media_search_success.json`.

---

## Phase 3: MediaList Service (`MediaList`)
- [x] Implement `MediaListService`:
  - [x] Lookup entry (`GetEntry`) by `userId` and `mediaId`.
  - [x] Query user collection (`GetCollection`) filtered by status.
  - [x] Mutation `SaveEntry` to update progress, status (`CURRENT`, `COMPLETED`), and score.
  - [x] Mutation `DeleteEntry`.
- [x] Unit test suite for list queries and mutations with mock fixtures.

---

## Phase 4: OAuth2 Authentication
- [x] Implement `AuthService`:
  - [x] Generate authorization code grant URL (`AuthorizeURL`).
  - [x] Exchange code (`ExchangeCode`) against `/oauth/token`.
  - [x] Renew token via `RefreshToken`.
- [x] Unit tests with local mock OAuth2 endpoints.

---

## Phase 5: Stabilization, CI, and Initial Release
- [x] Full unit test validation (`go test -v -race ./...`).
- [x] Static code analysis with `go vet ./...`.
- [x] Automated GitHub Actions CI workflow (`.github/workflows/ci.yml`).
- [x] Initial release tag `v0.1.0`.