package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
)

type Client struct {
	VM VMClient
}

type azureVMClient struct {
	client *armcompute.VirtualMachinesClient
}

func (c *azureVMClient) NewListAllPager(
	options *armcompute.VirtualMachinesClientListAllOptions,
) VMPager {
	return &azureVMPager{
		pager: c.client.NewListAllPager(options),
	}
}

func (c *azureVMClient) GetInstanceView(
	ctx context.Context,
	resourceGroup string,
	vmName string,
	options *armcompute.VirtualMachinesClientInstanceViewOptions,
) (armcompute.VirtualMachinesClientInstanceViewResponse, error) {
	return c.client.InstanceView(ctx, resourceGroup, vmName, options)
}

func (c *azureVMClient) BeginStart(
	ctx context.Context,
	resourceGroup string,
	vmName string,
	options *armcompute.VirtualMachinesClientBeginStartOptions,
) (*runtime.Poller[armcompute.VirtualMachinesClientStartResponse], error) {
	return c.client.BeginStart(ctx, resourceGroup, vmName, options)
}

func (c *azureVMClient) BeginPowerOff(
	ctx context.Context,
	resourceGroup string,
	vmName string,
	options *armcompute.VirtualMachinesClientBeginPowerOffOptions,
) (*runtime.Poller[armcompute.VirtualMachinesClientPowerOffResponse], error) {
	return c.client.BeginPowerOff(ctx, resourceGroup, vmName, options)
}

func (c *azureVMClient) BeginRestart(
	ctx context.Context,
	resourceGroup string,
	vmName string,
	options *armcompute.VirtualMachinesClientBeginRestartOptions,
) (*runtime.Poller[armcompute.VirtualMachinesClientRestartResponse], error) {
	return c.client.BeginRestart(ctx, resourceGroup, vmName, options)
}

type azureVMPager struct {
	pager *runtime.Pager[armcompute.VirtualMachinesClientListAllResponse]
}

func (p *azureVMPager) More() bool {
	return p.pager.More()
}

func (p *azureVMPager) NextPage(
	ctx context.Context,
) (armcompute.VirtualMachinesClientListAllResponse, error) {
	return p.pager.NextPage(ctx)
}

func NewClient(subscriptionID string) (*Client, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	factory, err := armcompute.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		VM: &azureVMClient{
			client: factory.NewVirtualMachinesClient(),
		},
	}, nil
}
