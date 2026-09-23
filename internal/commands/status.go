package commands

import (
	"context"
	"fmt"

	"github.com/leprpht/azvm/internal/azure"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <vm>",
	Short: "Show the status of an Azure Virtual Machine",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID, err := azure.GetSubscriptionID()
		if err != nil {
			return err
		}

		client, err := azure.NewClient(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to create Azure client: %w", err)
		}

		vms, err := client.ListVMs(context.Background())
		if err != nil {
			return fmt.Errorf("failed to list VMs: %w", err)
		}

		var target *azure.VM

		for i := range vms {
			if vms[i].Name == args[0] {
				target = &vms[i]
				break
			}
		}

		if target == nil {
			return fmt.Errorf("VM %q not found", args[0])
		}

		status, err := client.GetVMStatus(context.Background(), *target)
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
