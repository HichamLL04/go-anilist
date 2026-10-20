package anilist

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient(nil)

	if c.client == nil {
		t.Fatal("expected non-nil http.Client")
	}
	if c.client.Timeout != DefaultTimeout {
		t.Errorf("Timeout = %v, want %v", c.client.Timeout, DefaultTimeout)
	}
	if c.baseURL != DefaultBaseURL {
		t.Errorf("baseURL = %q, want %q", c.baseURL, DefaultBaseURL)
	}
	if c.userAgent != DefaultUserAgent {
		t.Errorf("userAgent = %q, want %q", c.userAgent, DefaultUserAgent)
	}
	if c.token != "" {
		t.Errorf("token = %q, want empty", c.token)
	}
	if c.Media == nil || c.MediaList == nil || c.User == nil || c.Auth == nil {
		t.Error("expected all services to be initialized")
	}
}

func TestNewClient_Options(t *testing.T) {
	customHTTP := &http.Client{Timeout: 5 * time.Second}
	customBaseURL := "https://example.com/graphql"
	customToken := "test-secret-token"
	customUA := "custom-agent/1.0"

	c := NewClient(
		customHTTP,
		WithBaseURL(customBaseURL),
		WithToken(customToken),
		WithUserAgent(customUA),
	)

	if c.client != customHTTP {
		t.Errorf("client = %v, want %v", c.client, customHTTP)
	}
	if c.baseURL != customBaseURL {
		t.Errorf("baseURL = %q, want %q", c.baseURL, customBaseURL)
	}
	if c.token != customToken {
		t.Errorf("token = %q, want %q", c.token, customToken)
	}
	if c.userAgent != customUA {
		t.Errorf("userAgent = %q, want %q", c.userAgent, customUA)
	}
}
