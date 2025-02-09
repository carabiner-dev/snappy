// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
)

type Options struct {
	Host  string
	Token string
}

var defaultOptions = Options{
	Host: "api.github.com",
}

type Caller interface {
	RequestWithContext(context.Context, string, string, io.Reader) (*http.Response, error)
}

func buildGithubRestClient(opts Options) (*api.RESTClient, error) {
	return api.NewRESTClient(api.ClientOptions{
		AuthToken: opts.Token,
		Host:      opts.Host,
	})
}

func readToken(opts *Options) {
	opts.Token = os.Getenv("GITHUB_TOKEN")
}

func NewClient() (*Client, error) {
	return NewClientWithOptions(defaultOptions)
}

func NewClientWithOptions(opts Options) (*Client, error) {
	if opts.Token == "" {
		readToken(&defaultOptions)
	}
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
