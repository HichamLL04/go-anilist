package anilist

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTransport_Success(t *testing.T) {
	wantQuery := "query { Viewer { id name } }"
	wantVars := map[string]interface{}{"id": float64(123)}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %q, want POST", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("Authorization") != "Bearer test-secret" {
			t.Errorf("Authorization = %q, want Bearer test-secret", r.Header.Get("Authorization"))
		}

		var reqBody graphQLRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if reqBody.Query != wantQuery {
			t.Errorf("query = %q, want %q", reqBody.Query, wantQuery)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-RateLimit-Limit", "90")
		w.Header().Set("X-RateLimit-Remaining", "88")
		w.Header().Set("X-RateLimit-Reset", "1700000000")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"Viewer":{"id":123,"name":"TestUser"}}}`))
	}))
	defer server.Close()

	client := NewClient(nil, WithBaseURL(server.URL), WithToken("test-secret"))

	var result struct {
		Viewer struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"Viewer"`
	}

	resp, err := client.doGraphQL(context.Background(), wantQuery, wantVars, &result)
	if err != nil {
		t.Fatalf("doGraphQL() unexpected error: %v", err)
	}

	if resp == nil {
		t.Fatal("expected non-nil Response")
	}
	if resp.RateLimit.Limit != 90 {
		t.Errorf("RateLimit.Limit = %d, want 90", resp.RateLimit.Limit)
	}
	if resp.RateLimit.Remaining != 88 {
		t.Errorf("RateLimit.Remaining = %d, want 88", resp.RateLimit.Remaining)
	}
	if resp.RateLimit.Reset.Unix() != 1700000000 {
		t.Errorf("RateLimit.Reset = %v, want unix 1700000000", resp.RateLimit.Reset)
	}
	if result.Viewer.ID != 123 || result.Viewer.Name != "TestUser" {
		t.Errorf("result = %+v, want ID 123 Name TestUser", result)
	}
}

func TestTransport_GraphQLError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"errors":[{"message":"Validation failed","status":400,"locations":[{"line":1,"column":2}]}]}`))
	}))
	defer server.Close()

	client := NewClient(nil, WithBaseURL(server.URL))

	var result map[string]interface{}
	_, err := client.doGraphQL(context.Background(), "query { invalid }", nil, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var errResp *ErrorResponse
	if !errors.As(err, &errResp) {
		t.Fatalf("expected *ErrorResponse, got %T: %v", err, err)
	}
	if len(errResp.Errors) != 1 {
		t.Fatalf("got %d errors, want 1", len(errResp.Errors))
	}
	if errResp.Errors[0].Message != "Validation failed" {
		t.Errorf("error message = %q, want Validation failed", errResp.Errors[0].Message)
	}
}

func TestTransport_RateLimitError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"errors":[{"message":"Too Many Requests","status":429}]}`))
	}))
	defer server.Close()

	client := NewClient(nil, WithBaseURL(server.URL))

	resp, err := client.doGraphQL(context.Background(), "query { test }", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrRateLimited) {
		t.Errorf("errors.Is(err, ErrRateLimited) = false, want true")
	}
	if resp.RateLimit.RetryAfter != 30*time.Second {
		t.Errorf("RetryAfter = %v, want 30s", resp.RateLimit.RetryAfter)
	}
}

func TestTransport_ContextCancellation(t *testing.T) {
	client := NewClient(nil)

	// Test nil context
	_, err := client.doGraphQL(nil, "query { test }", nil, nil) //nolint:staticcheck
	if err == nil {
		t.Error("expected error for nil context, got nil")
	}

	// Test pre-cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = client.doGraphQL(ctx, "query { test }", nil, nil)
	if err == nil {
		t.Error("expected error for cancelled context, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}