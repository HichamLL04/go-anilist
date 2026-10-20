package anilist

import (
	"net/http"
	"strconv"
	"time"
)

// Response wraps http.Response and includes rate limit metadata returned by AniList.
type Response struct {
	*http.Response
	RateLimit RateLimitInfo
}

// RateLimitInfo contains the rate limiting details extracted from AniList HTTP headers.
type RateLimitInfo struct {
	Limit      int           // Maximum allowed requests in current window (X-RateLimit-Limit)
	Remaining  int           // Remaining requests in current window (X-RateLimit-Remaining)
	Reset      time.Time     // Time at which the current rate limit window resets (X-RateLimit-Reset)
	RetryAfter time.Duration // Recommended wait time when rate limited (Retry-After header)
}

func newResponse(r *http.Response) *Response {
	resp := &Response{Response: r}
	if r == nil {
		return resp
	}

	if limitStr := r.Header.Get("X-RateLimit-Limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			resp.RateLimit.Limit = limit
		}
	}

	if remStr := r.Header.Get("X-RateLimit-Remaining"); remStr != "" {
		if rem, err := strconv.Atoi(remStr); err == nil {
			resp.RateLimit.Remaining = rem
		}
	}

	if resetStr := r.Header.Get("X-RateLimit-Reset"); resetStr != "" {
		if resetUnix, err := strconv.ParseInt(resetStr, 10, 64); err == nil {
			resp.RateLimit.Reset = time.Unix(resetUnix, 0)
		}
	}

	if retryStr := r.Header.Get("Retry-After"); retryStr != "" {
		if seconds, err := strconv.Atoi(retryStr); err == nil {
			resp.RateLimit.RetryAfter = time.Duration(seconds) * time.Second
		}
	}

	return resp
}