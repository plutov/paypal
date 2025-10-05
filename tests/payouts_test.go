package tests

import (
	"context"
	"testing"

	"github.com/plutov/paypal/v4"
)

func TestCreatePayout(t *testing.T) {
	ctx := context.Background()
	server := createTestServer()
	defer server.Close()

	client, err := paypal.NewClient(clientId, clientSecret, server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	client.Token = &paypal.TokenResponse{Token: "dummy"}

	payout := paypal.Payout{
		Items: []paypal.PayoutItem{
			{
				Amount: &paypal.AmountPayout{
					Currency: "USD",
					Value:    "10.00",
				},
				Receiver: "test@example.com",
			},
		},
		SenderBatchHeader: &paypal.SenderBatchHeader{
			SenderBatchID: "BATCH123",
		},
	}

	response, err := client.CreatePayout(ctx, payout)
	assertNoError(t, err)

	assertEqual(t, "BATCH123", response.BatchHeader.PayoutBatchID)
	assertEqual(t, "SUCCESS", response.BatchHeader.BatchStatus)
}

func TestGetPayout(t *testing.T) {
	ctx := context.Background()
	server := createTestServer()
	defer server.Close()

	client, err := paypal.NewClient(clientId, clientSecret, server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	client.Token = &paypal.TokenResponse{Token: "dummy"}

	response, err := client.GetPayout(ctx, "BATCH123")
	assertNoError(t, err)

	assertEqual(t, "BATCH123", response.BatchHeader.PayoutBatchID)
	assertEqual(t, "SUCCESS", response.BatchHeader.BatchStatus)
}

func TestGetPayoutItem(t *testing.T) {
	ctx := context.Background()
	server := createTestServer()
	defer server.Close()

	client, err := paypal.NewClient(clientId, clientSecret, server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	client.Token = &paypal.TokenResponse{Token: "dummy"}

	response, err := client.GetPayoutItem(ctx, "ITEM123")
	assertNoError(t, err)

	assertEqual(t, "ITEM123", response.PayoutItemID)
	assertEqual(t, "BATCH123", response.PayoutBatchID)
	assertEqual(t, "SUCCESS", response.TransactionStatus)
}

func TestCancelPayoutItem(t *testing.T) {
	ctx := context.Background()
	server := createTestServer()
	defer server.Close()

	client, err := paypal.NewClient(clientId, clientSecret, server.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	client.Token = &paypal.TokenResponse{Token: "dummy"}

	response, err := client.CancelPayoutItem(ctx, "ITEM123")
	assertNoError(t, err)

	assertEqual(t, "ITEM123", response.PayoutItemID)
	assertEqual(t, "BATCH123", response.PayoutBatchID)
	assertEqual(t, "SUCCESS", response.TransactionStatus)
}
