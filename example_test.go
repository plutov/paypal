package paypal_test

import "github.com/plutov/paypal/v3"

func Example() {
	// Initialize client
	c, err := paypal.NewClient("clientID", "secretID", paypal.APIBaseSandBox)
	if err != nil {
		panic(err)
	}

	// Retrieve access token
	_, err = c.GetAccessToken()
	if err != nil {
		panic(err)
	}
}
