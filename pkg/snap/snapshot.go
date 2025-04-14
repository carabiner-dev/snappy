// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package snap

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/carabiner-dev/ampel/pkg/attestation"
	"github.com/carabiner-dev/ampel/pkg/formats/statement/intoto"
	"github.com/carabiner-dev/hasher"
	gointoto "github.com/in-toto/attestation/go/v1"
	"github.com/sirupsen/logrus"
)

type Snapshot struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Url      string            `json:"url"`
	Type     string            `json:"type"`
	Metadata Metadata          `json:"metadata"`
	Headers  map[string]string `json:"headers"`
	Values   any               `json:"values"`
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

func (s *Snapshot) SetSource(attestation.Subject) {

}

func (s *Snapshot) GetSource() attestation.Subject {
	return nil
}

func (s *Snapshot) GetVerifications() []*attestation.SignatureVerification {
	return nil
}

// AsStatement converts the snapshot to an intoto attestation
func (s *Snapshot) AsStatement() attestation.Statement {
	// Create the attestation with the snapshot as predicate
	statement := intoto.NewStatement(
		intoto.WithPredicate(s),
	)

	// Create a hasher to hash the ID
	reader := strings.NewReader(s.ID)
	hshr := hasher.New()

	// The GitHub attestations store only supports a single
	// digest. This means that for now we'll just use sha256.
	hshr.Options.Algorithms = []gointoto.HashAlgorithm{
		gointoto.AlgorithmSHA256,
	}

	var sbj *gointoto.ResourceDescriptor

	hashes, err := hshr.HashReaders([]io.Reader{reader})
	if err != nil {
		logrus.Errorf("error hashing snapshot id: %v", err)
		sbj = &gointoto.ResourceDescriptor{
			Name: s.ID,
			Uri:  s.ID,
		}
	} else {
		sbj = hashes.ToResourceDescriptors()[0]
		sbj.Name = s.Name
		sbj.Uri = s.Url
	}

	statement.AddSubject(sbj)
	return statement
}
