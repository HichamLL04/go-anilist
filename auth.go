package anilist

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultAuthURL is the standard AniList OAuth2 authorization endpoint.
	DefaultAuthURL = "https://anilist.co/api/v2/oauth/authorize"

	// DefaultTokenURL is the standard AniList OAuth2 token endpoint.
	DefaultTokenURL = "https://anilist.co/api/v2/oauth/token"
)

// Token represents OAuth2 credentials returned by the AniList authorization server.
type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

// AuthorizeURL returns the URL to direct users to for initiating the OAuth2 code grant flow.
func (s *AuthService) AuthorizeURL(clientID, redirectURI, state string) string {
	baseURL := s.client.authURL
	if baseURL == "" {
		baseURL = DefaultAuthURL
	}

	v := url.Values{}
	v.Set("client_id", clientID)
	v.Set("response_type", "code")
	if redirectURI != "" {
		v.Set("redirect_uri", redirectURI)
	}
	if state != "" {
		v.Set("state", state)
	}

	return fmt.Sprintf("%s?%s", baseURL, v.Encode())
}

type tokenExchangeRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri,omitempty"`
	Code         string `json:"code,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// ExchangeCode exchanges an authorization code for an OAuth2 Token.
func (s *AuthService) ExchangeCode(ctx context.Context, clientID, clientSecret, redirectURI, code string) (*Token, error) {
	if ctx == nil {
		return nil, errors.New("anilist: context must not be nil")
	}

	payload := tokenExchangeRequest{
		GrantType:    "authorization_code",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Code:         code,
	}

	return s.doTokenRequest(ctx, payload)
}

// RefreshToken requests a new access token using an existing refresh token.
func (s *AuthService) RefreshToken(ctx context.Context, clientID, clientSecret, refreshToken string) (*Token, error) {
	if ctx == nil {
		return nil, errors.New("anilist: context must not be nil")
	}

	payload := tokenExchangeRequest{
		GrantType:    "refresh_token",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
	}

	return s.doTokenRequest(ctx, payload)
}

func (s *AuthService) doTokenRequest(ctx context.Context, payload tokenExchangeRequest) (*Token, error) {
	tokenURL := s.client.tokenURL
	if tokenURL == "" {
		tokenURL = DefaultTokenURL
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return nil, fmt.Errorf("anilist: encode token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, &buf)
	if err != nil {
		return nil, fmt.Errorf("anilist: create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if s.client.userAgent != "" {
		req.Header.Set("User-Agent", s.client.userAgent)
	}

	resp, err := s.client.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("anilist: token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errData struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errData)
		msg := errData.Message
		if msg == "" {
			msg = errData.Error
		}
		if msg == "" {
			msg = fmt.Sprintf("status code %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("anilist: token request rejected: %s", msg)
	}

	var token Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("anilist: decode token response: %w", err)
	}

	if token.ExpiresIn > 0 {
		token.Expiry = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	}

	return &token, nil
}