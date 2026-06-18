package main

import (
	"embed"
	"io/fs"

	"github.com/DevShedLabs/New/cmd"
	"github.com/DevShedLabs/New/internal/file"
	"github.com/DevShedLabs/New/internal/generator"
	"github.com/DevShedLabs/New/internal/scaffold"
)

//go:embed templates
var templatesFS embed.FS

func main() {
	// Wire embedded templates into each generator package.
	filesFS, _ := fs.Sub(templatesFS, "templates/files")
	projectsFS, _ := fs.Sub(templatesFS, "templates/projects")

	file.SetTemplates(filesFS)
	scaffold.SetTemplates(projectsFS)

	// Register generators in priority order.
	// Scaffold is checked first so --template always wins.
	generator.Register(&scaffold.Generator{})
	generator.Register(&file.Generator{})

	cmd.Execute()
}
