package diagnostics

import (
	"github.com/leprpht/azvm/internal/azure"
	"reflect"
	"strings"
	"testing"
)

func report(status string, network ...azure.NetworkInfo) azure.VMReport {
	return azure.VMReport{VM: azure.VM{Name: "api-01"}, Status: azure.VMStatus{Status: status}, Network: network}
}

func TestRunScenarios(t *testing.T) {
	tests := []struct {
		name, status, want string
		network            []azure.NetworkInfo
		problems           int
	}{
		{"private only", "VM running", "The VM is running and has a private network address,", []azure.NetworkInfo{{NIC: "nic", PrivateIP: "10.0.2.14", PublicIP: "none"}}, 0},
		{"private and public", "VM running", "No immediate configuration problem", []azure.NetworkInfo{{NIC: "nic", PrivateIP: "10.0.2.14", PublicIP: "20.0.0.1"}}, 0},
		{"stopped", "VM stopped", "The VM is not running.", nil, 3},
		{"missing private", "VM running", "no private network address", []azure.NetworkInfo{{NIC: "nic"}}, 1},
		{"missing nic", "VM running", "no network interface was found", nil, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Run(report(tt.status, tt.network...))
			if len(d.Problems) != tt.problems {
				t.Fatalf("problems=%d, want %d", len(d.Problems), tt.problems)
			}
			if !strings.Contains(d.Summary, tt.want) {
				t.Fatalf("summary %q does not contain %q", d.Summary, tt.want)
			}
		})
	}
}

func TestRunIsDeterministic(t *testing.T) {
	r := report("VM running", azure.NetworkInfo{NIC: "nic", PrivateIP: "10.0.0.1"})
	if !reflect.DeepEqual(Run(r), Run(r)) {
		t.Fatal("diagnosis changed")
	}
}
