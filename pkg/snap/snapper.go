// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package snap

import (
	"context"
	"fmt"

	"github.com/carabiner-dev/snappy/pkg/platform"
)

// New creates a new Snapper
func New() *Snapper {
	return &Snapper{
		Options:        Options{},
		implementation: &defaultImplementation{},
		detector:       platform.NewDetector(),
	}
}

type Options struct {
	// ResponseHeaders is the list of response headers the snapper will record
	ResponseHeaders []string

	// Platform is the explicitly selected platform (optional, auto-detected if not set)
	Platform platform.Type

	// SpecPath is the path to the spec file (used for auto-detection)
	SpecPath string

	// Client is a pre-built API client the snapper uses instead of creating
	// one through the platform registry. Set it when the credentials come from
	// somewhere other than the process environment — for example a server
	// seeding each call with a short-lived installation token. When set, its
	// platform also resolves the snapshot defaults unless Platform is set
	// explicitly.
	Client platform.Client
}

type Snapper struct {
	Options        Options
	implementation SnapperImplementation
	detector       *platform.Detector
}

// Take grabs a snapshot of the repo status
func (s *Snapper) Take(ctx context.Context, spec *Spec) (*Snapshot, error) {
	if err := s.implementation.ValidateSpec(spec); err != nil {
		return nil, fmt.Errorf("validating spec: %w", err)
	}

	// Determine platform type
	platformType := s.Options.Platform
	if platformType == "" && s.Options.Client != nil {
		// A pre-built client knows its platform.
		platformType = s.Options.Client.Platform()
	}
	if platformType == "" {
		// Try to detect from spec path first
		if s.Options.SpecPath != "" {
			detected, err := s.detector.DetectFromSpec(s.Options.SpecPath)
			if err == nil {
				platformType = detected
			}
		}
		// Fallback to detecting from endpoint
		if platformType == "" {
			detected, err := s.detector.DetectFromEndpoint(spec.Endpoint)
			if err != nil {
				return nil, fmt.Errorf("could not auto-detect platform, please specify --platform flag: %w", err)
			}
			platformType = detected
		}
	}

	// Get platform factory and set default headers if not configured
	factory, err := platform.Get(platformType)
	if err != nil {
		return nil, fmt.Errorf("getting platform factory: %w", err)
	}
	if len(s.Options.ResponseHeaders) == 0 {
		s.Options.ResponseHeaders = factory.DefaultResponseHeaders()
	}

	// Use the pre-built client or create one for the detected/specified platform
	client := s.Options.Client
	if client == nil {
		client, err = s.implementation.GetClient(platformType)
		if err != nil {
			return nil, fmt.Errorf("creating api client: %w", err)
		}
	}

	resp, err := s.implementation.CallAPI(ctx, client, spec)
	if err != nil {
		return nil, fmt.Errorf("calling API: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	snapshot, err := s.implementation.ParseResponse(&s.Options, spec, resp)
	if err != nil {
		return nil, fmt.Errorf("parsing API response: %w", err)
	}

	if spec.PayloadType == PayloadTypeStruct {
		snapshot, err = s.implementation.ApplyFieldMask(snapshot, spec.Mask)
		if err != nil {
			return nil, fmt.Errorf("applying field mask: %w", err)
		}
	}

	return snapshot, nil
}
