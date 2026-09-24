package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
)

type VMPager interface {
	More() bool
	NextPage(context.Context) (armcompute.VirtualMachinesClientListAllResponse, error)
}

type VMClient interface {
	NewListAllPager(
		*armcompute.VirtualMachinesClientListAllOptions,
	) VMPager

	GetInstanceView(
		context.Context,
		string,
		string,
		*armcompute.VirtualMachinesClientInstanceViewOptions,
	) (armcompute.VirtualMachinesClientInstanceViewResponse, error)

	BeginStart(
		context.Context,
		string,
		string,
		*armcompute.VirtualMachinesClientBeginStartOptions,
	) (*runtime.Poller[armcompute.VirtualMachinesClientStartResponse], error)

	BeginPowerOff(
		context.Context,
		string,
		string,
		*armcompute.VirtualMachinesClientBeginPowerOffOptions,
	) (*runtime.Poller[armcompute.VirtualMachinesClientPowerOffResponse], error)

	BeginRestart(
		context.Context,
		string,
		string,
		*armcompute.VirtualMachinesClientBeginRestartOptions,
	) (*runtime.Poller[armcompute.VirtualMachinesClientRestartResponse], error)
}
