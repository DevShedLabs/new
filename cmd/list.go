package cmd

import (
	"fmt"

	"github.com/DevShedLabs/new/internal/blueprint"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates and blueprints",
	RunE: func(cmd *cobra.Command, args []string) error {
		builtins, err := blueprint.ListEmbedded(blueprint.Embedded())
		if err != nil {
			return err
		}

		userBlueprints, err := blueprint.ListUserBlueprints()
		if err != nil {
			return err
		}

		fmt.Println("Built-in templates:")
		if len(builtins) == 0 {
			fmt.Println("  (none)")
		}
		for _, name := range builtins {
			bp, _ := blueprint.Resolve(name, blueprint.Embedded())
			desc := ""
			if bp != nil && bp.Manifest.Description != "" {
				desc = "  — " + bp.Manifest.Description
			}
			fmt.Printf("  %-20s%s\n", name, desc)
		}

		fmt.Println()
		fmt.Println("User blueprints (~/.new/blueprints/):")
		if len(userBlueprints) == 0 {
			fmt.Println("  (none)")
		}
		for _, name := range userBlueprints {
			bp, _ := blueprint.Resolve(name, nil)
			desc := ""
			if bp != nil && bp.Manifest.Description != "" {
				desc = "  — " + bp.Manifest.Description
			}
			fmt.Printf("  %-20s%s\n", name, desc)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
