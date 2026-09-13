package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	const subscriptionID = "00000000-0000-0000-0000-000000000000"

	client, err := NewClient(subscriptionID)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}

	if client.VM == nil {
		t.Fatal("NewClient() returned client with nil VM")
	}

	vmClient, ok := client.VM.(*azureVMClient)
	if !ok {
		t.Fatalf("client.VM has type %T, want *azureVMClient", client.VM)
	}

	if vmClient.client == nil {
		t.Fatal("azureVMClient.client is nil")
	}
}

func TestAzureVMClient_NewListAllPager(t *testing.T) {
	t.Parallel()

	const subscriptionID = "00000000-0000-0000-0000-000000000000"

	client, err := NewClient(subscriptionID)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	vmClient, ok := client.VM.(*azureVMClient)
	if !ok {
		t.Fatalf("client.VM has type %T, want *azureVMClient", client.VM)
	}

	pager := vmClient.NewListAllPager(nil)
	if pager == nil {
		t.Fatal("NewListAllPager() returned nil")
	}

	azurePager, ok := pager.(*azureVMPager)
	if !ok {
		t.Fatalf("pager has type %T, want *azureVMPager", pager)
	}

	if azurePager.pager == nil {
		t.Fatal("azureVMPager.pager is nil")
	}
}

func TestAzureVMPager_More(t *testing.T) {
	t.Parallel()

	const subscriptionID = "00000000-0000-0000-0000-000000000000"

	client, err := NewClient(subscriptionID)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	vmClient := client.VM.(*azureVMClient)

	pager := vmClient.NewListAllPager(nil)
	azurePager := pager.(*azureVMPager)

	if !azurePager.More() {
		t.Fatal("More() = false, want true")
	}
}

func TestAzureVMPager_NextPage(t *testing.T) {
	t.Parallel()

	const subscriptionID = "00000000-0000-0000-0000-000000000000"

	client, err := NewClient(subscriptionID)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	vmClient := client.VM.(*azureVMClient)

	pager := vmClient.NewListAllPager(nil)
	azurePager := pager.(*azureVMPager)

	_, err = azurePager.NextPage(t.Context())
	if err == nil {
		t.Fatal("NextPage() error = nil, want error")
	}
}

func TestAzureVMClient_NewListAllPager_WithOptions(t *testing.T) {
	t.Parallel()

	const subscriptionID = "00000000-0000-0000-0000-000000000000"

	client, err := NewClient(subscriptionID)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	vmClient, ok := client.VM.(*azureVMClient)
	if !ok {
		t.Fatalf("client.VM has type %T, want *azureVMClient", client.VM)
	}

	options := &armcompute.VirtualMachinesClientListAllOptions{}

	pager := vmClient.NewListAllPager(options)
	if pager == nil {
		t.Fatal("NewListAllPager() returned nil")
	}

	azurePager, ok := pager.(*azureVMPager)
	if !ok {
		t.Fatalf("pager has type %T, want *azureVMPager", pager)
	}

	if azurePager.pager == nil {
		t.Fatal("azureVMPager.pager is nil")
	}
}
