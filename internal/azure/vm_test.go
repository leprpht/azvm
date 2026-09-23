package azure

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
)

type mockVMClient struct {
	pager VMPager
}

func (m *mockVMClient) NewListAllPager(_ *armcompute.VirtualMachinesClientListAllOptions) VMPager {
	return m.pager
}

func (m *mockVMClient) GetInstanceView(
	ctx context.Context,
	resourceGroup string,
	vmName string,
	options *armcompute.VirtualMachinesClientInstanceViewOptions,
) (armcompute.VirtualMachinesClientInstanceViewResponse, error) {
	return armcompute.VirtualMachinesClientInstanceViewResponse{}, nil
}

type mockVMPager struct {
	pages []armcompute.VirtualMachinesClientListAllResponse
	err   error
	index int
}

func (m *mockVMPager) More() bool {
	if m.err != nil {
		return m.index == 0
	}

	return m.index < len(m.pages)
}

func (m *mockVMPager) NextPage(_ context.Context) (armcompute.VirtualMachinesClientListAllResponse, error) {
	if m.err != nil {
		return armcompute.VirtualMachinesClientListAllResponse{}, m.err
	}

	page := m.pages[m.index]
	m.index++

	return page, nil
}

func TestListVMs(t *testing.T) {
	tests := []struct {
		name    string
		pager   VMPager
		want    []VM
		wantErr string
	}{
		{
			name: "returns VMs",
			pager: &mockVMPager{
				pages: []armcompute.VirtualMachinesClientListAllResponse{
					{
						VirtualMachineListResult: armcompute.VirtualMachineListResult{
							Value: []*armcompute.VirtualMachine{
								{
									Name:     new("vm-1"),
									ID:       new("/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"),
									Location: new("westeurope"),
								},
								{
									Name:     new("vm-2"),
									ID:       new("/subscriptions/123/resourceGroups/rg-2/providers/Microsoft.Compute/virtualMachines/vm-2"),
									Location: new("northeurope"),
								},
							},
						},
					},
				},
			},
			want: []VM{
				{
					Name:          "vm-1",
					ResourceGroup: "rg-1",
					Location:      "westeurope",
				},
				{
					Name:          "vm-2",
					ResourceGroup: "rg-2",
					Location:      "northeurope",
				},
			},
		},
		{
			name: "handles multiple pages",
			pager: &mockVMPager{
				pages: []armcompute.VirtualMachinesClientListAllResponse{
					{
						VirtualMachineListResult: armcompute.VirtualMachineListResult{
							Value: []*armcompute.VirtualMachine{
								{
									Name:     new("vm-1"),
									ID:       new("/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"),
									Location: new("westeurope"),
								},
							},
						},
					},
					{
						VirtualMachineListResult: armcompute.VirtualMachineListResult{
							Value: []*armcompute.VirtualMachine{
								{
									Name:     new("vm-2"),
									ID:       new("/subscriptions/123/resourceGroups/rg-2/providers/Microsoft.Compute/virtualMachines/vm-2"),
									Location: new("northeurope"),
								},
								{
									Name:     new("vm-3"),
									ID:       new("/subscriptions/123/resourceGroups/rg-3/providers/Microsoft.Compute/virtualMachines/vm-3"),
									Location: new("uksouth"),
								},
							},
						},
					},
				},
			},
			want: []VM{
				{
					Name:          "vm-1",
					ResourceGroup: "rg-1",
					Location:      "westeurope",
				},
				{
					Name:          "vm-2",
					ResourceGroup: "rg-2",
					Location:      "northeurope",
				},
				{
					Name:          "vm-3",
					ResourceGroup: "rg-3",
					Location:      "uksouth",
				},
			},
		},
		{
			name: "skips VM with nil name",
			pager: &mockVMPager{
				pages: []armcompute.VirtualMachinesClientListAllResponse{
					{
						VirtualMachineListResult: armcompute.VirtualMachineListResult{
							Value: []*armcompute.VirtualMachine{
								{
									Name:     nil,
									ID:       new("/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"),
									Location: new("westeurope"),
								},
								{
									Name:     new("vm-2"),
									ID:       new("/subscriptions/123/resourceGroups/rg-2/providers/Microsoft.Compute/virtualMachines/vm-2"),
									Location: new("northeurope"),
								},
							},
						},
					},
				},
			},
			want: []VM{
				{
					Name:          "vm-2",
					ResourceGroup: "rg-2",
					Location:      "northeurope",
				},
			},
		},
		{
			name: "handles nil ID and location",
			pager: &mockVMPager{
				pages: []armcompute.VirtualMachinesClientListAllResponse{
					{
						VirtualMachineListResult: armcompute.VirtualMachineListResult{
							Value: []*armcompute.VirtualMachine{
								{
									Name:     new("vm-1"),
									ID:       nil,
									Location: nil,
								},
							},
						},
					},
				},
			},
			want: []VM{
				{
					Name:          "vm-1",
					ResourceGroup: "",
					Location:      "",
				},
			},
		},
		{
			name: "returns error when Azure request fails",
			pager: &mockVMPager{
				err: errors.New("Azure API unavailable"),
			},
			wantErr: "failed to list VMs: Azure API unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				VM: &mockVMClient{
					pager: tt.pager,
				},
			}

			got, err := client.ListVMs(context.Background())

			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if err.Error() != tt.wantErr {
					t.Fatalf(
						"error = %q, want %q",
						err.Error(),
						tt.wantErr,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tt.want) {
				t.Fatalf(
					"got %d VMs, want %d",
					len(got),
					len(tt.want),
				)
			}

			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf(
						"VM[%d] = %+v, want %+v",
						i,
						got[i],
						tt.want[i],
					)
				}
			}
		})
	}
}

func TestGetString(t *testing.T) {
	tests := []struct {
		name  string
		value *string
		want  string
	}{
		{
			name:  "nil",
			value: nil,
			want:  "",
		},
		{
			name:  "empty",
			value: new(""),
			want:  "",
		},
		{
			name:  "value",
			value: new("westeurope"),
			want:  "westeurope",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getString(tt.value)

			if got != tt.want {
				t.Fatalf(
					"getString() = %q, want %q",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestGetResourceGroup(t *testing.T) {
	tests := []struct {
		name string
		id   *string
		want string
	}{
		{
			name: "nil",
			id:   nil,
			want: "",
		},
		{
			name: "empty",
			id:   new(""),
			want: "",
		},
		{
			name: "valid resource ID",
			id: new(
				"/subscriptions/123/resourceGroups/my-rg/providers/Microsoft.Compute/virtualMachines/my-vm",
			),
			want: "my-rg",
		},
		{
			name: "case insensitive",
			id: new(
				"/subscriptions/123/RESOURCEGROUPS/my-rg/providers/Microsoft.Compute/virtualMachines/my-vm",
			),
			want: "my-rg",
		},
		{
			name: "missing resource group",
			id: new(
				"/subscriptions/123/providers/Microsoft.Compute/virtualMachines/my-vm",
			),
			want: "",
		},
		{
			name: "resourceGroups is last segment",
			id: new(
				"/subscriptions/123/resourceGroups",
			),
			want: "",
		},
		{
			name: "empty resource group",
			id: new(
				"/subscriptions/123/resourceGroups//providers/Microsoft.Compute/virtualMachines/my-vm",
			),
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getResourceGroup(tt.id)

			if got != tt.want {
				t.Fatalf(
					"getResourceGroup() = %q, want %q",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestFindVM(t *testing.T) {
	t.Run("finds VM by name", func(t *testing.T) {
		// mock ListVMs response
	})

	t.Run("returns error when VM does not exist", func(t *testing.T) {
		// mock ListVMs response
	})
}
