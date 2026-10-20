package anilist

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type graphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type graphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors,omitempty"`
}

// doGraphQL sends a GraphQL request to AniList, reads response headers,
// streams the JSON response and decodes data into v.
func (c *Client) doGraphQL(ctx context.Context, query string, variables map[string]interface{}, v interface{}) (*Response, error) {
	if ctx == nil {
		return nil, errors.New("anilist: context must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("anilist: context cancelled: %w", err)
	}

	reqBody := graphQLRequest{
		Query:     query,
		Variables: variables,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(reqBody); err != nil {
		return nil, fmt.Errorf("anilist: failed to encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, &buf)
	if err != nil {
		return nil, fmt.Errorf("anilist: failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	httpResp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("anilist: request failed: %w", err)
	}
	defer httpResp.Body.Close()

	resp := newResponse(httpResp)

	var gqlResp graphQLResponse
	decErr := json.NewDecoder(httpResp.Body).Decode(&gqlResp)
	if decErr != nil && decErr != io.EOF {
		if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
			return resp, &ErrorResponse{
				Response: httpResp,
				Message:  fmt.Sprintf("http status %d: failed to decode response", httpResp.StatusCode),
			}
		}
		return resp, fmt.Errorf("anilist: failed to decode response: %w", decErr)
	}

	if len(gqlResp.Errors) > 0 {
		return resp, &ErrorResponse{
			Response: httpResp,
			Errors:   gqlResp.Errors,
		}
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return resp, &ErrorResponse{
			Response: httpResp,
		}
	}

	if v != nil && len(gqlResp.Data) > 0 {
		if err := json.Unmarshal(gqlResp.Data, v); err != nil {
			return resp, fmt.Errorf("anilist: failed to unmarshal data: %w", err)
		}
	}

	return resp, nil
}