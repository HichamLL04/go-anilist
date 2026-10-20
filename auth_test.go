package anilist

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAuthService_AuthorizeURL(t *testing.T) {
	client := NewClient(nil)

	authURL := client.Auth.AuthorizeURL("client123", "https://example.com/callback", "xyzstate")

	u, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("failed to parse authURL: %v", err)
	}

	if !strings.HasPrefix(authURL, DefaultAuthURL) {
		t.Errorf("expected URL to start with %q, got %q", DefaultAuthURL, authURL)
	}

	q := u.Query()
	if q.Get("client_id") != "client123" {
		t.Errorf("client_id = %q, want client123", q.Get("client_id"))
	}
	if q.Get("response_type") != "code" {
		t.Errorf("response_type = %q, want code", q.Get("response_type"))
	}
	if q.Get("redirect_uri") != "https://example.com/callback" {
		t.Errorf("redirect_uri = %q, want https://example.com/callback", q.Get("redirect_uri"))
	}
	if q.Get("state") != "xyzstate" {
		t.Errorf("state = %q, want xyzstate", q.Get("state"))
	}
}

func TestAuthService_ExchangeCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %q, want POST", r.Method)
		}

		var req tokenExchangeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.GrantType != "authorization_code" {
			t.Errorf("GrantType = %q, want authorization_code", req.GrantType)
		}
		if req.Code != "authcode123" {
			t.Errorf("Code = %q, want authcode123", req.Code)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"access_token": "mock-access-token",
			"token_type": "Bearer",
			"expires_in": 31536000,
			"refresh_token": "mock-refresh-token"
		}`))
	}))
	defer server.Close()

	client := NewClient(nil, WithTokenURL(server.URL))

	token, err := client.Auth.ExchangeCode(
		context.Background(),
		"client123",
		"secret456",
		"https://example.com/callback",
		"authcode123",
	)
	if err != nil {
		t.Fatalf("ExchangeCode() error: %v", err)
	}

	if token == nil {
		t.Fatal("expected non-nil Token")
	}
	if token.AccessToken != "mock-access-token" {
		t.Errorf("AccessToken = %q, want mock-access-token", token.AccessToken)
	}
	if token.RefreshToken != "mock-refresh-token" {
		t.Errorf("RefreshToken = %q, want mock-refresh-token", token.RefreshToken)
	}
	if token.ExpiresIn != 31536000 {
		t.Errorf("ExpiresIn = %d, want 31536000", token.ExpiresIn)
	}
	if token.Expiry.IsZero() {
		t.Error("expected non-zero Expiry time")
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req tokenExchangeRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		if req.GrantType != "refresh_token" {
			t.Errorf("GrantType = %q, want refresh_token", req.GrantType)
		}
		if req.RefreshToken != "old-refresh-token" {
			t.Errorf("RefreshToken = %q, want old-refresh-token", req.RefreshToken)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"access_token": "new-access-token",
			"token_type": "Bearer",
			"expires_in": 31536000,
			"refresh_token": "new-refresh-token"
		}`))
	}))
	defer server.Close()

	client := NewClient(nil, WithTokenURL(server.URL))

	token, err := client.Auth.RefreshToken(
		context.Background(),
		"client123",
		"secret456",
		"old-refresh-token",
	)
	if err != nil {
		t.Fatalf("RefreshToken() error: %v", err)
	}

	if token.AccessToken != "new-access-token" {
		t.Errorf("AccessToken = %q, want new-access-token", token.AccessToken)
	}
}

func TestAuthService_ExchangeCode_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid_grant","message":"Authorization code expired"}`))
	}))
	defer server.Close()

	client := NewClient(nil, WithTokenURL(server.URL))

	token, err := client.Auth.ExchangeCode(
		context.Background(),
		"client123",
		"secret456",
		"https://example.com/callback",
		"badcode",
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if token != nil {
		t.Errorf("expected nil token, got %v", token)
	}
	if !strings.Contains(err.Error(), "Authorization code expired") {
		t.Errorf("error = %q, want containing 'Authorization code expired'", err.Error())
	}
}