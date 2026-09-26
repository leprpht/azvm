package commands

import (
	"fmt"
	"strings"

	"github.com/leprpht/azvm/internal/azure"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <vm>",
	Short: "Show VM compute and network information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		stopSpinner := startSpinner(cmd.Context())
		defer stopSpinner()
		sub, err := azure.GetSubscriptionID()
		if err != nil {
			return err
		}

		c, err := azure.NewClient(sub)
		if err != nil {
			return fmt.Errorf("failed to create Azure client: %w", err)
		}

		r, err := c.InspectVM(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to inspect VM %q: %w", args[0], err)
		}

		fmt.Print(formatReport(r))
		return nil
	},
}

func formatReport(r azure.VMReport) string {
	var b strings.Builder

	fmt.Fprintf(
		&b,
		"VM: %s\n\n"+
			"Compute\n"+
			"  Status:          %s\n"+
			"  Location:        %s\n"+
			"  Resource Group:  %s\n"+
			"  Size:            %s\n"+
			"  OS:              %s\n\n"+
			"Network\n",
		r.VM.Name,
		r.Status.Status,
		r.VM.Location,
		r.VM.ResourceGroup,
		r.VM.Size,
		r.VM.OS,
	)

	for i, n := range r.Network {
		if i > 0 {
			b.WriteString("\n")
		}

		fmt.Fprintf(
			&b,
			"  NIC:              %s\n"+
				"  VNet:             %s\n"+
				"  Subnet:           %s\n"+
				"  Private IP:       %s\n"+
				"  Public IP:        %s\n",
			n.NIC,
			n.VNet,
			n.Subnet,
			n.PrivateIP,
			n.PublicIP,
		)
	}

	return b.String()
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
