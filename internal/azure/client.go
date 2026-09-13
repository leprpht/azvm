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

func (c *azureVMClient) NewListAllPager(options *armcompute.VirtualMachinesClientListAllOptions) VMPager {
	return &azureVMPager{
		pager: c.client.NewListAllPager(options),
	}
}

type azureVMPager struct {
	pager *runtime.Pager[armcompute.VirtualMachinesClientListAllResponse]
}

func (p *azureVMPager) More() bool {
	return p.pager.More()
}

func (p *azureVMPager) NextPage(ctx context.Context) (armcompute.VirtualMachinesClientListAllResponse, error) {
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
