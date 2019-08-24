package paypal

func Example() {
	// Initialize client
	c, err := NewClient("clientID", "secretID", APIBaseSandBox)
	if err != nil {
		panic(err)
	}

	// Retrieve access token
	_, err = c.GetAccessToken()
	if err != nil {
		panic(err)
	}
}
