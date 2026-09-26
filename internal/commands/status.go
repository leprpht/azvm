package commands

import (
	"fmt"

	"github.com/leprpht/azvm/internal/azure"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <vm>",
	Short: "Show the status of an Azure Virtual Machine",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		stopSpinner := startSpinner(cmd.Context())
		defer stopSpinner()
		ctx := cmd.Context()

		subscriptionID, err := azure.GetSubscriptionID()
		if err != nil {
			return err
		}

		client, err := azure.NewClient(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to create Azure client: %w", err)
		}

		vm, err := client.FindVM(ctx, args[0])
		if err != nil {
			return err
		}

		status, err := client.GetVMStatus(ctx, vm)
		if err != nil {
			return err
		}

		fmt.Printf("Name: %s, Status: %s\n", status.Name, status.Status)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
