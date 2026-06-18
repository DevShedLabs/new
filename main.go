package main

import (
	"embed"
	"io/fs"

	"github.com/DevShedLabs/new/cmd"
	"github.com/DevShedLabs/new/internal/blueprint"
	"github.com/DevShedLabs/new/internal/file"
	"github.com/DevShedLabs/new/internal/generator"
	"github.com/DevShedLabs/new/internal/scaffold"
)

//go:embed templates
var templatesFS embed.FS

func main() {
	// Wire embedded templates into each generator package.
	filesFS, _ := fs.Sub(templatesFS, "templates/files")
	projectsFS, _ := fs.Sub(templatesFS, "templates/projects")

	file.SetTemplates(filesFS)
	scaffold.SetTemplates(projectsFS)
	blueprint.SetEmbedded(projectsFS)

	// Register generators in priority order.
	// Scaffold is checked first so --template always wins.
	generator.Register(&scaffold.Generator{})
	generator.Register(&file.Generator{})

	cmd.Execute()
}
