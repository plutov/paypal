package main

import (
	"fmt"

	"github.com/logpacker/PayPal-Go-SDK"
)

func main() {
	c, err := paypalsdk.NewClient("AZgwu4yt5Ba0gyTu1dGBH3txHCJbMuFNvrmQxBaQbfDncDiCs6W_rwJD8Ir-0pZrN-_eq7n9zVd8Y-5f", "EBzA1wRl5t73OMugOieDj_tI3vihfJmGl47ukQT-cpctooIzDu0K7IPESNC0cKodlLSOXzwI8qXSM0rd", paypalsdk.APIBaseSandBox)
	if err == nil {
		fmt.Println("DEBUG: ClientID=" + c.ClientID + " APIBase=" + c.APIBase)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}

	token, err := c.GetAccessToken()
	if err == nil {
		fmt.Println("DEBUG: AccessToken=" + token.Token)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}

	payment, err := c.GetPayment("PAY-TEST-123")
	fmt.Println("DEBUG: PaymentID=" + payment.ID)

	payments, err := c.GetPayments()
	if err == nil {
		fmt.Printf("DEBUG: PaymentsCount=%d\n", len(payments))
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
	paymentResponse, err := c.CreatePayment(p)
	if err == nil {
		fmt.Println("DEBUG: CreatedPaymentID=" + paymentResponse.Payment.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")

	sale, err := c.GetSale("1")
	if err == nil {
		fmt.Println("DEBUG: SaleID=" + sale.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")

	refund, err := c.RefundSale("1", &paypalsdk.Amount{Total: "7.00", Currency: "USD"})
	if err == nil {
		fmt.Println("DEBUG: RefundID=" + refund.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")

	refund, err = c.GetRefund("1")
	if err == nil {
		fmt.Println("DEBUG: RefundID=" + refund.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")

	order, err := c.GetOrder("1")
	if err == nil {
		fmt.Println("DEBUG: OrderID=" + order.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")

	auth, err := c.AuthorizeOrder("1", &paypalsdk.Amount{Total: "7.00", Currency: "USD"})
	if err == nil {
		fmt.Println("DEBUG: AuthID=" + auth.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")

	capture, err := c.CaptureOrder("1", &paypalsdk.Amount{Total: "7.00", Currency: "USD"}, true, nil)
	if err == nil {
		fmt.Println("DEBUG: CaptureID=" + capture.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")

	order, err = c.VoidOrder("1")
	if err == nil {
		fmt.Println("DEBUG: OrderID=" + order.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
	fmt.Println("OK")

	token1, err := c.GrantNewAccessTokenFromAuthCode("123", "http://example.com/myapp/return.php")
	if err == nil {
		fmt.Println("DEBUG: Token=" + token1.Token)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}

	token2, err := c.GrantNewAccessTokenFromRefreshToken("123")
	if err == nil {
		fmt.Println("DEBUG: Token=" + token2.Token)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}

	u, err := c.GetUserInfo("openid")
	if err == nil {
		fmt.Println("DEBUG: UserID=" + u.ID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}

	payout := paypalsdk.Payout{
		SenderBatchHeader: &paypalsdk.SenderBatchHeader{
			EmailSubject: "Subject will be displayed on PayPal",
		},
		Items: []paypalsdk.PayoutItem{
			paypalsdk.PayoutItem{
				RecipientType: "EMAIL",
				Receiver:      "single-email-payout@mail.com",
				Amount: &paypalsdk.AmountPayout{
					Value:    "15.11",
					Currency: "USD",
				},
				Note:         "Optional note",
				SenderItemID: "Optional Item ID",
			},
		},
	}

	payoutResp, err := c.CreateSinglePayout(payout)
	if err == nil && len(payoutResp.Items) > 0 {
		fmt.Println("DEBUG: PayoutItemID=" + payoutResp.Items[0].PayoutItemID)
	} else {
		fmt.Println("ERROR: " + err.Error())
	}
}
