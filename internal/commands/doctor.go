package commands

import (
	"fmt"
	"strings"

	"github.com/leprpht/azvm/internal/azure"
	"github.com/leprpht/azvm/internal/diagnostics"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor <vm>",
	Short: "Diagnose available VM and network information",
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

		fmt.Print(formatDiagnosis(r, diagnostics.Run(r)))
		return nil
	},
}

func formatDiagnosis(r azure.VMReport, d diagnostics.Diagnosis) string {
	var b strings.Builder

	fmt.Fprintf(
		&b,
		"VM: %s\n\n"+
			"Compute\n"+
			"  Status:          %s %s\n\n"+
			"Network\n",
		r.VM.Name,
		mark(find(d, "VM running")),
		strings.ToLower(strings.TrimPrefix(r.Status.Status, "VM ")),
	)

	nic, privateIP, publicIP := "none", "none", "none"

	for _, n := range r.Network {
		if n.NIC != "" {
			nic = "attached"
		}
		if privateIP == "none" && n.PrivateIP != "" {
			privateIP = n.PrivateIP
		}
		if publicIP == "none" && n.PublicIP != "" {
			publicIP = n.PublicIP
		}
	}

	fmt.Fprintf(
		&b,
		"  NIC attached:    %s %s\n"+
			"  Private IP:      %s %s\n"+
			"  Public IP:       %s %s\n\n"+
			"Checks\n",
		mark(find(d, "Network interface")),
		nic,
		mark(find(d, "Private IP")),
		privateIP,
		mark(find(d, "Public IP")),
		publicIP,
	)

	for _, c := range d.Checks {
		fmt.Fprintf(
			&b,
			"  %-20s %s %s\n",
			c.Name,
			statusMark(c.Status),
			c.Detail,
		)
	}

	b.WriteString("\nDiagnosis\n")

	for _, line := range strings.Split(d.Summary, "\n") {
		fmt.Fprintf(&b, "  %s\n", line)
	}

	return b.String()
}

func find(d diagnostics.Diagnosis, name string) diagnostics.Check {
	for _, c := range d.Checks {
		if c.Name == name {
			return c
		}
	}

	return diagnostics.Check{Status: diagnostics.CheckSkip}
}

func statusMark(s diagnostics.CheckStatus) string {
	switch s {
	case diagnostics.CheckPass:
		return "✓"
	case diagnostics.CheckFail:
		return "✗"
	default:
		return "!"
	}
}

func mark(c diagnostics.Check) string {
	return statusMark(c.Status)
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
