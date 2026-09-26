package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
)

type Client struct {
	VM      VMClient
	Network NetworkClient
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
func (c *azureVMClient) Get(ctx context.Context, rg, name string, options *armcompute.VirtualMachinesClientGetOptions) (armcompute.VirtualMachinesClientGetResponse, error) {
	return c.client.Get(ctx, rg, name, options)
}

type azureNetworkClient struct {
	interfaces *armnetwork.InterfacesClient
	publicIPs  *armnetwork.PublicIPAddressesClient
}

func (c *azureNetworkClient) GetInterface(ctx context.Context, rg, name string) (string, []NetworkIPConfiguration, error) {
	r, err := c.interfaces.Get(ctx, rg, name, nil)
	if err != nil {
		return "", nil, err
	}
	var out []NetworkIPConfiguration
	if r.Interface.Properties != nil {
		for _, p := range r.Interface.Properties.IPConfigurations {
			if p == nil || p.Properties == nil {
				continue
			}
			out = append(out, NetworkIPConfiguration{PrivateIP: getString(p.Properties.PrivateIPAddress), SubnetID: getSubnetID(p.Properties.Subnet), PublicIPID: getPublicIPID(p.Properties.PublicIPAddress)})
		}
	}
	return getString(r.Interface.Name), out, nil
}
func (c *azureNetworkClient) GetPublicIP(ctx context.Context, rg, name string) (string, error) {
	r, err := c.publicIPs.Get(ctx, rg, name, nil)
	if err != nil {
		return "", err
	}
	if r.PublicIPAddress.Properties == nil {
		return "", nil
	}
	return getString(r.PublicIPAddress.Properties.IPAddress), nil
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
		Network: func() NetworkClient {
			f, _ := armnetwork.NewClientFactory(subscriptionID, cred, nil)
			return &azureNetworkClient{f.NewInterfacesClient(), f.NewPublicIPAddressesClient()}
		}(),
	}, nil
}
