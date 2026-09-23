package azure

import "testing"

func TestGetSubscriptionID(t *testing.T) {
	subscriptionID, err := GetSubscriptionID()
	if err != nil {
		t.Fatalf("GetSubscriptionID() returned error: %v", err)
	}

	if subscriptionID == "" {
		t.Fatal("GetSubscriptionID() returned empty subscription ID")
	}

	t.Logf("subscription ID: %s", subscriptionID)
}
