// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/carabiner-dev/snappy/pkg/snap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type snapOptions struct {
	SpecPath         string
	Attest           bool
	VarSubstitutions []string
}

// Validates the options in context with arguments
func (to *snapOptions) Validate() error {
	errs := []error{}
	if to.SpecPath == "" {
		errs = append(errs, errors.New("spec path not defined"))
	}

	for _, val := range to.VarSubstitutions {
		if !strings.Contains(val, "=") {
			errs = append(errs, fmt.Errorf("variable substitution not well formed: %q", val))
		}
	}
	return errors.Join(errs...)
}

// AddFlags adds the subcommands flags
func (to *snapOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(
		&to.SpecPath, "spec", "s", "", "list of hashes to add as subjects ",
	)
	cmd.PersistentFlags().BoolVarP(
		&to.Attest, "attest", "a", false, "write the output as an (unsigned) in-toto attestation",
	)
	cmd.PersistentFlags().StringSliceVarP(
		&to.VarSubstitutions, "var", "v", []string{}, "spec variable subsitutions (--var name=value)",
	)
}

func addSnap(parentCmd *cobra.Command) {
	opts := &snapOptions{}
	attCmd := &cobra.Command{
		Short:             "takes a snapshot of an API response",
		Use:               "snap spec.yaml",
		Example:           fmt.Sprintf(`%s snap --var REPO=example spec.yaml`, appname),
		SilenceUsage:      false,
		SilenceErrors:     true,
		PersistentPreRunE: initLogging,
		PreRunE: func(_ *cobra.Command, args []string) error {
			if len(args) > 0 && opts.SpecPath == "" {
				opts.SpecPath = args[0]
			}

			if len(args) > 0 && args[0] != opts.SpecPath {
				return fmt.Errorf("spec specified twice (-p and argument)")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate the options
			if err := opts.Validate(); err != nil {
				return err
			}
			cmd.SilenceUsage = true

			f, err := os.Open(opts.SpecPath)
			if err != nil {
				return fmt.Errorf("opening spec: %w", err)
			}
			defer f.Close()

			parseOpts := &snap.ParseOptions{
				Variables: map[string]string{},
			}

			// Slice the variable substitutions
			for _, vs := range opts.VarSubstitutions {
				name, val, _ := strings.Cut(vs, "=")
				parseOpts.Variables[name] = val
			}

			parser := snap.SpecParser{}
			spec, err := parser.ParseWithOptions(f, parseOpts)
			if err != nil {
				return fmt.Errorf("parsing predicate file: %w", err)
			}

			// Create a snapper
			snapper := snap.New()
			// ... and snapshot the repo
			snapshot, err := snapper.Take(context.Background(), spec)
			if err != nil {
				logrus.Errorf("taking snapshot: %v", err)
			}

			if opts.Attest {
				return encodeJSON(snapshot.AsStatement(), os.Stdout)
			}

			return encodeJSON(snapshot, os.Stdout)
		},
	}
	opts.AddFlags(attCmd)
	parentCmd.AddCommand(attCmd)
}

func encodeJSON(what any, output io.Writer) error {
	// Marshal to JSON
	enc := json.NewEncoder(output)
	enc.SetIndent("", "  ")
	if err := enc.Encode(what); err != nil {
		return fmt.Errorf("encoding to JSON: %w", err)
	}
	return nil
}
