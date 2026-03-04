// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package gitlab

import (
	"errors"
	"fmt"
	"os"
)

type Options struct {
	Host        string
	Token       string
	TokenReader TokenReader
}

var defaultOptions = Options{
	Host:        getDefaultHost(),
	TokenReader: &EnvTokenReader{VarName: "GITLAB_TOKEN"},
}

// getDefaultHost returns the GitLab host from GITLAB_HOST env var, or gitlab.com
func getDefaultHost() string {
	if host := os.Getenv("GITLAB_HOST"); host != "" {
		return host
	}
	return "gitlab.com"
}

// ensureToken makes sure we have a token. If there is none set, we
// read it using the token reader
func (o *Options) ensureToken() error {
	if o.Token != "" {
		return nil
	}
	if o.TokenReader == nil {
		return errors.New("no token set and no token reader configured")
	}

	token, err := o.TokenReader.ReadToken()
	if err != nil {
		return fmt.Errorf("reading token: %w", err)
	}

	if token == "" {
		return fmt.Errorf("unable to get a token from the token reader")
	}

	o.Token = token
	return nil
}
