// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package github

import "github.com/carabiner-dev/snappy/pkg/platform"

func init() {
	platform.Register(NewFactory())
}

type Factory struct{}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) CreateClient() (platform.Client, error) {
	return NewClient()
}

func (f *Factory) Platform() platform.Type {
	return platform.GitHub
}

func (f *Factory) DefaultResponseHeaders() []string {
	return []string{
		"Date",
		"Etag",
		"Server",
		"X-Accepted-Oauth-Scopes",
		"X-Github-Api-Version-Selected",
		"X-Github-Request-Id",
	}
}
