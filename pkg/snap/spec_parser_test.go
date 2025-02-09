// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package snap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVariableExtract(t *testing.T) {
	for _, tc := range []struct {
		name   string
		line   string
		expect []varSub
	}{
		{"no-subs", " Just a line", []varSub{}},
		{"1-sub", ` Just a line {$HALO}`, []varSub{{"HALO", `{$HALO}`}}},
		{"1-sub-spaces", ` Just a line { $HALO }`, []varSub{{"HALO", `{ $HALO }`}}},
		{"1-sub-mixedcase", ` Just a line { $Halo }`, []varSub{{"Halo", `{ $Halo }`}}},
		{"2-subs", ` Just a line {$HALO} { $BYE }`, []varSub{{"HALO", `{$HALO}`}, {"BYE", `{ $BYE }`}}},
		{"no-dollar", ` Just a line {HALO} `, []varSub{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			vars := extractLineVariables(tc.line)
			require.Len(t, vars, len(tc.expect))
			for i := range vars {
				require.Equal(t, tc.expect[i].Name, vars[i].Name)
				require.Equal(t, tc.expect[i].Replace, vars[i].Replace)
			}
		})
	}
}
