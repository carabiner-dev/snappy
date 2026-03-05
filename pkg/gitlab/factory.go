// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package gitlab

import "github.com/carabiner-dev/snappy/pkg/platform"

type Factory struct{}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) CreateClient() (platform.Client, error) {
	return NewClient()
}

func (f *Factory) Platform() platform.Type {
	return platform.GitLab
}

func (f *Factory) DefaultResponseHeaders() []string {
	return []string{
		"Date",
		"Etag",
		"X-Request-Id",
		"X-Runtime",
		"RateLimit-Limit",
		"RateLimit-Observed",
		"RateLimit-Remaining",
		"RateLimit-Reset",
		"RateLimit-ResetTime",
	}
}
