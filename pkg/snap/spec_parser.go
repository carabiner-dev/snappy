// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package snap

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type SpecParser struct{}

var variableRegex *regexp.Regexp

type ParseOptions struct {
	Variables map[string]string
}

func (sp *SpecParser) ParseWithOptions(r io.Reader, opts *ParseOptions) (*Spec, error) {
	// Find the variables
	fscanner := bufio.NewScanner(r)
	yamlCode := ""
	varErrs := []error{}
	for fscanner.Scan() {
		line := fscanner.Text()
		vars := extractLineVariables(line)
		for _, sub := range vars {
			if value, ok := opts.Variables[sub.Name]; ok {
				line = strings.Replace(line, sub.Replace, value, 1)
			} else {
				varErrs = append(varErrs, fmt.Errorf("no variable substitution defined for %q", sub.Name))
			}
		}
		yamlCode += line + "\n"
		if err := errors.Join(varErrs...); err != nil {
			return nil, err
		}
	}

	logrus.Debugf("Parsed Spec:\n%s", yamlCode)

	// Parse the transformed YAML
	spec := &Spec{}
	if err := yaml.Unmarshal([]byte(yamlCode), spec); err != nil {
		return nil, fmt.Errorf("parsing yaml spec: %w", err)
	}
	return spec, nil
}

type varSub struct {
	Name    string
	Replace string
}

func extractLineVariables(line string) []varSub {
	if variableRegex == nil {
		variableRegex = regexp.MustCompile(`\{\s*\$([-A-Z0-1a-z_]+)\s*\}`)
	}

	res := variableRegex.FindAllStringSubmatch(line, -1)
	ret := []varSub{}
	for _, m := range res {
		ret = append(ret, varSub{
			Name:    m[1],
			Replace: m[0],
		})
	}
	return ret
}
