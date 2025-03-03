// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package snap

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

const (
	PayloadTypeStruct = "struct"
	PayloadTypeArray  = "array"
)

var PayloadTypes = []string{PayloadTypeStruct, PayloadTypeArray}

type Spec struct {
	// ID is the string that will be used to generate the subject's
	// hash when generating an attestation. It should identify the
	// object described by the data returned by the API call. This
	// ID MUST be unique for each instance the API returns.
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Url         string   `json:"url"`
	Type        string   `json:"type"`
	Endpoint    string   `json:"endpoint"`
	Method      string   `json:"method"`
	PayloadType string   `json:"payload" yaml:"payload"`
	Mask        []string `json:"mask"`
	Data        string   `json:"data"`
	TrimNL      bool     `json:"trimNL" yaml:"trimNL"`
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

	if spec.PayloadType != "" && !slices.Contains(PayloadTypes, spec.PayloadType) {
		errs = append(errs, fmt.Errorf("unsupported payload type must be: %+v", PayloadTypes))
	}

	if spec.PayloadType == PayloadTypeArray && len(spec.Mask) > 0 {
		errs = append(errs, errors.New("field mask not supported when payload is an array"))
	}

	if spec.PayloadType == PayloadTypeStruct && len(spec.Mask) == 0 {
		errs = append(errs, fmt.Errorf("at least one entry in the field mask should be set"))
	}

	if spec.Data != "" && spec.Method != "POST" {
		errs = append(errs, errors.New("data can only be specified with method POST"))
	}

	return errors.Join(errs...)
}
