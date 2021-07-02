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
) (*paypal.TokenResponse, error) {
	c, err := paypal.NewClient(clientID, secret, APIBase)
	if err != nil {
		return nil, err
	}
	return c.GetAccessToken(context.Background())
}

// DispPaymentGWAccessToken shows access token of payment gateway
//
// - example:
//   + Refersh Token :
//   + Token         :  A21AAJ-hLFVkIB2bELTGteZAIUsivRx7tYr6-mnMugWPGfZMzdA7WFN9tkMPkNm9FuSvWEOQZ5EvhJwPHFf4iDAS2zAKIoz-A
//   + Type          :  Bearer
//   + ExpiresIn     :  32012
//   + PayPal Sandbox  :  https://api.sandbox.paypal.com
//     - https://api-m.sandbox.paypal.com unless paypal.APIBaseSandBox is used.
//     - https://developer.paypal.com/docs/api/overview/#make-rest-api-calls
//
func DispPaymentGWAccessToken(tok *paypal.TokenResponse) {
	fmt.Println("- Payment Gateway Access Token & PayPal Sandbox:")
	fmt.Println("  + Refersh Token : ", tok.RefreshToken)
	fmt.Println("  + Token         : ", tok.Token)
	fmt.Println("  + Type          : ", tok.Type)
	fmt.Println("  + ExpiresIn     : ", tok.ExpiresIn)
	fmt.Println("  + PayPal Sandbox: ", paypal.APIBaseSandBox)
}
