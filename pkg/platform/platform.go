package platform

import (
	"context"
	"io"
	"net/http"
)

// Type represents a supported platform
type Type string

const (
	GitHub Type = "github"
	GitLab Type = "gitlab"
)

// Client defines the interface for platform-specific API clients
type Client interface {
	// Call makes an HTTP request to the platform's API
	Call(ctx context.Context, method, path string, body io.Reader) (*http.Response, error)

	// Platform returns the platform type this client supports
	Platform() Type
}

// Factory creates platform-specific clients and provides platform metadata
type Factory interface {
	// CreateClient creates a new platform client
	CreateClient() (Client, error)

	// Platform returns the platform type this factory creates clients for
	Platform() Type

	// DefaultResponseHeaders returns the default headers to include in snapshots
	DefaultResponseHeaders() []string
}
