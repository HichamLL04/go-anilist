# Open Questions and Pending Decisions — `go-anilist`

Log of technical and design decisions under evaluation for `go-anilist`.

---

## 1. Proactive Rate Limiter vs Delegated to Consumer

- **Question**: AniList imposes a strict quota of 90 requests/minute. Should `go-anilist` include an optional token-bucket rate limiter (e.g. `golang.org/x/time/rate`), or should it remain zero-dependency by passively parsing `RateLimitInfo` and returning `ErrRateLimited`?
- **Options**:
  - *Option A (Adopted)*: Maintain zero external dependencies. The client passively extracts `X-RateLimit-*` headers and handles HTTP `429 Retry-After`. Queue throttling is managed by the consumer (`aniSync`).
  - *Option B*: Add an optional client-side rate limiter.
- **Status**: Resolved (Option A).

---

## 2. Normalization of Score Formats (`ScoreFormat`)

- **Question**: AniList allows users to configure different rating formats:
  - `POINT_100`: Score 0-100.
  - `POINT_10_DECIMAL`: Score 0.0-10.0 (one decimal place).
  - `POINT_10`: Integer 0-10.
  - `POINT_5`: 1-5 stars.
  - `POINT_3`: 3 smiley faces.
  Should `go-anilist` automatically normalize scores or pass the raw float64 value?
- **Options**:
  - *Option A (Adopted)*: Accept the raw float64 in `SaveMediaListEntryInput.Score (*float64)` and let the consumer adapt according to `User.MediaListOptions.ScoreFormat`.
  - *Option B*: Provide conversion helper methods.
- **Status**: Resolved (Option A).

---

## 3. Retaining Raw HTTP Response Body in `Response`

- **Question**: `go-myanimelist` retains `Response.Body []byte`. In memory-constrained environments, this duplicates heap allocations.
- **Options**:
  - *Option A (Adopted)*: Do not buffer `Body []byte` on successful responses (decode via streaming). Only retain diagnostic error messages on failure in `ErrorResponse`.
  - *Option B*: Provide a debug logging option (`WithDebug(true)`).
- **Status**: Resolved (Option A).