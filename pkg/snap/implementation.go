// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package snap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/carabiner-dev/snappy/pkg/github"
)

type SnapperImplementation interface {
	ValidateSpec(*Spec) error
	GetClient() (*github.Client, error)
	CallAPI(context.Context, *github.Client, *Spec) (*http.Response, error)
	ParseResponse(*Options, *Spec, *http.Response) (*Snapshot, error)
	ApplyFieldMask(*Snapshot, []string) (*Snapshot, error)
}

// defaultImplementation implements the logic for a SnapperImplementation
type defaultImplementation struct{}

func (di *defaultImplementation) ValidateSpec(spec *Spec) error {
	return spec.Validate()
}

func (di *defaultImplementation) GetClient() (*github.Client, error) {
	return github.NewClient()
}

func (di *defaultImplementation) CallAPI(ctx context.Context, client *github.Client, spec *Spec) (*http.Response, error) {
	return client.Call(ctx, spec.Method, spec.Endpoint, nil)
}

// ParseResponse extracts the data from the response and returns the snapshot
func (di *defaultImplementation) ParseResponse(opts *Options, spec *Spec, resp *http.Response) (*Snapshot, error) {
	values := map[string]any{}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http error %d received from API", resp.StatusCode)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing response body: %w", err)
	}

	snapshot := &Snapshot{
		ID:   spec.ID,
		Type: spec.Type,
		Metadata: Metadata{
			Date:     time.Now(),
			Endpoint: spec.Endpoint,
			Method:   spec.Method,
		},
		Headers: map[string]string{},
		Values:  values,
	}

	for k, v := range resp.Header {
		if slices.Contains(opts.ResponseHeaders, k) {
			for _, content := range v {
				if _, ok := snapshot.Headers[k]; !ok {
					snapshot.Headers[k] = content
				} else {
					snapshot.Headers[k] = "; " + content
				}
			}
		}
	}

	// Umarshal
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, fmt.Errorf("unmarshaling response data: %w", err)
	}

	// Done, return the new snapshot
	return snapshot, nil
}

func (di *defaultImplementation) ApplyFieldMask(snapshot *Snapshot, fieldmap []string) (*Snapshot, error) {
	newValues := map[string]any{}
	for k, val := range snapshot.Values {
		if slices.Contains(fieldmap, k) {
			newValues[k] = val
		}
	}

	snapshot.Values = newValues
	return snapshot, nil
}
