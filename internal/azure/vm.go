package azure

import (
	"context"
	"fmt"
	"strings"
)

type VM struct {
	Name          string
	ResourceGroup string
	Location      string
}

type VMStatus struct {
	Name   string
	Status string
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
			})
		}
	}

	return vms, nil
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
