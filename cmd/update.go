package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update new to the latest version",
	Long:  `Installs the latest released version of new directly from GitHub, bypassing any local cache.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("updating new...")

		gobin, err := exec.LookPath("go")
		if err != nil {
			return fmt.Errorf("go not found on PATH — cannot self-update")
		}

		c := exec.Command(gobin, "install", "github.com/DevShedLabs/new@latest")
		c.Env = append(os.Environ(), "GOFLAGS=-mod=mod")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			return fmt.Errorf("update failed: %w", err)
		}

		fmt.Println("done — run `new --version` to confirm")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
