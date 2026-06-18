package file

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/DevShedLabs/New/internal/generator"
)

// embeddedTemplates is set by the main package via SetTemplates.
var embeddedTemplates fs.FS

// SetTemplates wires in the embedded templates/files sub-filesystem.
func SetTemplates(fsys fs.FS) {
	embeddedTemplates = fsys
}

// Generator handles single-file creation.
type Generator struct{}

func (g *Generator) Name() string { return "file" }

// Detect returns true when no --template flag is set, meaning the user wants
// a single file rather than a full project.
func (g *Generator) Detect(ctx *generator.Context) bool {
	return ctx.Template == ""
}

func (g *Generator) Generate(ctx *generator.Context) error {
	dest := filepath.Join(ctx.OutputDir, ctx.Name)

	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("file already exists: %s", dest)
	}

	content, err := boilerplate(ctx)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}

	if err := os.WriteFile(dest, []byte(content), 0o644); err != nil {
		return err
	}

	fmt.Printf("created %s\n", dest)
	return nil
}

// boilerplate returns the rendered template content for the given file, or an
// empty string if no template exists for the extension.
func boilerplate(ctx *generator.Context) (string, error) {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(ctx.Name)), ".")
	if ext == "" || embeddedTemplates == nil {
		return "", nil
	}

	tmplPath := ext + ".tmpl"
	data, err := fs.ReadFile(embeddedTemplates, tmplPath)
	if err != nil {
		// No template for this extension — create an empty file.
		return "", nil
	}

	tmpl, err := template.New(tmplPath).Parse(string(data))
	if err != nil {
		return "", fmt.Errorf("parsing template %s: %w", tmplPath, err)
	}

	// Strip extension from name for display inside templates (e.g. "index" not "index.html").
	baseName := strings.TrimSuffix(ctx.Name, filepath.Ext(ctx.Name))

	vars := map[string]string{"Name": baseName}
	for k, v := range ctx.Vars {
		vars[k] = v
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("rendering template %s: %w", tmplPath, err)
	}

	return buf.String(), nil
}
