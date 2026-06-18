package scaffold

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/DevShedLabs/New/internal/blueprint"
	"github.com/DevShedLabs/New/internal/generator"
)

// embeddedProjects is set by the main package via SetTemplates.
var embeddedProjects fs.FS

// SetTemplates wires in the embedded templates/projects sub-filesystem.
func SetTemplates(fsys fs.FS) {
	embeddedProjects = fsys
}

// Generator handles project scaffolding from built-in templates and user blueprints.
type Generator struct{}

func (g *Generator) Name() string { return "scaffold" }

// Detect returns true when a --template flag is provided.
func (g *Generator) Detect(ctx *generator.Context) bool {
	return ctx.Template != ""
}

func (g *Generator) Generate(ctx *generator.Context) error {
	bp, err := blueprint.Resolve(ctx.Template, embeddedProjects)
	if err != nil {
		return fmt.Errorf("template %q not found — check built-in templates or ~/.new/blueprints/", ctx.Template)
	}

	destRoot := filepath.Join(ctx.OutputDir, ctx.Name)
	if _, err := os.Stat(destRoot); err == nil {
		return fmt.Errorf("destination already exists: %s", destRoot)
	}

	vars := buildVars(ctx, bp)

	if err := walkAndWrite(bp.FS, destRoot, vars); err != nil {
		return err
	}

	fmt.Printf("created project %s from template %q\n", destRoot, ctx.Template)
	return nil
}

// walkAndWrite copies every file in the blueprint FS into destRoot,
// rendering each file path and content through text/template.
func walkAndWrite(fsys fs.FS, destRoot string, vars map[string]string) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the manifest — it's metadata, not a project file.
		if path == "blueprint.yaml" {
			return nil
		}

		// Render the destination path so blueprint authors can use {{.Name}} in filenames.
		renderedPath, err := renderString(path, vars)
		if err != nil {
			return fmt.Errorf("rendering path %q: %w", path, err)
		}

		dest := filepath.Join(destRoot, filepath.FromSlash(renderedPath))

		if d.IsDir() {
			return os.MkdirAll(dest, 0o755)
		}

		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}

		rendered, err := renderString(string(data), vars)
		if err != nil {
			return fmt.Errorf("rendering file %q: %w", path, err)
		}

		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}

		return os.WriteFile(dest, []byte(rendered), 0o644)
	})
}

func renderString(s string, vars map[string]string) (string, error) {
	tmpl, err := template.New("").Option("missingkey=zero").Parse(s)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func buildVars(ctx *generator.Context, bp *blueprint.Blueprint) map[string]string {
	vars := map[string]string{
		"Name":     ctx.Name,
		"Template": ctx.Template,
	}
	// Apply blueprint defaults first, then caller-supplied vars override.
	for k, v := range bp.Manifest.Defaults {
		vars[k] = v
	}
	for k, v := range ctx.Vars {
		vars[k] = v
	}
	return vars
}
