// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package snap

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/carabiner-dev/ampel/pkg/attestation"
	"github.com/sirupsen/logrus"
)

type Snapshot struct {
	ID       string              `json:"id"`
	Type     string              `json:"type"`
	Metadata Metadata            `json:"metadata"`
	Headers  map[string][]string `json:"headers"`
	Values   map[string]any      `json:"values"`
}

type Metadata struct {
	Date     time.Time `json:"date"`
	Endpoint string    `json:"endpoint"`
	Method   string    `json:"method"`
}

func (s *Snapshot) GetData() []byte {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		logrus.Errorf("marshaling data: %v", err)
		return nil
	}
	return b.Bytes()
}

func (s *Snapshot) GetParsed() any {
	return s
}

func (s *Snapshot) GetType() attestation.PredicateType {
	return attestation.PredicateType(s.Type)
}

func (s *Snapshot) SetType(t attestation.PredicateType) error {
	s.Type = string(t)
	return nil
}
