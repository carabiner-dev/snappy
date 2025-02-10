// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package github

import (
	"errors"
	"os"
)

type TokenReader interface {
	ReadToken() (string, error)
}

type EnvTokenReader struct {
	VarName string
}

func (etr *EnvTokenReader) ReadToken() (string, error) {
	if etr.VarName == "" {
		return "", errors.New("environment variable name to read token not set")
	}
	return os.Getenv(etr.VarName), nil
}
