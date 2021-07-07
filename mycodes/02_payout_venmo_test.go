package main

import (
	"context"
	"fmt"
	"testing"

	paypal "github.com/plutov/paypal/v4"
	"github.com/shawnhankim/paypal/mycodes/util"
)

func TestPayoutVenmo(t *testing.T) {
	c, res, err := util.GetPaymentGWAccessToken(
		TestClientID, TestSecret, paypal.APIBaseSandBox,
	)
	if err != nil {
		panic(err)
	}
	util.DispPaymentGWAccessToken(res)
	resPayout, err := util.PayoutVenmo(c)
	if err != nil {
		fmt.Println(resPayout)
		panic(err)
	}
	fmt.Println(resPayout)
}

func TestClient_CreatePayout_Venmo(t *testing.T) {
	// Initialize client
	c, err := paypal.NewClient(TestClientID, TestSecret, paypal.APIBaseSandBox)
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
				// Receiver:        "receiver@example.com",
				Receiver: "h.kim3@f5.com",
				Amount: &paypal.AmountPayout{
					Value:    "9.87",
					Currency: "USD",
				},
				Note:         "Thanks for your patronage!",
				SenderItemID: "201403140001",
			},
		},
	}
	// c.CreatePayout(context.Background(), payout)

	res, err := c.CreatePayout(context.Background(), payout)
	fmt.Println(res)
}
