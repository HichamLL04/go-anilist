# AniList GraphQL API (v2) Specification

Technical reference documentation for the AniList GraphQL API based on the
official documentation (https://docs.anilist.co) **[DOC-ANILIST]**.

---

## 1. HTTP Endpoint and Protocol

- **Production Endpoint**: `https://graphql.anilist.co` **[DOC-ANILIST]**
- **Method**: Exclusively `POST`. `GET` requests are rejected.
- **Mandatory Headers**:
  ```http
  Content-Type: application/json
  Accept: application/json
  ```
- **Authentication Header** (for user-authenticated queries and mutations):
  ```http
  Authorization: Bearer <access_token>
  ```
- **Request Body Format**:
  ```json
  {
    "query": "query ($id: Int) { Media(id: $id) { id title { romaji english } } }",
    "variables": { "id": 1 }
  }
  ```

---

## 2. Authentication and OAuth2

AniList uses the OAuth 2.0 protocol for user authorization and mutations **[DOC-ANILIST]**:

- **Authorize URL**: `https://anilist.co/api/v2/oauth/authorize`
- **Token URL**: `https://anilist.co/api/v2/oauth/token`

### Supported Grant Types

1. **Authorization Code Grant** (Recommended for automated backends like aniSync):
   - Redirect user to `/oauth/authorize?client_id={CLIENT_ID}&response_type=code&redirect_uri={REDIRECT_URI}`.
   - Exchange code via POST to `/oauth/token`:
     ```json
     {
       "grant_type": "authorization_code",
       "client_id": "{CLIENT_ID}",
       "client_secret": "{CLIENT_SECRET}",
       "redirect_uri": "{REDIRECT_URI}",
       "code": "{AUTHORIZATION_CODE}"
     }
     ```
   - Returns `access_token`, `token_type: "Bearer"`, `expires_in`, and `refresh_token`.
2. **Implicit Grant** (For manual static tokens or single-user desktop utilities):
   - Direct authorization returning token in URL hash fragment (`#access_token=...`).

### Token Expiration and Validity

- AniList access tokens are typically valid for **1 year (31,536,000 seconds)** **[DOC-ANILIST]**.
- When using the Authorization Code grant, the `refresh_token` allows automated renewal without requiring manual re-authentication.

---

## 3. Rate Limiting (Actual Quotas)

- **General Quota**: **90 requests per minute** per IP or per OAuth2 token **[DOC-ANILIST]**.
- When exceeding the quota, the API responds with HTTP **`429 Too Many Requests`**.
- **HTTP Response Headers**:
  - `X-RateLimit-Limit`: Maximum requests per window (typically `90`).
  - `X-RateLimit-Remaining`: Remaining requests in current 60-second window.
  - `X-RateLimit-Reset`: Unix timestamp (in seconds) when the rate limit window resets.
  - `Retry-After`: In HTTP 429 responses, specifies the exact seconds to wait before retrying.

---

## 4. Error Handling and Response Envelope

AniList returns standard GraphQL envelopes **[DOC-ANILIST]**:

### Successful Response
```json
{
  "data": { ... }
}
```

### Error Response
```json
{
  "errors": [
    {
      "message": "Validation failed: Score must be between 0 and 100",
      "status": 400,
      "locations": [{ "line": 2, "column": 3 }]
    }
  ]
}
```

### HTTP Status Codes
- **200 OK**: Request completed (even if partial field errors exist in `errors`).
- **400 Bad Request**: Invalid GraphQL syntax, variable validation failures, or type mismatches.
- **401 Unauthorized**: Missing, invalid, or expired Bearer token.
- **404 Not Found**: Resource not found (e.g. Media ID does not exist).
- **429 Too Many Requests**: Rate limit of 90 requests/minute exceeded.
- **500 / 502 / 503 Server Error**: AniList or upstream Cloudflare failure.

---

## 5. Pagination

- Collections are wrapped in `Page(page: $page, perPage: $perPage)` **[DOC-ANILIST]**.
- **Strict Limit**: Maximum `perPage` allowed by the API is **50**.
- **`pageInfo` Structure**:
  ```graphql
  pageInfo {
    total
    perPage
    currentPage
    lastPage
    hasNextPage
  }
  ```

---

## 6. Required Queries and Mutations for aniSync

### 1. Authenticated User (`Viewer`)
```graphql
query CurrentUser {
  Viewer {
    id
    name
    avatar {
      large
      medium
    }
    mediaListOptions {
      scoreFormat
    }
  }
}
```

### 2. Anime Search and Lookup (`Media`)
```graphql
query SearchAnime($search: String, $format: MediaFormat, $seasonYear: Int, $perPage: Int) {
  Page(page: 1, perPage: $perPage) {
    media(search: $search, format: $format, seasonYear: $seasonYear, type: ANIME) {
      id
      idMal
      title {
        romaji
        english
        native
      }
      format
      status
      episodes
      duration
      synonyms
      seasonYear
    }
  }
}
```

### 3. User Media List Entry (`MediaList`)
```graphql
query GetMediaListEntry($userId: Int, $mediaId: Int) {
  MediaList(userId: $userId, mediaId: $mediaId, type: ANIME) {
    id
    mediaId
    status
    progress
    score
    repeat
    updatedAt
  }
}
```

### 4. User Collection (`MediaListCollection`)
```graphql
query GetUserAnimeList($userId: Int, $status: MediaListStatus) {
  MediaListCollection(userId: $userId, type: ANIME, status: $status) {
    lists {
      name
      isCustomList
      status
      entries {
        id
        mediaId
        status
        progress
        score
        repeat
        updatedAt
        media {
          id
          idMal
          title {
            romaji
            english
          }
          episodes
        }
      }
    }
  }
}
```

### 5. Update Watch Progress (`SaveMediaListEntry`)
```graphql
mutation UpdateProgress(
  $mediaId: Int
  $status: MediaListStatus
  $progress: Int
  $score: Float
  $repeat: Int
) {
  SaveMediaListEntry(
    mediaId: $mediaId
    status: $status
    progress: $progress
    score: $score
    repeat: $repeat
  ) {
    id
    mediaId
    status
    progress
    score
    repeat
    updatedAt
  }
}
```

### 6. Delete Entry (`DeleteMediaListEntry`)
```graphql
mutation DeleteEntry($id: Int) {
  DeleteMediaListEntry(id: $id) {
    deleted
  }
}
```