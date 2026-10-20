package anilist

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMediaService_GetByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"Media":{"id":154587,"title":{"romaji":"Sousou no Frieren"},"format":"TV","episodes":28}}}`))
	}))
	defer server.Close()

	client := NewClient(nil, WithBaseURL(server.URL))

	media, _, err := client.Media.GetByID(context.Background(), 154587)
	if err != nil {
		t.Fatalf("GetByID() error: %v", err)
	}

	if media == nil {
		t.Fatal("expected non-nil Media")
	}
	if media.ID != 154587 {
		t.Errorf("ID = %d, want 154587", media.ID)
	}
	if media.Title.Romaji != "Sousou no Frieren" {
		t.Errorf("Title.Romaji = %q, want Sousou no Frieren", media.Title.Romaji)
	}
	if media.Episodes == nil || *media.Episodes != 28 {
		t.Errorf("Episodes = %v, want 28", media.Episodes)
	}
}

func TestMediaService_GetByIDMal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"Media":{"id":154587,"idMal":52991,"title":{"romaji":"Sousou no Frieren"}}}}`))
	}))
	defer server.Close()

	client := NewClient(nil, WithBaseURL(server.URL))

	media, _, err := client.Media.GetByIDMal(context.Background(), 52991)
	if err != nil {
		t.Fatalf("GetByIDMal() error: %v", err)
	}

	if media == nil {
		t.Fatal("expected non-nil Media")
	}
	if media.IDMal == nil || *media.IDMal != 52991 {
		t.Errorf("IDMal = %v, want 52991", media.IDMal)
	}
}

func TestMediaService_Search(t *testing.T) {
	fixtureData, err := os.ReadFile("testdata/media_search_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixtureData)
	}))
	defer server.Close()

	client := NewClient(nil, WithBaseURL(server.URL))

	results, _, err := client.Media.Search(
		context.Background(),
		"Frieren",
		WithFormat(MediaFormatTV),
		WithSeasonYear(2023),
		WithLimit(10),
	)
	if err != nil {
		t.Fatalf("Search() error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("results count = %d, want 1", len(results))
	}
	m := results[0]
	if m.ID != 154587 {
		t.Errorf("ID = %d, want 154587", m.ID)
	}
	if m.Title.English == nil || *m.Title.English != "Frieren: Beyond Journey's End" {
		t.Errorf("Title.English = %v, want Frieren: Beyond Journey's End", m.Title.English)
	}
	if m.SeasonYear == nil || *m.SeasonYear != 2023 {
		t.Errorf("SeasonYear = %v, want 2023", m.SeasonYear)
	}
}

func TestMediaService_GetByID_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"Media":null}}`))
	}))
	defer server.Close()

	client := NewClient(nil, WithBaseURL(server.URL))

	media, _, err := client.Media.GetByID(context.Background(), 999999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if media != nil {
		t.Errorf("expected nil media, got %v", media)
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}