package anilist

import (
	"context"
	"fmt"
)

// User represents a user profile in AniList.
type User struct {
	ID               int              `json:"id"`
	Name             string           `json:"name"`
	Avatar           *UserAvatar      `json:"avatar,omitempty"`
	MediaListOptions MediaListOptions `json:"mediaListOptions"`
}

// UserAvatar holds URLs for the user's avatar.
type UserAvatar struct {
	Large  string `json:"large,omitempty"`
	Medium string `json:"medium,omitempty"`
}

// MediaListOptions contains user preferences for media lists.
type MediaListOptions struct {
	ScoreFormat string `json:"scoreFormat"` // e.g. "POINT_100", "POINT_10_DECIMAL", "POINT_10", "POINT_5", "POINT_3"
}

const currentUserQuery = `query CurrentUser {
  Viewer {
    id
    name
    avatar {
      large
      medium
    }
    mediaListOptions {
      scoreFormat
    }
  }
}`

const getUserByIDQuery = `query GetUser($id: Int) {
  User(id: $id) {
    id
    name
    avatar {
      large
      medium
    }
    mediaListOptions {
      scoreFormat
    }
  }
}`

// CurrentUser retrieves the authenticated user's profile and list settings (Viewer query).
func (s *UserService) CurrentUser(ctx context.Context) (*User, *Response, error) {
	var data struct {
		Viewer *User `json:"Viewer"`
	}

	resp, err := s.client.doGraphQL(ctx, currentUserQuery, nil, &data)
	if err != nil {
		return nil, resp, fmt.Errorf("anilist: get current user: %w", err)
	}

	if data.Viewer == nil {
		return nil, resp, fmt.Errorf("anilist: get current user: %w", ErrNotFound)
	}

	return data.Viewer, resp, nil
}

// GetByID retrieves a user's public profile by their numeric AniList ID.
func (s *UserService) GetByID(ctx context.Context, id int) (*User, *Response, error) {
	var data struct {
		User *User `json:"User"`
	}

	vars := map[string]interface{}{
		"id": id,
	}

	resp, err := s.client.doGraphQL(ctx, getUserByIDQuery, vars, &data)
	if err != nil {
		return nil, resp, fmt.Errorf("anilist: get user by id %d: %w", id, err)
	}

	if data.User == nil {
		return nil, resp, fmt.Errorf("anilist: get user by id %d: %w", id, ErrNotFound)
	}

	return data.User, resp, nil
}