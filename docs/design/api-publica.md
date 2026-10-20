# Go Public API Specification — `go-anilist`

This document defines the exact Go public API signatures, models, options, and contracts
for `go-anilist`. The design is specifically shaped to allow `aniSync` to wrap it cleanly
behind its `Tracker` domain interface.

---

## 1. Package and Core Client

```go
package anilist

import (
    "context"
    "net/http"
    "time"
)

// Client manages communication with the AniList GraphQL API.
type Client struct {
    client    *http.Client
    baseURL   string
    authURL   string
    tokenURL  string
    token     string
    userAgent string

    Media     *MediaService
    MediaList *MediaListService
    User      *UserService
    Auth      *AuthService
}

// NewClient returns a new AniList API client.
// If httpClient is nil, a default http.Client with a 30s timeout is used.
func NewClient(httpClient *http.Client, opts ...ClientOption) *Client

// ClientOption configures a Client using functional options.
type ClientOption func(*Client)

// WithBaseURL overrides the default GraphQL base URL (useful for httptest.Server).
func WithBaseURL(url string) ClientOption

// WithAuthURL overrides the default OAuth2 authorization URL.
func WithAuthURL(url string) ClientOption

// WithTokenURL overrides the default OAuth2 token endpoint URL.
func WithTokenURL(url string) ClientOption

// WithToken sets a static OAuth2 Bearer token for authenticating requests.
func WithToken(token string) ClientOption

// WithUserAgent sets a custom User-Agent header for API requests.
func WithUserAgent(ua string) ClientOption
```

---

## 2. Domain Types and Models

```go
// MediaType represents if the media is anime or manga.
type MediaType string

const (
    MediaTypeAnime MediaType = "ANIME"
    MediaTypeManga MediaType = "MANGA"
)

// MediaFormat represents the broadcast or release format.
type MediaFormat string

const (
    MediaFormatTV      MediaFormat = "TV"
    MediaFormatTVShort MediaFormat = "TV_SHORT"
    MediaFormatMovie   MediaFormat = "MOVIE"
    MediaFormatSpecial MediaFormat = "SPECIAL"
    MediaFormatOVA     MediaFormat = "OVA"
    MediaFormatONA     MediaFormat = "ONA"
    MediaFormatMusic   MediaFormat = "MUSIC"
)

// MediaStatus represents the current releasing status.
type MediaStatus string

const (
    MediaStatusFinished       MediaStatus = "FINISHED"
    MediaStatusReleasing      MediaStatus = "RELEASING"
    MediaStatusNotYetReleased MediaStatus = "NOT_YET_RELEASED"
    MediaStatusCancelled      MediaStatus = "CANCELLED"
    MediaStatusHiatus         MediaStatus = "HIATUS"
)

// MediaListStatus represents the status of a media entry in a user's list.
type MediaListStatus string

const (
    MediaListStatusCurrent   MediaListStatus = "CURRENT"
    MediaListStatusPlanning  MediaListStatus = "PLANNING"
    MediaListStatusCompleted MediaListStatus = "COMPLETED"
    MediaListStatusDropped   MediaListStatus = "DROPPED"
    MediaListStatusPaused    MediaListStatus = "PAUSED"
    MediaListStatusRepeating MediaListStatus = "REPEATING"
)

// Media represents an anime or manga entry in AniList.
type Media struct {
    ID         int         `json:"id"`
    IDMal      *int        `json:"idMal,omitempty"`
    Title      MediaTitle  `json:"title"`
    Format     MediaFormat `json:"format"`
    Status     MediaStatus `json:"status"`
    Episodes   *int        `json:"episodes,omitempty"`
    Duration   *int        `json:"duration,omitempty"`
    Synonyms   []string    `json:"synonyms"`
    SeasonYear *int        `json:"seasonYear,omitempty"`
}

// MediaTitle contains titles in multiple language scripts.
type MediaTitle struct {
    Romaji  string  `json:"romaji"`
    English *string `json:"english,omitempty"`
    Native  *string `json:"native,omitempty"`
}

// MediaListEntry represents the status and progress of an anime in a user's list.
type MediaListEntry struct {
    ID        int             `json:"id"`
    MediaID   int             `json:"mediaId"`
    Status    MediaListStatus `json:"status"`
    Progress  int             `json:"progress"`
    Score     float64         `json:"score"`
    Repeat    int             `json:"repeat"`
    UpdatedAt int64           `json:"updatedAt"`
    Media     *Media          `json:"media,omitempty"`
}

// User represents a user profile in AniList.
type User struct {
    ID               int              `json:"id"`
    Name             string           `json:"name"`
    Avatar           *UserAvatar      `json:"avatar,omitempty"`
    MediaListOptions MediaListOptions `json:"mediaListOptions"`
}

// UserAvatar holds URLs for the user's avatar.
type UserAvatar struct {
    Large  string `json:"large,omitempty"`
    Medium string `json:"medium,omitempty"`
}

// MediaListOptions contains user preferences for media lists.
type MediaListOptions struct {
    ScoreFormat string `json:"scoreFormat"`
}
```

---

## 3. Services and Method Signatures

### 3.1. MediaService
```go
type MediaService struct {
    client *Client
}

func (s *MediaService) GetByID(ctx context.Context, id int) (*Media, *Response, error)
func (s *MediaService) GetByIDMal(ctx context.Context, idMal int) (*Media, *Response, error)
func (s *MediaService) Search(ctx context.Context, query string, opts ...SearchOption) ([]Media, *Response, error)

type SearchOption func(*searchParams)
func WithFormat(f MediaFormat) SearchOption
func WithSeasonYear(year int) SearchOption
func WithLimit(limit int) SearchOption
```

### 3.2. MediaListService
```go
type MediaListService struct {
    client *Client
}

func (s *MediaListService) GetEntry(ctx context.Context, userID, mediaID int) (*MediaListEntry, *Response, error)
func (s *MediaListService) GetCollection(ctx context.Context, userID int, opts ...CollectionOption) ([]MediaListEntry, *Response, error)
func (s *MediaListService) SaveEntry(ctx context.Context, input SaveMediaListEntryInput) (*MediaListEntry, *Response, error)
func (s *MediaListService) DeleteEntry(ctx context.Context, id int) (bool, *Response, error)

type CollectionOption func(*collectionParams)
func WithStatus(status MediaListStatus) CollectionOption
```

### 3.3. UserService
```go
type UserService struct {
    client *Client
}

func (s *UserService) CurrentUser(ctx context.Context) (*User, *Response, error)
func (s *UserService) GetByID(ctx context.Context, id int) (*User, *Response, error)
```

### 3.4. AuthService
```go
type AuthService struct {
    client *Client
}

type Token struct {
    AccessToken  string    `json:"access_token"`
    TokenType    string    `json:"token_type"`
    ExpiresIn    int       `json:"expires_in"`
    RefreshToken string    `json:"refresh_token"`
    Expiry       time.Time `json:"expiry"`
}

func (s *AuthService) AuthorizeURL(clientID, redirectURI, state string) string
func (s *AuthService) ExchangeCode(ctx context.Context, clientID, clientSecret, redirectURI, code string) (*Token, error)
func (s *AuthService) RefreshToken(ctx context.Context, clientID, clientSecret, refreshToken string) (*Token, error)
```

---

## 4. Response Metadata and Rate Limiting

```go
type Response struct {
    *http.Response
    RateLimit RateLimitInfo
}

type RateLimitInfo struct {
    Limit      int
    Remaining  int
    Reset      time.Time
    RetryAfter time.Duration
}
```

---

## 5. Structured and Sentinel Errors

```go
var (
    ErrNotFound     = errors.New("anilist: resource not found")
    ErrUnauthorized = errors.New("anilist: unauthorized or invalid token")
    ErrRateLimited  = errors.New("anilist: rate limit exceeded (429)")
)

type GraphQLError struct {
    Message   string     `json:"message"`
    Status    int        `json:"status"`
    Locations []Location `json:"locations,omitempty"`
}

type ErrorResponse struct {
    Response *http.Response
    Errors   []GraphQLError
    Message  string
}
```

---

## 6. Integration in `aniSync` via `Tracker` Interface

`aniSync` wraps `go-anilist` cleanly behind its domain interface:

```go
type Tracker interface {
    Name() string
    SearchAnime(ctx context.Context, title string, year int) ([]domain.AnimeRef, error)
    GetWatchStatus(ctx context.Context, animeID string) (*domain.WatchStatus, error)
    UpdateProgress(ctx context.Context, animeID string, episode int, completed bool) error
}
```