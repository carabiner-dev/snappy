package platform

import (
	"fmt"
	"strings"
)

// Detector detects the platform type from various sources
type Detector struct{}

// NewDetector creates a new platform detector
func NewDetector() *Detector {
	return &Detector{}
}

// DetectFromSpec detects platform from spec file path
// Looks for patterns like "github/" or "gitlab/" in the path
func (d *Detector) DetectFromSpec(specPath string) (Type, error) {
	// Check for builtin: prefix patterns
	if strings.Contains(specPath, "github/") || strings.Contains(specPath, ":github/") {
		return GitHub, nil
	}
	if strings.Contains(specPath, "gitlab/") || strings.Contains(specPath, ":gitlab/") {
		return GitLab, nil
	}

	return "", fmt.Errorf("could not detect platform from spec path: %s", specPath)
}

// DetectFromEndpoint detects platform from API endpoint pattern
// GitHub uses patterns like "repos/", "orgs/", "users/"
// GitLab uses patterns like "projects/", "groups/"
func (d *Detector) DetectFromEndpoint(endpoint string) (Type, error) {
	// Normalize endpoint
	endpoint = strings.TrimPrefix(endpoint, "/")

	// GitHub-specific patterns
	if strings.HasPrefix(endpoint, "repos/") ||
		strings.HasPrefix(endpoint, "orgs/") ||
		strings.HasPrefix(endpoint, "user/") ||
		strings.HasPrefix(endpoint, "users/") {
		return GitHub, nil
	}

	// GitLab-specific patterns
	if strings.HasPrefix(endpoint, "projects/") ||
		strings.HasPrefix(endpoint, "groups/") ||
		strings.HasPrefix(endpoint, "merge_requests/") {
		return GitLab, nil
	}

	return "", fmt.Errorf("could not detect platform from endpoint: %s", endpoint)
}
