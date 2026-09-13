package azure

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
)

type Client struct {
	VM *armcompute.VirtualMachinesClient
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

	return &Client{VM: factory.NewVirtualMachinesClient()}, nil
}
