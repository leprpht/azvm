package azure

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
	"strings"
)

type VM struct {
	Name          string
	ResourceGroup string
	Location      string
	Size, OS      string
	NetworkIDs    string
}
type NetworkInfo struct{ NIC, VNet, Subnet, PrivateIP, PublicIP string }
type VMReport struct {
	VM      VM
	Status  VMStatus
	Network []NetworkInfo
}

type VMStatus struct {
	Name   string
	Status string
}

func (c *Client) ListVMs(ctx context.Context) ([]VM, error) {
	var vms []VM

	pager := c.VM.NewListAllPager(nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list VMs: %w", err)
		}

		for _, vm := range page.Value {
			if vm.Name == nil {
				continue
			}

			vms = append(vms, VM{
				Name:          *vm.Name,
				ResourceGroup: getResourceGroup(vm.ID),
				Location:      getString(vm.Location),
				Size:          getVMSize(vm), OS: getVMOS(vm), NetworkIDs: getNetworkIDs(vm),
			})
		}
	}

	return vms, nil
}

func (c *Client) InspectVM(ctx context.Context, name string) (VMReport, error) {
	vm, err := c.FindVM(ctx, name)
	if err != nil {
		return VMReport{}, err
	}
	status, err := c.GetVMStatus(ctx, vm)
	if err != nil {
		return VMReport{}, err
	}
	r := VMReport{VM: vm, Status: status}
	for _, id := range strings.Split(vm.NetworkIDs, ";") {
		if id == "" {
			continue
		}
		nic := lastID(id)
		nn, cs, err := c.Network.GetInterface(ctx, vm.ResourceGroup, nic)
		if err != nil {
			return VMReport{}, fmt.Errorf("failed to get network interface %q: %w", nic, err)
		}
		for _, cfg := range cs {
			n := NetworkInfo{NIC: nn, VNet: vnetFromSubnet(cfg.SubnetID), Subnet: lastID(cfg.SubnetID), PrivateIP: cfg.PrivateIP, PublicIP: "none"}
			if cfg.PublicIPID != "" {
				n.PublicIP, err = c.Network.GetPublicIP(ctx, vm.ResourceGroup, lastID(cfg.PublicIPID))
				if err != nil {
					return VMReport{}, fmt.Errorf("failed to get public IP: %w", err)
				}
				if n.PublicIP == "" {
					n.PublicIP = "none"
				}
			}
			r.Network = append(r.Network, n)
		}
	}
	return r, nil
}
func getNetworkIDs(vm *armcompute.VirtualMachine) string {
	var out []string
	if vm.Properties != nil && vm.Properties.NetworkProfile != nil {
		for _, n := range vm.Properties.NetworkProfile.NetworkInterfaces {
			if n != nil && n.ID != nil {
				out = append(out, *n.ID)
			}
		}
	}
	return strings.Join(out, ";")
}
func getVMSize(vm *armcompute.VirtualMachine) string {
	if vm.Properties != nil && vm.Properties.HardwareProfile != nil {
		if vm.Properties.HardwareProfile.VMSize != nil {
			return string(*vm.Properties.HardwareProfile.VMSize)
		}
	}
	return ""
}
func getVMOS(vm *armcompute.VirtualMachine) string {
	if vm.Properties != nil && vm.Properties.OSProfile != nil {
		if vm.Properties.OSProfile.LinuxConfiguration != nil {
			return "Linux"
		}
		if vm.Properties.OSProfile.WindowsConfiguration != nil {
			return "Windows"
		}
	}
	return ""
}
func getSubnetID(v *armnetwork.Subnet) string {
	if v == nil {
		return ""
	}
	return getString(v.ID)
}
func getPublicIPID(v *armnetwork.PublicIPAddress) string {
	if v == nil {
		return ""
	}
	return getString(v.ID)
}
func lastID(id string) string {
	p := strings.Split(strings.TrimSuffix(id, "/"), "/")
	if len(p) == 0 {
		return ""
	}
	return p[len(p)-1]
}
func vnetFromSubnet(id string) string {
	p := strings.Split(id, "/")
	for i := 0; i < len(p)-1; i++ {
		if strings.EqualFold(p[i], "virtualNetworks") {
			return p[i+1]
		}
	}
	return ""
}

func (c *Client) FindVM(ctx context.Context, name string) (VM, error) {
	vms, err := c.ListVMs(ctx)
	if err != nil {
		return VM{}, err
	}

	for _, vm := range vms {
		if vm.Name == name {
			return vm, nil
		}
	}

	return VM{}, fmt.Errorf("VM %q not found", name)
}

func (c *Client) GetVMStatus(ctx context.Context, vm VM) (VMStatus, error) {
	response, err := c.VM.GetInstanceView(
		ctx,
		vm.ResourceGroup,
		vm.Name,
		nil,
	)
	if err != nil {
		return VMStatus{}, fmt.Errorf("failed to get VM status: %w", err)
	}

	for _, status := range response.Statuses {
		if status.Code == nil || status.DisplayStatus == nil {
			continue
		}

		if strings.HasPrefix(*status.Code, "PowerState/") {
			return VMStatus{
				Name:   vm.Name,
				Status: *status.DisplayStatus,
			}, nil
		}
	}

	return VMStatus{
		Name:   vm.Name,
		Status: "unknown",
	}, nil
}

func (c *Client) StartVM(ctx context.Context, vm VM) error {
	poller, err := c.VM.BeginStart(
		ctx,
		vm.ResourceGroup,
		vm.Name,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start VM: %w", err)
	}

	if _, err := poller.PollUntilDone(ctx, nil); err != nil {
		return fmt.Errorf("failed to start VM: %w", err)
	}

	return nil
}

func (c *Client) StopVM(ctx context.Context, vm VM) error {
	poller, err := c.VM.BeginPowerOff(
		ctx,
		vm.ResourceGroup,
		vm.Name,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to stop VM: %w", err)
	}

	if _, err := poller.PollUntilDone(ctx, nil); err != nil {
		return fmt.Errorf("failed to stop VM: %w", err)
	}

	return nil
}

func (c *Client) RestartVM(ctx context.Context, vm VM) error {
	poller, err := c.VM.BeginRestart(
		ctx,
		vm.ResourceGroup,
		vm.Name,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to restart VM: %w", err)
	}

	if _, err := poller.PollUntilDone(ctx, nil); err != nil {
		return fmt.Errorf("failed to restart VM: %w", err)
	}

	return nil
}

func getString(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func getResourceGroup(id *string) string {
	if id == nil {
		return ""
	}

	parts := strings.Split(*id, "/")

	for i := 0; i < len(parts)-1; i++ {
		if strings.EqualFold(parts[i], "resourceGroups") {
			return parts[i+1]
		}
	}

	return ""
}
