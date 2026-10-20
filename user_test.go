package anilist

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUserService_CurrentUser(t *testing.T) {
	fixtureData, err := os.ReadFile("testdata/viewer_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-RateLimit-Remaining", "89")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixtureData)
	}))
	defer server.Close()

	client := NewClient(nil, WithBaseURL(server.URL), WithToken("dummy-token"))

	user, resp, err := client.User.CurrentUser(context.Background())
	if err != nil {
		t.Fatalf("CurrentUser() error: %v", err)
	}

	if user == nil {
		t.Fatal("expected non-nil User")
	}
	if user.ID != 98765 {
		t.Errorf("user.ID = %d, want 98765", user.ID)
	}
	if user.Name != "Hicham" {
		t.Errorf("user.Name = %q, want Hicham", user.Name)
	}
	if user.Avatar == nil || user.Avatar.Large == "" {
		t.Errorf("expected non-nil Avatar.Large, got %v", user.Avatar)
	}
	if user.MediaListOptions.ScoreFormat != "POINT_100" {
		t.Errorf("ScoreFormat = %q, want POINT_100", user.MediaListOptions.ScoreFormat)
	}
	if resp.RateLimit.Remaining != 89 {
		t.Errorf("RateLimit.Remaining = %d, want 89", resp.RateLimit.Remaining)
	}
}

func TestUserService_CurrentUser_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"errors":[{"message":"Invalid token","status":401}]}`))
	}))
	defer server.Close()

	client := NewClient(nil, WithBaseURL(server.URL))

	user, _, err := client.User.CurrentUser(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if user != nil {
		t.Errorf("expected nil user, got %v", user)
	}
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestUserService_GetByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"User":{"id":1234,"name":"Alice","mediaListOptions":{"scoreFormat":"POINT_10_DECIMAL"}}}}`))
	}))
	defer server.Close()

	client := NewClient(nil, WithBaseURL(server.URL))

	user, _, err := client.User.GetByID(context.Background(), 1234)
	if err != nil {
		t.Fatalf("GetByID() error: %v", err)
	}

	if user == nil {
		t.Fatal("expected non-nil User")
	}
	if user.ID != 1234 || user.Name != "Alice" {
		t.Errorf("user = %+v, want ID 1234 Name Alice", user)
	}
	if user.MediaListOptions.ScoreFormat != "POINT_10_DECIMAL" {
		t.Errorf("ScoreFormat = %q, want POINT_10_DECIMAL", user.MediaListOptions.ScoreFormat)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"User":null}}`))
	}))
	defer server.Close()

	client := NewClient(nil, WithBaseURL(server.URL))

	user, _, err := client.User.GetByID(context.Background(), 999999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if user != nil {
		t.Errorf("expected nil user, got %v", user)
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}