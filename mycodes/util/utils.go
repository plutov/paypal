package util

import (
	"context"
	"fmt"

	paypal "github.com/plutov/paypal/v4"
)

// GetPaymentGWAccessToken initializes HTTP client for payment gateway, and
// returns access token that contains token, type, expiresIn.
func GetPaymentGWAccessToken(
	clientID string,
	secret string,
	APIBase string,
) (*paypal.Client, *paypal.TokenResponse, error) {
	c, err := paypal.NewClient(clientID, secret, APIBase)
	if err != nil {
		return c, nil, err
	}
	res, err := c.GetAccessToken(context.Background())
	return c, res, err
}

// DispPaymentGWAccessToken shows access token of payment gateway
//
// - example:
//   + Refersh Token :
//   + Token         : A21AAKOsE0jofF4TUddS3hS9GUOwcC_uOO8Uc3X--VFJCnK7iU6xPUR5T-zuXf9f6u5GNLXIAZYyiYO_NnDsA63xmYf5sPQlg
//   + Type          : Bearer
//   + ExpiresIn     : 32012
//   + PayPal Sandbox: https://api.sandbox.paypal.com
//     - https://api-m.sandbox.paypal.com unless paypal.APIBaseSandBox is used.
//     - https://developer.paypal.com/docs/api/overview/#make-rest-api-calls
//
func DispPaymentGWAccessToken(tok *paypal.TokenResponse) {
	fmt.Println("- Payment Gateway Access Token & PayPal Sandbox:")
	fmt.Println("  + Refersh Token :", tok.RefreshToken)
	fmt.Println("  + Token         :", tok.Token)
	fmt.Println("  + Type          :", tok.Type)
	fmt.Println("  + ExpiresIn     :", tok.ExpiresIn)
	fmt.Println("  + PayPal Sandbox:", paypal.APIBaseSandBox)
}

// PayoutVenmo sets payout item with Venmo wallet
func PayoutVenmo(c *paypal.Client) (*paypal.PayoutResponse, error) {
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
	return c.CreatePayout(context.Background(), payout)
}
