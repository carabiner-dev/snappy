// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package snap

import (
	"fmt"
	"io"
	"regexp"

	yaml "github.com/carabiner-dev/yamplate"
)

type SpecParser struct{}

var variableRegex *regexp.Regexp

type ParseOptions struct {
	Variables map[string]string
}

func (sp *SpecParser) ParseWithOptions(r io.Reader, opts *ParseOptions) (*Spec, error) {
	// Build the parser and
	decoder := yaml.NewDecoder(r)

	// pass the substitution table
	decoder.Options.Variables = opts.Variables

	spec := &Spec{}
	if err := decoder.Decode(spec); err != nil {
		return nil, fmt.Errorf("parsing yaml spec: %w", err)
	}

	return spec, nil
}
