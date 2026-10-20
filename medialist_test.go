package anilist

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMediaListService_GetEntry(t *testing.T) {
	fixtureData, err := os.ReadFile("testdata/medialist_get_success.json")
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

	entry, _, err := client.MediaList.GetEntry(context.Background(), 98765, 154587)
	if err != nil {
		t.Fatalf("GetEntry() error: %v", err)
	}

	if entry == nil {
		t.Fatal("expected non-nil MediaListEntry")
	}
	if entry.ID != 12345678 {
		t.Errorf("ID = %d, want 12345678", entry.ID)
	}
	if entry.Progress != 14 {
		t.Errorf("Progress = %d, want 14", entry.Progress)
	}
	if entry.Status != MediaListStatusCurrent {
		t.Errorf("Status = %v, want CURRENT", entry.Status)
	}
	if entry.Score != 9.5 {
		t.Errorf("Score = %v, want 9.5", entry.Score)
	}
}

func TestMediaListService_GetEntry_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"MediaList":null}}`))
	}))
	defer server.Close()

	client := NewClient(nil, WithBaseURL(server.URL))

	entry, _, err := client.MediaList.GetEntry(context.Background(), 98765, 999999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if entry != nil {
		t.Errorf("expected nil entry, got %v", entry)
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMediaListService_GetCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"MediaListCollection": {
					"lists": [
						{
							"name": "Watching",
							"entries": [
								{"id": 1, "mediaId": 100, "status": "CURRENT", "progress": 5},
								{"id": 2, "mediaId": 200, "status": "CURRENT", "progress": 12}
							]
						}
					]
				}
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(nil, WithBaseURL(server.URL))

	entries, _, err := client.MediaList.GetCollection(context.Background(), 98765, WithStatus(MediaListStatusCurrent))
	if err != nil {
		t.Fatalf("GetCollection() error: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("entries count = %d, want 2", len(entries))
	}
	if entries[0].MediaID != 100 || entries[1].MediaID != 200 {
		t.Errorf("entries = %+v, want MediaIDs 100 and 200", entries)
	}
}

func TestMediaListService_SaveEntry(t *testing.T) {
	fixtureData, err := os.ReadFile("testdata/medialist_save_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixtureData)
	}))
	defer server.Close()

	client := NewClient(nil, WithBaseURL(server.URL), WithToken("test-token"))

	progress := 15
	status := MediaListStatusCurrent
	score := 10.0

	entry, _, err := client.MediaList.SaveEntry(context.Background(), SaveMediaListEntryInput{
		MediaID:  154587,
		Progress: &progress,
		Status:   &status,
		Score:    &score,
	})
	if err != nil {
		t.Fatalf("SaveEntry() error: %v", err)
	}

	if entry == nil {
		t.Fatal("expected non-nil entry")
	}
	if entry.Progress != 15 {
		t.Errorf("Progress = %d, want 15", entry.Progress)
	}
	if entry.Score != 10.0 {
		t.Errorf("Score = %v, want 10.0", entry.Score)
	}
}

func TestMediaListService_DeleteEntry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"DeleteMediaListEntry":{"deleted":true}}}`))
	}))
	defer server.Close()

	client := NewClient(nil, WithBaseURL(server.URL), WithToken("test-token"))

	deleted, _, err := client.MediaList.DeleteEntry(context.Background(), 12345678)
	if err != nil {
		t.Fatalf("DeleteEntry() error: %v", err)
	}
	if !deleted {
		t.Error("expected deleted = true, got false")
	}
}