package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
)

type VMPager interface {
	More() bool
	NextPage(context.Context) (armcompute.VirtualMachinesClientListAllResponse, error)
}

type VMClient interface {
	NewListAllPager(*armcompute.VirtualMachinesClientListAllOptions) VMPager
}
