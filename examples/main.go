package main

import (
	paypalsdk "PayPal-Go-SDK"
	"strconv"

	"fmt"
)

func main() {
	client, err := paypalsdk.NewClient("AZgwu4yt5Ba0gyTu1dGBH3txHCJbMuFNvrmQxBaQbfDncDiCs6W_rwJD8Ir-0pZrN-_eq7n9zVd8Y-5f", "EBzA1wRl5t73OMugOieDj_tI3vihfJmGl47ukQT-cpctooIzDu0K7IPESNC0cKodlLSOXzwI8qXSM0rd", paypalsdk.APIBaseSandBox)
	if err == nil {
		fmt.Println("DEBUG: ClientID=" + client.ClientID + " APIBase=" + client.APIBase)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}

	token, err := client.GetAccessToken()
	if err == nil {
		fmt.Println("DEBUG: AccessToken=" + token.Token)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}

	payment, err := client.GetPayment("PAY-TEST-123")
	fmt.Println("DEBUG: PaymentID=" + payment.ID)

	payments, err := client.GetPayments()
	if err == nil {
		fmt.Println("DEBUG: PaymentsCount=" + strconv.Itoa(len(payments)))
	} else {
		fmt.Println("ERROR: " + err.Error())
	}

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
				Total:    "200",
			},
			Description: "My Payment",
		}},
		RedirectURLs: &paypalsdk.RedirectURLs{
			ReturnURL: "http://...",
			CancelURL: "http://...",
		},
	}
	paymentResponse, err := client.CreatePayment(p)
	if err == nil {
		fmt.Println("DEBUG: CreatedPaymentID=" + paymentResponse.Payment.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")

	sale, err := client.GetSale("1")
	if err == nil {
		fmt.Println("DEBUG: SaleID=" + sale.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")

	refund, err := client.RefundSale("1", &paypalsdk.Amount{Total: "7.00", Currency: "USD"})
	if err == nil {
		fmt.Println("DEBUG: RefundID=" + refund.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")

	refund, err = client.GetRefund("1")
	if err == nil {
		fmt.Println("DEBUG: RefundID=" + refund.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")

	order, err := client.GetOrder("1")
	if err == nil {
		fmt.Println("DEBUG: OrderID=" + order.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")

	auth, err := client.AuthorizeOrder("1", &paypalsdk.Amount{Total: "7.00", Currency: "USD"})
	if err == nil {
		fmt.Println("DEBUG: AuthID=" + auth.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")

	capture, err := client.CaptureOrder("1", &paypalsdk.Amount{Total: "7.00", Currency: "USD"}, true, nil)
	if err == nil {
		fmt.Println("DEBUG: CaptureID=" + capture.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")

	order, err = client.VoidOrder("1")
	if err == nil {
		fmt.Println("DEBUG: OrderID=" + order.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")
}
