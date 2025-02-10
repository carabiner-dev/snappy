// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"io"
	"net/http"

	"github.com/cli/go-gh/v2/pkg/api"
)

// Replaceable caller interface
type Caller interface {
	RequestWithContext(context.Context, string, string, io.Reader) (*http.Response, error)
}

func buildGithubRestClient(opts Options) (*api.RESTClient, error) {
	return api.NewRESTClient(api.ClientOptions{
		AuthToken: opts.Token,
		Host:      opts.Host,
	})
}

func NewClient() (*Client, error) {
	return NewClientWithOptions(defaultOptions)
}

func NewClientWithOptions(opts Options) (*Client, error) {
	// Ensure the client has a token to connect
	if err := opts.ensureToken(); err != nil {
		return nil, err
	}

	// Create the client
	rclient, err := buildGithubRestClient(opts)
	if err != nil {
		return nil, err
	}
	return &Client{
		Options: opts,
		caller:  rclient,
	}, nil
}

type Client struct {
	Options Options
	caller  Caller
}

func (c *Client) Call(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	return c.caller.RequestWithContext(ctx, method, path, body)
}
