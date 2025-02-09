// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package snap

import (
	"context"
	"fmt"
)

func New() *Snapper {
	return &Snapper{
		Options:        defaultOptions,
		implementation: &defaultImplementation{},
	}
}

var defaultOptions = Options{
	ResponseHeaders: []string{
		"Date", "Etag", "Server", "X-Accepted-Oauth-Scopes",
		"X-Github-Api-Version-Selected", "X-Github-Request-Id",
	},
}

type Options struct {
	// ResponseHeaders is the list of response headers the snapper will record
	ResponseHeaders []string
}

type Snapper struct {
	Options        Options
	implementation SnapperImplementation
}

// Take grabs a snapshot of the repo status
func (s *Snapper) Take(ctx context.Context, spec *Spec) (*Snapshot, error) {
	if err := s.implementation.ValidateSpec(spec); err != nil {
		return nil, fmt.Errorf("validating spec: %w", err)
	}
	client, err := s.implementation.GetClient()
	if err != nil {
		return nil, fmt.Errorf("creating api client: %w", err)
	}

	resp, err := s.implementation.CallAPI(ctx, client, spec)
	if err != nil {
		return nil, fmt.Errorf("calling API: %w", err)
	}
	defer resp.Body.Close()

	snapshot, err := s.implementation.ParseResponse(&s.Options, spec, resp)
	if err != nil {
		return nil, fmt.Errorf("parsing API response: %w", err)
	}

	snapshot, err = s.implementation.ApplyFieldMask(snapshot, spec.Mask)
	if err != nil {
		return nil, fmt.Errorf("applying field mask: %w", err)
	}

	return snapshot, nil
}
