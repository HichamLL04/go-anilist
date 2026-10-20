package anilist

import (
	"context"
	"fmt"
)

// MediaListStatus represents the status of a media entry in a user's list.
type MediaListStatus string

const (
	MediaListStatusCurrent   MediaListStatus = "CURRENT"
	MediaListStatusPlanning  MediaListStatus = "PLANNING"
	MediaListStatusCompleted MediaListStatus = "COMPLETED"
	MediaListStatusDropped   MediaListStatus = "DROPPED"
	MediaListStatusPaused    MediaListStatus = "PAUSED"
	MediaListStatusRepeating MediaListStatus = "REPEATING"
)

// MediaListEntry represents the status and progress of an anime in a user's list.
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

// SaveMediaListEntryInput provides parameters for updating or creating a list entry.
type SaveMediaListEntryInput struct {
	MediaID  int              `json:"mediaId"`
	Status   *MediaListStatus `json:"status,omitempty"`
	Progress *int             `json:"progress,omitempty"`
	Score    *float64         `json:"score,omitempty"`
	Repeat   *int             `json:"repeat,omitempty"`
}

type collectionParams struct {
	status MediaListStatus
}

// CollectionOption configures parameters for MediaListService.GetCollection.
type CollectionOption func(*collectionParams)

// WithStatus filters the collection by a specific list status.
func WithStatus(status MediaListStatus) CollectionOption {
	return func(p *collectionParams) {
		p.status = status
	}
}

const getMediaListEntryQuery = `query ($userId: Int, $mediaId: Int) {
  MediaList(userId: $userId, mediaId: $mediaId, type: ANIME) {
    id
    mediaId
    status
    progress
    score
    repeat
    updatedAt
    media {
      id
      idMal
      title {
        romaji
        english
        native
      }
      format
      status
      episodes
      duration
      synonyms
      seasonYear
    }
  }
}`

const getMediaListCollectionQuery = `query ($userId: Int, $status: MediaListStatus) {
  MediaListCollection(userId: $userId, type: ANIME, status: $status) {
    lists {
      name
      isCustomList
      status
      entries {
        id
        mediaId
        status
        progress
        score
        repeat
        updatedAt
        media {
          id
          idMal
          title {
            romaji
            english
            native
          }
          format
          status
          episodes
          duration
          synonyms
          seasonYear
        }
      }
    }
  }
}`

const saveMediaListEntryMutation = `mutation ($mediaId: Int, $status: MediaListStatus, $progress: Int, $score: Float, $repeat: Int) {
  SaveMediaListEntry(mediaId: $mediaId, status: $status, progress: $progress, score: $score, repeat: $repeat) {
    id
    mediaId
    status
    progress
    score
    repeat
    updatedAt
  }
}`

const deleteMediaListEntryMutation = `mutation ($id: Int) {
  DeleteMediaListEntry(id: $id) {
    deleted
  }
}`

// GetEntry retrieves the list entry for a specific user and anime.
// Returns ErrNotFound if the anime is not in the user's list.
func (s *MediaListService) GetEntry(ctx context.Context, userID, mediaID int) (*MediaListEntry, *Response, error) {
	var data struct {
		MediaList *MediaListEntry `json:"MediaList"`
	}

	vars := map[string]interface{}{
		"userId":  userID,
		"mediaId": mediaID,
	}

	resp, err := s.client.doGraphQL(ctx, getMediaListEntryQuery, vars, &data)
	if err != nil {
		return nil, resp, fmt.Errorf("anilist: get entry (user %d, media %d): %w", userID, mediaID, err)
	}

	if data.MediaList == nil {
		return nil, resp, fmt.Errorf("anilist: get entry (user %d, media %d): %w", userID, mediaID, ErrNotFound)
	}

	return data.MediaList, resp, nil
}

// GetCollection retrieves the user's anime list collection, optionally filtered by status.
func (s *MediaListService) GetCollection(ctx context.Context, userID int, opts ...CollectionOption) ([]MediaListEntry, *Response, error) {
	params := collectionParams{}
	for _, opt := range opts {
		opt(&params)
	}

	vars := map[string]interface{}{
		"userId": userID,
	}
	if params.status != "" {
		vars["status"] = params.status
	}

	var data struct {
		MediaListCollection struct {
			Lists []struct {
				Name    string           `json:"name"`
				Entries []MediaListEntry `json:"entries"`
			} `json:"lists"`
		} `json:"MediaListCollection"`
	}

	resp, err := s.client.doGraphQL(ctx, getMediaListCollectionQuery, vars, &data)
	if err != nil {
		return nil, resp, fmt.Errorf("anilist: get collection for user %d: %w", userID, err)
	}

	var entries []MediaListEntry
	for _, l := range data.MediaListCollection.Lists {
		entries = append(entries, l.Entries...)
	}

	return entries, resp, nil
}

// SaveEntry creates or updates an anime entry in the authenticated user's list.
func (s *MediaListService) SaveEntry(ctx context.Context, input SaveMediaListEntryInput) (*MediaListEntry, *Response, error) {
	vars := map[string]interface{}{
		"mediaId": input.MediaID,
	}
	if input.Status != nil {
		vars["status"] = *input.Status
	}
	if input.Progress != nil {
		vars["progress"] = *input.Progress
	}
	if input.Score != nil {
		vars["score"] = *input.Score
	}
	if input.Repeat != nil {
		vars["repeat"] = *input.Repeat
	}

	var data struct {
		SaveMediaListEntry *MediaListEntry `json:"SaveMediaListEntry"`
	}

	resp, err := s.client.doGraphQL(ctx, saveMediaListEntryMutation, vars, &data)
	if err != nil {
		return nil, resp, fmt.Errorf("anilist: save entry for media %d: %w", input.MediaID, err)
	}

	if data.SaveMediaListEntry == nil {
		return nil, resp, fmt.Errorf("anilist: save entry for media %d: unexpected empty result", input.MediaID)
	}

	return data.SaveMediaListEntry, resp, nil
}

// DeleteEntry deletes an entry from the authenticated user's list.
func (s *MediaListService) DeleteEntry(ctx context.Context, id int) (bool, *Response, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var data struct {
		DeleteMediaListEntry struct {
			Deleted bool `json:"deleted"`
		} `json:"DeleteMediaListEntry"`
	}

	resp, err := s.client.doGraphQL(ctx, deleteMediaListEntryMutation, vars, &data)
	if err != nil {
		return false, resp, fmt.Errorf("anilist: delete entry %d: %w", id, err)
	}

	return data.DeleteMediaListEntry.Deleted, resp, nil
}