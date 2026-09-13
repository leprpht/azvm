package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/leprpht/azvm/internal/azure"
)

func main() {
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")

	if subscriptionID == "" {
		log.Fatal("AZURE_SUBSCRIPTION_ID environment variable is not set")
	}

	client, err := azure.NewClient(subscriptionID)
	if err != nil {
		log.Fatal("AZURE_SUBSCRIPTION_ID environment variable is not set")
	}

	fmt.Println("Azure client created successfully")
	fmt.Println(client)

	vms, err := client.ListVMs(context.Background())
	if err != nil {
		log.Fatalf("Failed to list VMs: %v", err)
	}

	if len(vms) == 0 {
		fmt.Println("No virtual machines found.")
		return
	}

	for _, vm := range vms {
		fmt.Printf("%s (%s)\n", vm.Name, vm.Location)
	}
}
