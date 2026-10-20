# Go Style Guidelines — go-anilist

Technical and code style standards for Go code in the `go-anilist` library.

---

## 1. Error Handling

- **Always wrap with `%w`**: All propagated errors must be wrapped using
  `fmt.Errorf("anilist: <op>: %w", err)`. Never discard errors silently with `_`
  in critical execution paths.
- **Typed and Sentinel Errors**:
  - Define sentinel errors for predictable conditions (e.g. `ErrNotFound`, `ErrRateLimited`, `ErrUnauthorized`).
  - Provide structured types `GraphQLError` and `ErrorResponse` exposing HTTP status codes, AniList error messages, and query locations.
  - Implement `errors.Is` / `errors.As` support so consumers can inspect underlying causes.

---

## 2. Context and Cancellation

- **`context.Context` as First Parameter**: Any method or function performing HTTP, I/O, or asynchronous work must accept `ctx context.Context`:
  ```go
  func (s *MediaService) GetByID(ctx context.Context, id int) (*Media, *Response, error)
  ```
- **Respect Cancellation**: Propagate the context to `http.NewRequestWithContext(ctx, ...)` and check `ctx.Err()`.
- If `ctx == nil`, reject the request immediately with an error before initiating any network traffic.

---

## 3. State, Global Variables, and `init()` Functions

- **No Mutable Global Variables**: Zero shared mutable global state. Each `Client` instance manages its own configuration and HTTP transport independently.
- **No `init()` Functions**: No hidden registration in `init()`. Configuration must be explicit via `NewClient` functional options.
- Constants for endpoints and URLs must be immutable.

---

## 4. Interface and API Design

- **Small, Consumer-Driven Interfaces**:
  - `go-anilist` exposes concrete structs for services (`MediaService`, `MediaListService`, etc.) and returned domain types.
  - Do not force giant interfaces inside the library; interfaces should be defined by the consumer (e.g., in `aniSync` via its `Tracker` interface).
- **Typed Functional Options**: Use typed functional options for optional parameters and filters (e.g., `WithFormat(MediaFormatTV)`, `WithLimit(10)`).
- **Dependency Injection**:
  - The client must accept an injected `*http.Client`:
    ```go
    func NewClient(httpClient *http.Client, opts ...ClientOption) *Client
    ```
  - If `httpClient == nil`, default to an `http.Client` with standard timeouts (30s).

---

## 5. Testing Strategy

- **Table-Driven Tests**: Use test tables with named subtests via `t.Run(tc.name, func(t *testing.T) { ... })`.
- **Zero Real Network Calls**: Unit tests must never depend on live AniList servers. Rely exclusively on `net/http/httptest.Server` and mock JSON fixtures.
- **Race Detector**: All tests must pass cleanly under `go test -v -race ./...`.

---

## 6. Build and Compilation

- **`CGO_ENABLED=0`**: All code must build cleanly without CGO for static portability.
- **External Dependencies**:
  - Keep dependencies to the absolute minimum. Always prefer the Go standard library (`net/http`, `encoding/json`, `context`, `time`).
  - Do not add third-party GraphQL client libraries without formal justification in `docs/design/graphql.md`.