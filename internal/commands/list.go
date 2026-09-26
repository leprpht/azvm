package commands

import (
	"fmt"

	"github.com/leprpht/azvm/internal/azure"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Azure Virtual Machines",
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

		vms, err := client.ListVMs(ctx)
		if err != nil {
			return fmt.Errorf("failed to list VMs: %w", err)
		}

		if len(vms) == 0 {
			fmt.Println("No virtual machines found.")
			return nil
		}

		for _, vm := range vms {
			fmt.Printf(
				"Name: %s, Resource Group: %s, Location: %s\n",
				vm.Name,
				vm.ResourceGroup,
				vm.Location,
			)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
