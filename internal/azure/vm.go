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
