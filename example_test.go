package paypalsdk_test

import paypalsdk "github.com/plutov/PayPal-Go-SDK"

func Example() {
	// Initialize client
	c, err := paypalsdk.NewClient("clientID", "secretID", paypalsdk.APIBaseSandBox)
	if err != nil {
		panic(err)
	}

	// Retrieve access token
	_, err = c.GetAccessToken()
	if err != nil {
		panic(err)
	}
}
