package paypalsdk_test

import paypalsdk "github.com/logpacker/PayPal-Go-SDK"

func Example() {
	// Initialize client
	c, err := paypalsdk.NewClient("clientID", "secretID", paypalsdk.APIBaseSandBox)
	if err != nil {
		panic(err)
	}

	// Retreive access token
	_, err = c.GetAccessToken()
	if err != nil {
		panic(err)
	}

	// Create credit card payment
	p := paypalsdk.Payment{
		Intent: "sale",
		Payer: &paypalsdk.Payer{
			PaymentMethod: "credit_card",
			FundingInstruments: []paypalsdk.FundingInstrument{paypalsdk.FundingInstrument{
				CreditCard: &paypalsdk.CreditCard{
					Number:      "4111111111111111",
					Type:        "visa",
					ExpireMonth: "11",
					ExpireYear:  "2020",
					CVV2:        "777",
					FirstName:   "John",
					LastName:    "Doe",
				},
			}},
		},
		Transactions: []paypalsdk.Transaction{paypalsdk.Transaction{
			Amount: &paypalsdk.Amount{
				Currency: "USD",
				Total:    "7.00",
			},
			Description: "My Payment",
		}},
		RedirectURLs: &paypalsdk.RedirectURLs{
			ReturnURL: "http://...",
			CancelURL: "http://...",
		},
	}
	_, err = c.CreatePayment(p)
	if err != nil {
		panic(err)
	}
}
