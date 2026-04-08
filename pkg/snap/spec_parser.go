// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package snap

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	yaml "github.com/carabiner-dev/yamplate"
)

type SpecParser struct{}

type ParseOptions struct {
	Variables map[string]string
}

// ParseWithOptions parses a spec yaml definition.
func (sp *SpecParser) ParseWithOptions(r io.Reader, opts *ParseOptions) (*Spec, error) {
	// Build the parser and
	decoder := yaml.NewDecoder(r)

	// pass the substitution table
	decoder.Options.Variables = opts.Variables

	spec := &Spec{}
	if err := decoder.Decode(spec); err != nil {
		return nil, fmt.Errorf("parsing yaml spec: %w", err)
	}

	// Default to payloads of type struct
	if spec.PayloadType == "" {
		spec.PayloadType = PayloadTypeStruct
	}

	// Post-process GitLab endpoints to URL-encode project paths
	// This happens AFTER variable substitution so we can encode the full path
	spec.Endpoint = encodeGitLabEndpoint(spec.Endpoint)

	return spec, nil
}

// encodeGitLabEndpoint URL-encodes the project path portion of GitLab API endpoints.
// Handles endpoints like "projects/{group}/{project}/resource/..." where the group/project
// path may contain slashes that need to be encoded as %2F for the API call.
//
// This runs after variable substitution, so it works on the final endpoint string.
// Detection strategy: the project path ends when we encounter either:
// 1. A segment containing underscore (snake_case API resources like protected_branches)
// 2. A numeric segment (resource IDs)
// 3. A known unambiguous API resource name
func encodeGitLabEndpoint(endpoint string) string {
	// Only process GitLab projects endpoints
	if !strings.HasPrefix(endpoint, "projects/") {
		return endpoint
	}

	// Remove prefix and check if already encoded
	rest := strings.TrimPrefix(endpoint, "projects/")
	if strings.Contains(rest, "%") {
		// Assume already encoded, don't double-encode
		return endpoint
	}

	// Split into segments
	segments := strings.Split(rest, "/")
	if len(segments) == 0 {
		return endpoint
	}

	// Find where project path ends
	projectEndIndex := findProjectPathEnd(segments)

	// Extract and encode project path
	projectPath := strings.Join(segments[:projectEndIndex], "/")
	encodedPath := url.PathEscape(projectPath)

	// Reconstruct endpoint
	if projectEndIndex < len(segments) {
		remaining := strings.Join(segments[projectEndIndex:], "/")
		return "projects/" + encodedPath + "/" + remaining
	}

	return "projects/" + encodedPath
}

// findProjectPathEnd determines where the GitLab project path ends in the segment list.
// Returns the index of the first segment that is definitely an API resource (not part of project path).
func findProjectPathEnd(segments []string) int {
	for i, segment := range segments {
		if isDefinitelyGitLabResource(segment) {
			return i
		}
	}
	// If no resource found, entire path is the project path
	return len(segments)
}

// isDefinitelyGitLabResource checks if a segment is definitely a GitLab API resource.
// Uses high-confidence heuristics to minimize false positives.
func isDefinitelyGitLabResource(segment string) bool {
	// Numeric IDs are definitely resource identifiers, not project names
	if isNumeric(segment) {
		return true
	}

	// Snake_case (with underscores) is used for API resources, rarely for project names
	if strings.Contains(segment, "_") {
		return true
	}

	// Unambiguous single-word API resources that are very unlikely to be project names
	unambiguousResources := []string{
		"repository", "branches", "commits", "tags", "tree",
		"members", "variables", "pipelines", "jobs", "artifacts",
		"runners", "hooks", "labels", "milestones", "releases",
		"approvals",
	}

	for _, resource := range unambiguousResources {
		if segment == resource {
			return true
		}
	}

	return false
}

// isNumeric checks if a string contains only digits
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
