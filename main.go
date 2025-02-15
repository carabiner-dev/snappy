package main

import (
	"embed"

	"github.com/carabiner-dev/snappy/internal/cmd"
)

//go:embed specs/*/*.yaml
var SpecsFS embed.FS

func main() {
	cmd.Execute(&SpecsFS)
}
