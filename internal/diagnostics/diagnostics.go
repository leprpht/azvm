package diagnostics

import (
	"strings"

	"github.com/leprpht/azvm/internal/azure"
)

type CheckStatus string

const (
	CheckPass CheckStatus = "pass"
	CheckWarn CheckStatus = "warn"
	CheckFail CheckStatus = "fail"
	CheckSkip CheckStatus = "skip"
)

type Check struct {
	Name   string
	Status CheckStatus
	Detail string
}
type Problem struct{ Severity, Title, Detail, Cause string }
type Diagnosis struct {
	Checks   []Check
	Problems []Problem
	Summary  string
}

func Run(r azure.VMReport) Diagnosis {
	d := Diagnosis{}
	running := strings.Contains(strings.ToLower(r.Status.Status), "running")
	if running {
		d.Checks = append(d.Checks, Check{"VM running", CheckPass, "running"})
	} else {
		d.Checks = append(d.Checks, Check{"VM running", CheckFail, statusDetail(r.Status.Status)})
		d.Problems = append(d.Problems, Problem{"high", "VM is not running", "Network diagnostics may be unavailable or unreliable until the VM is started.", "The VM power state is not running."})
	}

	nic, privateIP, publicIP := false, "", ""
	for _, n := range r.Network {
		if n.NIC != "" {
			nic = true
		}
		if privateIP == "" && n.PrivateIP != "" && n.PrivateIP != "none" {
			privateIP = n.PrivateIP
		}
		if publicIP == "" && n.PublicIP != "" && n.PublicIP != "none" {
			publicIP = n.PublicIP
		}
	}
	if nic {
		d.Checks = append(d.Checks, Check{"Network interface", CheckPass, "attached"})
	} else {
		d.Checks = append(d.Checks, Check{"Network interface", CheckFail, "none attached"})
		d.Problems = append(d.Problems, Problem{"high", "No network interface is attached", "The VM has no network interface available in the collected information.", "No attached network interface was found."})
	}
	if privateIP != "" {
		d.Checks = append(d.Checks, Check{"Private IP", CheckPass, privateIP})
	} else {
		d.Checks = append(d.Checks, Check{"Private IP", CheckWarn, "none"})
		d.Problems = append(d.Problems, Problem{"medium", "No private IP was found", "The VM has no private network address in the collected information.", "No private IP configuration was found."})
	}
	if publicIP != "" {
		d.Checks = append(d.Checks, Check{"Public endpoint", CheckPass, publicIP})
	} else {
		d.Checks = append(d.Checks, Check{"Public endpoint", CheckWarn, "none"})
	}
	d.Summary = summary(running, nic, privateIP != "", publicIP != "")
	return d
}

func statusDetail(s string) string {
	if s == "" {
		return "unknown"
	}
	return strings.ToLower(strings.TrimPrefix(s, "VM "))
}

func summary(running, nic, private, public bool) string {
	if !running {
		return "The VM is not running.\n\nNetwork connectivity cannot be reliably diagnosed until the VM is started."
	}
	if private && !public {
		return "The VM is running and has a private network address,\nbut it does not have a public endpoint.\n\nThe VM may only be reachable from networks that have\naccess to its private address."
	}
	if private {
		return "The VM is running and has a private network address.\n\nNo immediate configuration problem was detected from\nthe available VM information."
	}
	if !nic {
		return "The VM is running, but no network interface was found.\n\nNetwork connectivity cannot be reliably diagnosed from\nthe available VM information."
	}
	return "The VM is running, but no private network address was found.\n\nNetwork connectivity cannot be reliably diagnosed from\nthe available VM information."
}
