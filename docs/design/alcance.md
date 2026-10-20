# Library Scope and Boundaries — `go-anilist`

This document defines the functional perimeter of `go-anilist`.
The library does **not** attempt to cover the entirety of AniList's schema,
focusing strictly on the needs of **aniSync** and media tracking tools.

---

## 1. Boundary Principle

The official AniList v2 GraphQL schema contains hundreds of types (including forum threads,
social feeds, community reviews, character voice actors, and moderation logs).
Building a universal client would:
- Multiply compiled binary size and build times.
- Increase runtime allocations with zero added value.
- Introduce an unmanageable maintenance burden.

Therefore, `go-anilist` follows a **strict minimum viable scope**, focusing on
searching, querying, and updating anime watch progress.

---

## 2. In-Scope Components

| Component / Service | Methods & Entities | Technical Justification in aniSync |
|---------------------|--------------------|-----------------------------------|
| **Auth / OAuth2** | `AuthorizeURL`, `ExchangeCode`, `RefreshToken` | Enables linking user AniList accounts persistently and unattended. |
| **User (`UserService`)** | `CurrentUser(ctx)`, `GetByID(ctx, id)` | Identifies the authenticated user, numeric ID, and rating scale (`ScoreFormat`). |
| **Media (`MediaService`)** | `GetByID(ctx, id)`, `GetByIDMal(ctx, idMal)`, `Search(ctx, query, opts...)` | Matches Plex titles (Romaji, English, Native) and cross-references MyAnimeList IDs. |
| **MediaList (`MediaListService`)** | `GetEntry(ctx, userID, mediaID)`, `GetCollection(ctx, userID, opts...)`, `SaveEntry(ctx, input)`, `DeleteEntry(ctx, id)` | Core scrobbling: checks watch status, mutates episode progress and completion (`CURRENT`, `COMPLETED`). |
| **Transport & Quotas** | `doGraphQL`, `X-RateLimit-*` headers, HTTP 429 handling | Guarantees compliance with AniList's 90 req/min quota and reports rate limits to caller. |

---

## 3. Out-of-Scope Components

| Schema Area | Status | Reason for Exclusion |
|-------------|:------:|----------------------|
| **Forums & Threads (`Thread`, `Comment`)** | **Excluded** | Irrelevant for video watch progress sync. |
| **Social Feed (`Activity`, `Like`, `Follow`)** | **Excluded** | aniSync is an unattended daemon; does not consume social features. |
| **Reviews & Recommendations (`Review`)** | **Excluded** | High type overhead with zero utility in media matching. |
| **Characters & Staff (`Character`, `Staff`)** | **Excluded** | Large relational payload unused in scrobbling. |
| **System Notifications (`Notification`)** | **Excluded** | Does not participate in synchronization. |
| **Manga / Light Novels (Reading Lists)** | **Deferred** | aniSync focuses on video streams. Base enums (`MediaType`) are provided for future extension. |
| **Administrative Mutations (`UpdateUser`)** | **Excluded** | Account configuration is managed via the AniList web interface. |