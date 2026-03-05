package main

import (
	"embed"

	"github.com/carabiner-dev/snappy/internal/cmd"
	"github.com/carabiner-dev/snappy/pkg/github"
	"github.com/carabiner-dev/snappy/pkg/gitlab"
	"github.com/carabiner-dev/snappy/pkg/platform"
)

//go:embed specs/*/*.yaml
var SpecsFS embed.FS

func main() {
	platform.Register(github.NewFactory())
	platform.Register(gitlab.NewFactory())

	cmd.Execute(&SpecsFS)
}
