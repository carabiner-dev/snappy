// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package snap

import "time"

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
