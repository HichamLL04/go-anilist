package anilist

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithBaseURL overrides the default GraphQL base URL.
// Useful for targeting test servers created with net/http/httptest.
func WithBaseURL(rawURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = rawURL
	}
}

// WithToken sets a static OAuth2 Bearer token for authenticating requests.
func WithToken(token string) ClientOption {
	return func(c *Client) {
		c.token = token
	}
}

// WithUserAgent sets a custom User-Agent header for API requests.
func WithUserAgent(ua string) ClientOption {
	return func(c *Client) {
		c.userAgent = ua
	}
}
