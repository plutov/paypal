package paypal_test

import (
	"context"

	"github.com/plutov/paypal/v4"
)

func Example() {
	// Initialize client
	c, err := paypal.NewClient("clientID", "secretID", paypal.APIBaseSandBox)
	if err != nil {
		panic(err)
	}

	// Retrieve access token
	_, err = c.GetAccessToken(context.Background())
	if err != nil {
		panic(err)
	}
}

func ExampleClient_CreatePayout_Venmo() {
	// Initialize client
	c, err := paypal.NewClient("clientID", "secretID", paypal.APIBaseSandBox)
	if err != nil {
		panic(err)
	}

	// Retrieve access token
	_, err = c.GetAccessToken(context.Background())
	if err != nil {
		panic(err)
	}

	// Set payout item with Venmo wallet
	payout := paypal.Payout{
		SenderBatchHeader: &paypal.SenderBatchHeader{
			SenderBatchID: "Payouts_2018_100007",
			EmailSubject:  "You have a payout!",
			EmailMessage:  "You have received a payout! Thanks for using our service!",
		},
		Items: []paypal.PayoutItem{
			{
				RecipientType:   "EMAIL",
				RecipientWallet: paypal.VenmoRecipientWallet,
				Receiver:        "receiver@example.com",
				Amount: &paypal.AmountPayout{
					Value:    "9.87",
					Currency: "USD",
				},
				Note:         "Thanks for your patronage!",
				SenderItemID: "201403140001",
			},
		},
	}

	c.CreatePayout(context.Background(), payout)
}
