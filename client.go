package anilist

import (
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the default GraphQL API endpoint for AniList.
	DefaultBaseURL = "https://graphql.anilist.co"

	// DefaultTimeout is the default timeout for the HTTP client.
	DefaultTimeout = 30 * time.Second

	// DefaultUserAgent is the default User-Agent header sent with requests.
	DefaultUserAgent = "go-anilist/0.1.0"
)

// Client manages communication with the AniList GraphQL API.
type Client struct {
	client    *http.Client
	baseURL   string
	authURL   string
	tokenURL  string
	token     string
	userAgent string

	// Services
	Media     *MediaService
	MediaList *MediaListService
	User      *UserService
	Auth      *AuthService
}

// MediaService handles communication with the media related methods of the AniList API.
type MediaService struct {
	client *Client
}

// MediaListService handles communication with the media list methods of the AniList API.
type MediaListService struct {
	client *Client
}

// UserService handles communication with the user related methods of the AniList API.
type UserService struct {
	client *Client
}

// AuthService handles authentication with the AniList API.
type AuthService struct {
	client *Client
}

// NewClient returns a new AniList API client. If httpClient is nil, a default
// http.Client with DefaultTimeout will be used.
func NewClient(httpClient *http.Client, opts ...ClientOption) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: DefaultTimeout,
		}
	}

	c := &Client{
		client:    httpClient,
		baseURL:   DefaultBaseURL,
		authURL:   DefaultAuthURL,
		tokenURL:  DefaultTokenURL,
		userAgent: DefaultUserAgent,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Media = &MediaService{client: c}
	c.MediaList = &MediaListService{client: c}
	c.User = &UserService{client: c}
	c.Auth = &AuthService{client: c}

	return c
}