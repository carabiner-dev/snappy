// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package snap

import (
	"fmt"
	"io"

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

	return spec, nil
}
