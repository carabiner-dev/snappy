// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"embed"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sigs.k8s.io/release-utils/log"
	"sigs.k8s.io/release-utils/version"
)

const appname = "snappy"

var rootCmd = &cobra.Command{
	Short:             fmt.Sprintf("%s: create an attestable predicate of an API response", appname),
	Long:              fmt.Sprintf(`%s: create an attestable predicate of an API response`, appname),
	Use:               appname,
	SilenceUsage:      false,
	PersistentPreRunE: initLogging,
	Example: fmt.Sprintf(`
Create a snapshot of an API response:

	%s snap --ver REPO=example spec.yaml
	`, appname),
}

type commandLineOptions struct {
	logLevel string
}

var commandLineOpts = commandLineOptions{}

func initLogging(*cobra.Command, []string) error {
	return log.SetupGlobalLogger(commandLineOpts.logLevel)
}

// Execute builds the command
func Execute(fs *embed.FS) {
	rootCmd.PersistentFlags().StringVar(
		&commandLineOpts.logLevel,
		"log-level", "info", fmt.Sprintf("the logging verbosity, either %s", log.LevelNames()),
	)
	addSnap(rootCmd, fs)
	rootCmd.AddCommand(version.WithFont("doom"))

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
