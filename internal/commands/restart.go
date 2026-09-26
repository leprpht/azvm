package commands

import (
	"fmt"

	"github.com/leprpht/azvm/internal/azure"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart <vm>",
	Short: "Restart an Azure Virtual Machine",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		stopSpinner := startSpinner(cmd.Context())
		defer stopSpinner()
		ctx := cmd.Context()

		subscriptionID, err := azure.GetSubscriptionID()
		if err != nil {
			return fmt.Errorf("failed to get subscription ID: %w", err)
		}

		client, err := azure.NewClient(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to create Azure client: %w", err)
		}

		vm, err := client.FindVM(ctx, args[0])
		if err != nil {
			return fmt.Errorf("failed to find VM %q: %w", args[0], err)
		}

		if err := client.RestartVM(ctx, vm); err != nil {
			return fmt.Errorf("failed to restart VM %q: %w", args[0], err)
		}

		fmt.Printf("Restarted VM: %s\n", vm.Name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
