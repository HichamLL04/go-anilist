# Especificación de la API Pública de Go — `go-anilist`

Este documento define las firmas exactas, tipos de datos, opciones y contratos
de la API pública de `go-anilist`. El diseño está pensado específicamente para que
`aniSync` pueda consumirla y envolverla de forma limpia tras su interfaz `Tracker`.

---

## 1. Paquete y Cliente Principal

```go
package anilist

import (
    "context"
    "net/http"
    "time"
)

// Client gestiona la comunicación con la API GraphQL de AniList.
type Client struct {
    client  *http.Client
    baseURL string
    token   string

    Media     *MediaService
    MediaList *MediaListService
    User      *UserService
    Auth      *AuthService
}

// NewClient inicializa un nuevo cliente de AniList.
// Si httpClient es nil, se utiliza un http.Client por defecto con timeout de 30s.
func NewClient(httpClient *http.Client, opts ...ClientOption) *Client

// ClientOption permite configurar el cliente mediante opciones funcionales.
type ClientOption func(*Client)

// WithBaseURL sobreescribe la URL base (útil para httptest.Server en tests).
func WithBaseURL(url string) ClientOption

// WithToken configura un Bearer token estático para todas las peticiones.
func WithToken(token string) ClientOption

// WithUserAgent configura la cabecera User-Agent de las peticiones.
func WithUserAgent(ua string) ClientOption
```

---

## 2. Tipos de Dominio y Modelos de Datos

```go
// MediaType representa si el medio es anime o manga.
type MediaType string

const (
    MediaTypeAnime MediaType = "ANIME"
    MediaTypeManga MediaType = "MANGA"
)

// MediaFormat representa el formato de emisión.
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

// MediaStatus estado de emisión de la obra.
type MediaStatus string

const (
    MediaStatusFinished        MediaStatus = "FINISHED"
    MediaStatusReleasing       MediaStatus = "RELEASING"
    MediaStatusNotYetReleased  MediaStatus = "NOT_YET_RELEASED"
    MediaStatusCancelled       MediaStatus = "CANCELLED"
    MediaStatusHiatus          MediaStatus = "HIATUS"
)

// MediaListStatus estado de seguimiento de la entrada de un usuario.
type MediaListStatus string

const (
    MediaListStatusCurrent   MediaListStatus = "CURRENT"
    MediaListStatusPlanning  MediaListStatus = "PLANNING"
    MediaListStatusCompleted MediaListStatus = "COMPLETED"
    MediaListStatusDropped   MediaListStatus = "DROPPED"
    MediaListStatusPaused    MediaListStatus = "PAUSED"
    MediaListStatusRepeating MediaListStatus = "REPEATING"
)

// Media representa una obra (anime) en AniList.
type Media struct {
    ID          int         `json:"id"`
    IDMal       *int        `json:"idMal,omitempty"`
    Title       MediaTitle  `json:"title"`
    Format      MediaFormat `json:"format"`
    Status      MediaStatus `json:"status"`
    Episodes    *int        `json:"episodes,omitempty"`
    Duration    *int        `json:"duration,omitempty"`
    Synonyms    []string    `json:"synonyms"`
    SeasonYear  *int        `json:"seasonYear,omitempty"`
}

// MediaTitle almacena variantes de títulos de una obra.
type MediaTitle struct {
    Romaji  string  `json:"romaji"`
    English *string `json:"english,omitempty"`
    Native  *string `json:"native,omitempty"`
}

// MediaListEntry representa el estado de una obra en la lista del usuario.
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

// User representa al usuario en AniList.
type User struct {
    ID               int              `json:"id"`
    Name             string           `json:"name"`
    AvatarURL        string           `json:"avatarUrl"`
    MediaListOptions MediaListOptions `json:"mediaListOptions"`
}

// MediaListOptions opciones de configuración de listas de usuario.
type MediaListOptions struct {
    ScoreFormat string `json:"scoreFormat"` // "POINT_100", "POINT_10_DECIMAL", "POINT_10", "POINT_5", "POINT_3"
}
```

---

## 3. Servicios y Firmas de Métodos

### 3.1. MediaService

```go
type MediaService struct {
    client *Client
}

// GetByID obtiene un medio por su ID numérico de AniList.
func (s *MediaService) GetByID(ctx context.Context, id int) (*Media, *Response, error)

// GetByIDMal obtiene un medio a partir de su ID de MyAnimeList.
func (s *MediaService) GetByIDMal(ctx context.Context, idMal int) (*Media, *Response, error)

// Search busca medios por texto y filtros opcionales.
func (s *MediaService) Search(ctx context.Context, query string, opts ...SearchOption) ([]Media, *Response, error)

// SearchOption define opciones de búsqueda.
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

// GetEntry obtiene la entrada de un anime específico en la lista de un usuario.
// Devuelve ErrNotFound si el anime no está en su lista.
func (s *MediaListService) GetEntry(ctx context.Context, userID, mediaID int) (*MediaListEntry, *Response, error)

// GetCollection obtiene la lista de seguimiento de un usuario filtrada por estado.
func (s *MediaListService) GetCollection(ctx context.Context, userID int, opts ...CollectionOption) ([]MediaListEntry, *Response, error)

// SaveMediaListEntryInput parámetros de entrada para guardar o actualizar una entrada.
type SaveMediaListEntryInput struct {
    MediaID  int              `json:"mediaId"`
    Status   *MediaListStatus `json:"status,omitempty"`
    Progress *int             `json:"progress,omitempty"`
    Score    *float64         `json:"score,omitempty"`
    Repeat   *int             `json:"repeat,omitempty"`
}

// SaveEntry crea o actualiza el progreso y estado de un anime en la lista del usuario autenticado.
func (s *MediaListService) SaveEntry(ctx context.Context, input SaveMediaListEntryInput) (*MediaListEntry, *Response, error)

// DeleteEntry elimina una entrada de la lista por su entry ID.
func (s *MediaListService) DeleteEntry(ctx context.Context, id int) (bool, *Response, error)
```

### 3.3. UserService

```go
type UserService struct {
    client *Client
}

// CurrentUser obtiene los datos del usuario autenticado (query Viewer).
func (s *UserService) CurrentUser(ctx context.Context) (*User, *Response, error)

// GetByID obtiene el perfil público de un usuario por su ID.
func (s *UserService) GetByID(ctx context.Context, id int) (*User, *Response, error)
```

### 3.4. AuthService

```go
type AuthService struct {
    client *Client
}

// Token representa las credenciales OAuth2 devueltas por AniList.
type Token struct {
    AccessToken  string    `json:"access_token"`
    TokenType    string    `json:"token_type"`
    ExpiresIn    int       `json:"expires_in"`
    RefreshToken string    `json:"refresh_token"`
    Expiry       time.Time `json:"expiry"`
}

// AuthorizeURL genera la URL para iniciar el flujo de autorización OAuth2.
func (s *AuthService) AuthorizeURL(clientID, redirectURI, state string) string

// ExchangeCode intercambia el código de autorización por un Token de acceso.
func (s *AuthService) ExchangeCode(ctx context.Context, clientID, clientSecret, redirectURI, code string) (*Token, error)

// RefreshToken solicita un nuevo AccessToken utilizando el RefreshToken almacenado.
func (s *AuthService) RefreshToken(ctx context.Context, clientID, clientSecret, refreshToken string) (*Token, error)
```

---

## 4. Metadatos de Respuesta y Rate Limiting

```go
// Response envuelve http.Response y extrae las métricas de consumo de AniList.
type Response struct {
    *http.Response
    RateLimit RateLimitInfo
}

// RateLimitInfo informa sobre las cuotas de la API.
type RateLimitInfo struct {
    Limit      int           // Valor de X-RateLimit-Limit (normalmente 90)
    Remaining  int           // Valor de X-RateLimit-Remaining
    Reset      time.Time     // Instante de reinicio de ventana según X-RateLimit-Reset
    RetryAfter time.Duration // Tiempo a esperar si se recibe un 429
}
```

---

## 5. Errores Tipados y Centinela

```go
var (
    ErrNotFound     = errors.New("anilist: resource not found")
    ErrUnauthorized = errors.New("anilist: unauthorized or invalid token")
    ErrRateLimited  = errors.New("anilist: rate limit exceeded (429)")
)

// GraphQLError describe un error individual devuelto en el array "errors".
type GraphQLError struct {
    Message   string     `json:"message"`
    Status    int        `json:"status"`
    Locations []Location `json:"locations,omitempty"`
}

// Location indica línea y columna del error en la query GraphQL.
type Location struct {
    Line   int `json:"line"`
    Column int `json:"column"`
}

// ErrorResponse reporta un error HTTP o GraphQL de la API.
type ErrorResponse struct {
    Response *http.Response
    Errors   []GraphQLError
}

func (e *ErrorResponse) Error() string
```

---

## 6. Integración en `aniSync` mediante la interfaz `Tracker`

En `aniSync`, la interfaz del núcleo de sincronización envolverá a `go-anilist`
de forma transparente:

```go
// internal/sync/tracker.go en aniSync
type Tracker interface {
    Name() string
    SearchAnime(ctx context.Context, title string, year int) ([]domain.AnimeRef, error)
    GetWatchStatus(ctx context.Context, animeID string) (*domain.WatchStatus, error)
    UpdateProgress(ctx context.Context, animeID string, episode int, completed bool) error
}
```

La implementación concreta `internal/tracker/anilist` instanciará `anilist.Client`,
recibirá las llamadas de `aniSync`, invocará `SaveEntry` o `Search`, y traducirá los
modelos de `go-anilist` a los tipos de dominio normalizados de `aniSync`.
