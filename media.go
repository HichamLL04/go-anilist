package anilist

import (
	"context"
	"fmt"
)

// MediaType represents if the media is anime or manga.
type MediaType string

const (
	MediaTypeAnime MediaType = "ANIME"
	MediaTypeManga MediaType = "MANGA"
)

// MediaFormat represents the broadcast or release format.
type MediaFormat string

const (
	MediaFormatTV      MediaFormat = "TV"
	MediaFormatTVShort MediaFormat = "TV_SHORT"
	MediaFormatMovie   MediaFormat = "MOVIE"
	MediaFormatSpecial MediaFormat = "SPECIAL"
	MediaFormatOVA     MediaFormat = "OVA"
	MediaFormatONA     MediaFormat = "ONA"
	MediaFormatMusic   MediaFormat = "MUSIC"
)

// MediaStatus represents the current releasing status.
type MediaStatus string

const (
	MediaStatusFinished       MediaStatus = "FINISHED"
	MediaStatusReleasing      MediaStatus = "RELEASING"
	MediaStatusNotYetReleased MediaStatus = "NOT_YET_RELEASED"
	MediaStatusCancelled      MediaStatus = "CANCELLED"
	MediaStatusHiatus         MediaStatus = "HIATUS"
)

// MediaTitle contains titles in multiple language scripts.
type MediaTitle struct {
	Romaji  string  `json:"romaji"`
	English *string `json:"english,omitempty"`
	Native  *string `json:"native,omitempty"`
}

// Media represents an anime or manga entry in AniList.
type Media struct {
	ID         int         `json:"id"`
	IDMal      *int        `json:"idMal,omitempty"`
	Title      MediaTitle  `json:"title"`
	Format     MediaFormat `json:"format"`
	Status     MediaStatus `json:"status"`
	Episodes   *int        `json:"episodes,omitempty"`
	Duration   *int        `json:"duration,omitempty"`
	Synonyms   []string    `json:"synonyms"`
	SeasonYear *int        `json:"seasonYear,omitempty"`
}

const mediaFieldsFragment = `
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
`

const getMediaByIDQuery = `query ($id: Int) {
  Media(id: $id, type: ANIME) {` + mediaFieldsFragment + `}
}`

const getMediaByIDMalQuery = `query ($idMal: Int) {
  Media(idMal: $idMal, type: ANIME) {` + mediaFieldsFragment + `}
}`

const searchMediaQuery = `query ($search: String, $format: MediaFormat, $seasonYear: Int, $perPage: Int) {
  Page(page: 1, perPage: $perPage) {
    media(search: $search, format: $format, seasonYear: $seasonYear, type: ANIME) {` + mediaFieldsFragment + `}
  }
}`

// GetByID retrieves an anime by its AniList numeric ID.
func (s *MediaService) GetByID(ctx context.Context, id int) (*Media, *Response, error) {
	var data struct {
		Media *Media `json:"Media"`
	}

	vars := map[string]interface{}{
		"id": id,
	}

	resp, err := s.client.doGraphQL(ctx, getMediaByIDQuery, vars, &data)
	if err != nil {
		return nil, resp, fmt.Errorf("anilist: get media by id %d: %w", id, err)
	}

	if data.Media == nil {
		return nil, resp, fmt.Errorf("anilist: get media by id %d: %w", id, ErrNotFound)
	}

	return data.Media, resp, nil
}

// GetByIDMal retrieves an anime by its MyAnimeList numeric ID.
func (s *MediaService) GetByIDMal(ctx context.Context, idMal int) (*Media, *Response, error) {
	var data struct {
		Media *Media `json:"Media"`
	}

	vars := map[string]interface{}{
		"idMal": idMal,
	}

	resp, err := s.client.doGraphQL(ctx, getMediaByIDMalQuery, vars, &data)
	if err != nil {
		return nil, resp, fmt.Errorf("anilist: get media by mal id %d: %w", idMal, err)
	}

	if data.Media == nil {
		return nil, resp, fmt.Errorf("anilist: get media by mal id %d: %w", idMal, ErrNotFound)
	}

	return data.Media, resp, nil
}

// Search queries AniList for anime matching the given text and options.
func (s *MediaService) Search(ctx context.Context, query string, opts ...SearchOption) ([]Media, *Response, error) {
	params := searchParams{
		limit: 10,
	}
	for _, opt := range opts {
		opt(&params)
	}

	if params.limit > 50 {
		params.limit = 50
	}

	vars := map[string]interface{}{
		"search":  query,
		"perPage": params.limit,
	}
	if params.format != "" {
		vars["format"] = params.format
	}
	if params.seasonYear != 0 {
		vars["seasonYear"] = params.seasonYear
	}

	var data struct {
		Page struct {
			Media []Media `json:"media"`
		} `json:"Page"`
	}

	resp, err := s.client.doGraphQL(ctx, searchMediaQuery, vars, &data)
	if err != nil {
		return nil, resp, fmt.Errorf("anilist: search media %q: %w", query, err)
	}

	return data.Page.Media, resp, nil
}