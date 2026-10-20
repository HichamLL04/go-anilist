# Architectural Invariants — go-anilist

Structural rules and invariants that must be upheld at all times in `go-anilist`.

---

## 1. Independence and Decoupling

- **Completely Agnostic of aniSync**: `go-anilist` is an autonomous library.
  Under no circumstances may it import `aniSync` packages, reference Plex, SQLite,
  media matching algorithms, or webhook handlers.
- **Zero Global State**: Multiple `Client` instances must be able to run concurrently
  and safely within the same process without interference.

---

## 2. Layer Separation: GraphQL Transport vs Typed Public API

1. **Internal Transport Layer (`doGraphQL`)**:
   - Serializes GraphQL requests (`query`, `variables`).
   - Sends HTTP POST requests to `https://graphql.anilist.co`.
   - Parses AniList quota headers (`X-RateLimit-*`).
   - Deserializes the standard GraphQL envelope:
     ```json
     {
       "data": { ... },
       "errors": [ ... ]
     }
     ```
   - Normalizes HTTP status codes and GraphQL error arrays into enriched Go error types.

2. **Public Service Layer (`MediaService`, `MediaListService`, `UserService`, `AuthService`)**:
   - Exposes high-level methods with idiomatic Go domain signatures:
     - `media.GetByID(ctx, id)`
     - `media.Search(ctx, query, opts...)`
     - `mediaList.GetEntry(ctx, userID, mediaID)`
     - `mediaList.SaveEntry(ctx, input)`
   - Completely encapsulates and hides raw GraphQL query strings and mutations.
   - Returns strongly typed Go domain models (`*Media`, `*MediaListEntry`, `*User`).

---

## 3. Rate Limiting and Quota Management

- **Strict Header Inspection**: The transport layer extracts `X-RateLimit-Limit`,
  `X-RateLimit-Remaining`, and `X-RateLimit-Reset` on every response.
- **Response Metadata**: Public methods return `(*T, *Response, error)` where `Response`
  includes rate limit metadata so consumers (like aniSync) can schedule intelligent throttling.
- **HTTP 429 Compliance**: When receiving HTTP 429 Too Many Requests, extract `Retry-After`
  and return structured `ErrRateLimited`.

---

## 4. `http.Client` Lifecycle and Injection

- The client must never leave orphaned HTTP transport resources.
- Supports custom `http.Client` injection for:
  - Delegated OAuth2 authentication.
  - Proxies, metrics, or logging roundtrippers.
  - Local testing servers with `httptest.Server`.