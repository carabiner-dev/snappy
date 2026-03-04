// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package gitlab

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/carabiner-dev/snappy/pkg/platform"
)

type Client struct {
	Options    Options
	httpClient *http.Client
	baseURL    string
}

func NewClient() (*Client, error) {
	return NewClientWithOptions(defaultOptions)
}

func NewClientWithOptions(opts Options) (*Client, error) {
	// Ensure the client has a token to connect
	if err := opts.ensureToken(); err != nil {
		return nil, err
	}

	// Construct base URL
	baseURL := fmt.Sprintf("https://%s/api/v4", opts.Host)

	return &Client{
		Options:    opts,
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}, nil
}

func (c *Client) Call(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	// Construct full URL
	url := fmt.Sprintf("%s/%s", c.baseURL, path)

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Private-Token", c.Options.Token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Make request
	resp, err := c.httpClient.Do(req) //nolint:gosec // G704: Expected external input
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}

	return resp, nil
}

func (c *Client) Platform() platform.Type {
	return platform.GitLab
}
