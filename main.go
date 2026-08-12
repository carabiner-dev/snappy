package main

import (
	"github.com/carabiner-dev/snappy/internal/cmd"
	"github.com/carabiner-dev/snappy/pkg/github"
	"github.com/carabiner-dev/snappy/pkg/gitlab"
	"github.com/carabiner-dev/snappy/pkg/platform"
	"github.com/carabiner-dev/snappy/specs"
)

func main() {
	platform.Register(github.NewFactory())
	platform.Register(gitlab.NewFactory())

	cmd.Execute(&specs.FS)
}
