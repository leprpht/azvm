package commands

import (
	"fmt"

	"github.com/leprpht/azvm/internal/azure"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start <vm>",
	Short: "Start an Azure Virtual Machine",
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

		if err := client.StartVM(ctx, vm); err != nil {
			return fmt.Errorf("failed to start VM %q: %w", args[0], err)
		}

		fmt.Printf("Started VM: %s\n", vm.Name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
