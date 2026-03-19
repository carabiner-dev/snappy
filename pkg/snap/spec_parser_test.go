// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package snap

import (
	"testing"
)

func TestEncodeGitLabEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		expected string
	}{
		{
			name:     "simple project path with resource",
			endpoint: "projects/mygroup/myproject/protected_branches/main",
			expected: "projects/mygroup%2Fmyproject/protected_branches/main",
		},
		{
			name:     "nested group path",
			endpoint: "projects/my-org/dev-team/my-project/protected_branches/main",
			expected: "projects/my-org%2Fdev-team%2Fmy-project/protected_branches/main",
		},
		{
			name:     "project endpoint only",
			endpoint: "projects/mygroup/myproject",
			expected: "projects/mygroup%2Fmyproject",
		},
		{
			name:     "repository endpoint",
			endpoint: "projects/group/subgroup/project/repository/branches",
			expected: "projects/group%2Fsubgroup%2Fproject/repository/branches",
		},
		{
			name:     "merge requests endpoint with ID",
			endpoint: "projects/group/project/merge_requests/123",
			expected: "projects/group%2Fproject/merge_requests/123",
		},
		{
			name:     "pipelines with numeric ID",
			endpoint: "projects/org/team/proj/pipelines/456",
			expected: "projects/org%2Fteam%2Fproj/pipelines/456",
		},
		{
			name:     "non-gitlab endpoint unchanged",
			endpoint: "repos/owner/repo/branches/main",
			expected: "repos/owner/repo/branches/main",
		},
		{
			name:     "already encoded path unchanged",
			endpoint: "projects/group%2Fproject/protected_branches/main",
			expected: "projects/group%2Fproject/protected_branches/main",
		},
		{
			name:     "members endpoint",
			endpoint: "projects/a/b/c/members",
			expected: "projects/a%2Fb%2Fc/members",
		},
		{
			name:     "variables endpoint",
			endpoint: "projects/group/project/variables/VAR_NAME",
			expected: "projects/group%2Fproject/variables/VAR_NAME",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encodeGitLabEndpoint(tt.endpoint)
			if result != tt.expected {
				t.Errorf("encodeGitLabEndpoint() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"0", true},
		{"456789", true},
		{"12a3", false},
		{"abc", false},
		{"", false},
		{"12.3", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isNumeric(tt.input)
			if result != tt.expected {
				t.Errorf("isNumeric(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
