package main

import (
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
}
