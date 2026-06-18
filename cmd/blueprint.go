package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/DevShedLabs/new/internal/blueprint"
	"github.com/spf13/cobra"
)

var blueprintCmd = &cobra.Command{
	Use:   "blueprint",
	Short: "Manage user blueprints",
}

var captureCmd = &cobra.Command{
	Use:   "capture [name] [path]",
	Short: "Capture a project directory as a user blueprint",
	Long: `Snapshot an existing project into ~/.new/blueprints/<name>/ so it can be
used as a template with: new <project> --template <name>

If path is omitted, the current directory is used.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		srcPath := "."
		if len(args) == 2 {
			srcPath = args[1]
		}

		src, err := filepath.Abs(srcPath)
		if err != nil {
			return err
		}

		info, err := os.Stat(src)
		if err != nil || !info.IsDir() {
			return fmt.Errorf("source path is not a directory: %s", src)
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		dest := filepath.Join(home, ".new", "blueprints", name)

		if _, err := os.Stat(dest); err == nil {
			return fmt.Errorf("blueprint %q already exists at %s — remove it first to recapture", name, dest)
		}

		if err := copyDir(src, dest); err != nil {
			return err
		}

		if err := writeManifest(dest, name); err != nil {
			return err
		}

		fmt.Printf("captured %s → %s\n", src, dest)
		fmt.Printf("use it with: new <project> --template %s\n", name)
		return nil
	},
}

func init() {
	blueprintCmd.AddCommand(captureCmd)
	rootCmd.AddCommand(blueprintCmd)
}

// copyDir recursively copies src into dest, skipping common build artifacts
// and hidden dirs that shouldn't be part of a blueprint.
func copyDir(src, dest string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if shouldSkip(rel, d) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		target := filepath.Join(dest, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		return copyFile(path, target)
	})
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// writeManifest creates a blueprint.yaml in dest if one doesn't already exist.
func writeManifest(dest, name string) error {
	manifestPath := filepath.Join(dest, "blueprint.yaml")
	if _, err := os.Stat(manifestPath); err == nil {
		return nil
	}

	m := blueprint.Manifest{
		Name:        name,
		Description: fmt.Sprintf("Captured from project %q.", name),
		Vars:        []string{"Name"},
		Defaults:    map[string]string{},
	}

	data, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	return os.WriteFile(manifestPath, data, 0o644)
}

// shouldSkip returns true for paths that should not be included in a blueprint.
var skipDirs = map[string]bool{
	"node_modules": true,
	".git":         true,
	"vendor":       true,
	"dist":         true,
	"build":        true,
	".next":        true,
	".nuxt":        true,
	"__pycache__":  true,
	".venv":        true,
}

func shouldSkip(rel string, d fs.DirEntry) bool {
	if rel == "." {
		return false
	}
	top := strings.SplitN(rel, string(os.PathSeparator), 2)[0]
	if d.IsDir() && skipDirs[top] {
		return true
	}
	// Skip blueprint.yaml at the root — writeManifest will generate a fresh one.
	if rel == "blueprint.yaml" {
		return true
	}
	return false
}
