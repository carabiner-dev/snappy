// SPDX-FileCopyrightText: Copyright 2026 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package snap

import "testing"

// TestAsStatementSingleSubject guards against the snapshot subject being
// added twice, which the GitHub attestations store rejects as a duplicate.
func TestAsStatementSingleSubject(t *testing.T) {
	t.Parallel()
	snap := &Snapshot{
		ID:   "github.com/example/repo@main",
		Name: "github.com/example/repo@main",
		Url:  "git+https://github.com/example/repo/@main",
	}

	subjects := snap.AsStatement().GetSubjects()
	if len(subjects) != 1 {
		t.Fatalf("statement has %d subjects, want exactly 1", len(subjects))
	}
	if subjects[0].GetName() != snap.Name {
		t.Errorf("subject name = %q, want %q", subjects[0].GetName(), snap.Name)
	}
	if len(subjects[0].GetDigest()) != 1 || subjects[0].GetDigest()["sha256"] == "" {
		t.Errorf("subject digest = %v, want a single sha256", subjects[0].GetDigest())
	}
}
