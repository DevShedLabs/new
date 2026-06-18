package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/DevShedLabs/new/internal/generator"
	"github.com/spf13/cobra"
)

var (
	flagTemplate string
	flagOutput   string
	flagVars     []string
)

var rootCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Scaffold files and projects",
	Long: `new creates files with smart boilerplate or scaffolds full projects from templates.

Examples:
  new index.html                   Create an HTML file with boilerplate
  new main.go                      Create a Go file with package main
  new my-app --template react      Scaffold a React project
  new my-site --template html      Scaffold an HTML project
  new my-app --template my-bp      Scaffold from a user blueprint in ~/.new/blueprints/`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		outDir := cwd
		if flagOutput != "" {
			outDir = flagOutput
		}

		ctx := &generator.Context{
			Name:      args[0],
			Template:  flagTemplate,
			OutputDir: outDir,
			Vars:      parseVars(flagVars),
		}

		gen, err := generator.Resolve(ctx)
		if err != nil {
			return err
		}

		return gen.Generate(ctx)
	},
}

// Execute is the entrypoint called by main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&flagTemplate, "template", "t", "", "Template or blueprint name to scaffold a project")
	rootCmd.Flags().StringVarP(&flagOutput, "output", "o", "", "Output directory (defaults to current directory)")
	rootCmd.Flags().StringArrayVarP(&flagVars, "var", "v", nil, "Template variables as key=value pairs (repeatable)")
}

// parseVars converts ["key=value", ...] into a map.
func parseVars(raw []string) map[string]string {
	m := make(map[string]string, len(raw))
	for _, kv := range raw {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}
