package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Azure",
	RunE: func(cmd *cobra.Command, args []string) error {
		command := exec.Command("az", "login")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Stdin = os.Stdin

		if err := command.Run(); err != nil {
			return fmt.Errorf("Azure login failed: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
