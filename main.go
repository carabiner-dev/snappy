package main

import (
	"embed"

	"github.com/carabiner-dev/snappy/internal/cmd"
	// Register platform implementations
	_ "github.com/carabiner-dev/snappy/pkg/github"
	_ "github.com/carabiner-dev/snappy/pkg/gitlab"
)

//go:embed specs/*/*.yaml
var SpecsFS embed.FS

func main() {
	cmd.Execute(&SpecsFS)
}
