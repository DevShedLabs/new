package blueprint

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// embeddedProjects is set by main via SetEmbedded.
var embeddedProjects fs.FS

// SetEmbedded wires in the embedded templates/projects sub-filesystem so that
// List and Resolve can reference built-in templates.
func SetEmbedded(fsys fs.FS) {
	embeddedProjects = fsys
}

// Embedded returns the registered embedded projects FS.
func Embedded() fs.FS {
	return embeddedProjects
}

// Manifest is the optional blueprint.yaml descriptor inside a blueprint folder.
type Manifest struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	// Vars lists variable names the blueprint expects (used for prompting later).
	Vars        []string          `yaml:"vars"`
	// Defaults provides fallback values for vars.
	Defaults    map[string]string `yaml:"defaults"`
}

// Blueprint represents a resolved blueprint — either a user blueprint from
// ~/.new/blueprints/<name> or an embedded built-in template.
type Blueprint struct {
	Name     string
	Manifest *Manifest
	// FS is the filesystem rooted at the blueprint directory.
	FS fs.FS
}

// userBlueprintsDir returns ~/.new/blueprints.
func userBlueprintsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".new", "blueprints"), nil
}

// Resolve looks for a blueprint by name. User blueprints in ~/.new/blueprints/
// take precedence over the provided embedded FS (pass the embedded templates/projects
// sub-FS as fallback).
func Resolve(name string, embedded fs.FS) (*Blueprint, error) {
	// 1. Check user blueprints first.
	if bp, err := loadFromUserDir(name); err == nil {
		return bp, nil
	}

	// 2. Fall back to embedded templates.
	if embedded != nil {
		sub, err := fs.Sub(embedded, name)
		if err == nil {
			return &Blueprint{
				Name:     name,
				Manifest: loadManifestFromFS(sub),
				FS:       sub,
			}, nil
		}
	}

	return nil, errors.New("blueprint not found: " + name)
}

// ListEmbedded returns names of all built-in templates in the provided embedded FS.
func ListEmbedded(embedded fs.FS) ([]string, error) {
	if embedded == nil {
		return nil, nil
	}
	entries, err := fs.ReadDir(embedded, ".")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

// ListUserBlueprints returns names of all blueprints in ~/.new/blueprints/.
func ListUserBlueprints() ([]string, error) {
	dir, err := userBlueprintsDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func loadFromUserDir(name string) (*Blueprint, error) {
	dir, err := userBlueprintsDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, name)
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return nil, errors.New("not found")
	}
	fsys := os.DirFS(path)
	return &Blueprint{
		Name:     name,
		Manifest: loadManifestFromFS(fsys),
		FS:       fsys,
	}, nil
}

func loadManifestFromFS(fsys fs.FS) *Manifest {
	data, err := fs.ReadFile(fsys, "blueprint.yaml")
	if err != nil {
		return &Manifest{}
	}
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return &Manifest{}
	}
	return &m
}
