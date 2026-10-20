package anilist

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	// ErrNotFound is returned when a requested resource (Media, MediaListEntry, User) does not exist.
	ErrNotFound = errors.New("anilist: resource not found")

	// ErrUnauthorized is returned when the request lacks valid authentication credentials.
	ErrUnauthorized = errors.New("anilist: unauthorized or invalid token")

	// ErrRateLimited is returned when the AniList rate limit (90 req/min) is exceeded.
	ErrRateLimited = errors.New("anilist: rate limit exceeded (429)")
)

// Location points to the line and column in the GraphQL query where an error originated.
type Location struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// GraphQLError represents an individual error returned in the GraphQL response "errors" array.
type GraphQLError struct {
	Message   string     `json:"message"`
	Status    int        `json:"status"`
	Locations []Location `json:"locations,omitempty"`
}

func (e GraphQLError) Error() string {
	if e.Status != 0 {
		return fmt.Sprintf("graphql error (%d): %s", e.Status, e.Message)
	}
	return fmt.Sprintf("graphql error: %s", e.Message)
}

// ErrorResponse wraps HTTP and GraphQL errors returned by the AniList API.
type ErrorResponse struct {
	Response *http.Response `json:"-"`
	Errors   []GraphQLError `json:"errors,omitempty"`
	Message  string         `json:"message,omitempty"`
}

func (e *ErrorResponse) Error() string {
	code := 0
	if e.Response != nil {
		code = e.Response.StatusCode
	}

	if len(e.Errors) > 0 {
		msgs := make([]string, len(e.Errors))
		for i, err := range e.Errors {
			msgs[i] = err.Message
		}
		if code > 0 {
			return fmt.Sprintf("anilist: api error (%d): %s", code, strings.Join(msgs, "; "))
		}
		return fmt.Sprintf("anilist: api error: %s", strings.Join(msgs, "; "))
	}

	if e.Message != "" {
		if code > 0 {
			return fmt.Sprintf("anilist: api error (%d): %s", code, e.Message)
		}
		return fmt.Sprintf("anilist: api error: %s", e.Message)
	}

	if code > 0 {
		return fmt.Sprintf("anilist: api error (%d)", code)
	}
	return "anilist: unknown api error"
}

// Is reports whether any error in ErrorResponse matches target centinel errors.
func (e *ErrorResponse) Is(target error) bool {
	if e.Response != nil {
		if target == ErrRateLimited && e.Response.StatusCode == http.StatusTooManyRequests {
			return true
		}
		if target == ErrUnauthorized && e.Response.StatusCode == http.StatusUnauthorized {
			return true
		}
		if target == ErrNotFound && e.Response.StatusCode == http.StatusNotFound {
			return true
		}
	}

	for _, ge := range e.Errors {
		if target == ErrNotFound && ge.Status == http.StatusNotFound {
			return true
		}
		if target == ErrRateLimited && ge.Status == http.StatusTooManyRequests {
			return true
		}
		if target == ErrUnauthorized && ge.Status == http.StatusUnauthorized {
			return true
		}
	}

	return false
}