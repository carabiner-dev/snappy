// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

// Package specs embeds the builtin snapshot spec definitions so they can be
// used both by the snappy CLI (the "builtin:" spec paths) and by programs
// importing snappy as a library.
package specs

import "embed"

// FS holds the builtin specs, keyed by platform: "github/repo.yaml",
// "gitlab/project.yaml", and so on.
//
//go:embed github/*.yaml gitlab/*.yaml
var FS embed.FS
