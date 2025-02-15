// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package snap

import (
	"errors"
	"fmt"
	"strings"
)

type Spec struct {
	// ID is the string that will be used to generate the subject's
	// hash when generating an attestation. It should identify the
	// object described by the data returned by the API call. This
	// ID MUST be unique for each isntance.
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Url      string   `json:"url"`
	Type     string   `json:"type"`
	Endpoint string   `json:"endpoint"`
	Method   string   `json:"method"`
	Mask     []string `json:"mask"`
}

func (spec *Spec) Validate() error {
	errs := []error{}
	if spec.Method != "" && spec.Method != "POST" && spec.Method != "GET" {
		errs = append(errs, errors.New("wrong method, it shoud be POST or GET"))
	}
	if spec.Endpoint == "" {
		errs = append(errs, errors.New("no endpoint defined"))
	}

	if strings.HasPrefix(spec.Endpoint, "http") {
		errs = append(errs, errors.New("endpoint should be a relative path"))
	}

	if len(spec.Mask) == 0 {
		return fmt.Errorf("at least one entry in the field mask should be set")
	}

	return errors.Join(errs...)
}
