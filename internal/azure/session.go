package azure

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetSubscriptionID() (string, error) {
	cmd := exec.Command(
		"az",
		"account",
		"show",
		"--query",
		"id",
		"-o",
		"tsv",
	)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get Azure subscription: %w", err)
	}

	subscriptionID := strings.TrimSpace(string(output))

	if subscriptionID == "" {
		return "", fmt.Errorf("no Azure subscription is selected")
	}

	return subscriptionID, nil
}
