// +build integration

package paypalsdk

import (
	"os"
	"testing"
)

// All test values are defined here
var testClientID = "AZgwu4yt5Ba0gyTu1dGBH3txHCJbMuFNvrmQxBaQbfDncDiCs6W_rwJD8Ir-0pZrN-_eq7n9zVd8Y-5f"
var testSecret = "EBzA1wRl5t73OMugOieDj_tI3vihfJmGl47ukQT-cpctooIzDu0K7IPESNC0cKodlLSOXzwI8qXSM0rd"
var testAuthID = "2DC87612EK520411B"
var testFakeAuthID = "FAKE-2DC87612EK520411B"
var testOrderID = "O-4J082351X3132253H"
var testFakeOrderID = "FAKE-O-4J082351X3132253H"
var testSaleID = "4CF18861HF410323U"
var testPaymentID = "PAY-5YK922393D847794YKER7MUI"
var testPayerID = "CR87QHB7JTRSC"
var testUserID = "https://www.paypal.com/webapps/auth/identity/user/WEssgRpQij92sE99_F9MImvQ8FPYgUEjrvCja2qH2H8"
var testCardID = "CARD-54E6956910402550WKGRL6EA"

func TestGetAccessToken(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	token, err := c.GetAccessToken()
	if err != nil {
		t.Errorf("Not expected error for GetAccessToken()")
	}
	if token.Token == "" {
		t.Errorf("Expected non-empty token for GetAccessToken()")
	}
}

func TestGetAuthorization(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	a, err := c.GetAuthorization(testAuthID)
	if err != nil || a.ID != testAuthID {
		t.Errorf("GetAuthorization failed for ID=%s", testAuthID)
	}

	a, err = c.GetAuthorization(testFakeAuthID)
	if err == nil {
		t.Errorf("GetAuthorization must return error for ID=%s", testFakeAuthID)
	}
}

func TestCaptureAuthorization(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.CaptureAuthorization(testAuthID, &Amount{Total: "200", Currency: "USD"}, true)

	if err == nil {
		t.Errorf("Auth is expired, 400 error must be returned")
	}
}

func TestVoidAuthorization(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.VoidAuthorization(testAuthID)

	if err == nil {
		t.Errorf("Auth is expired, 400 error must be returned")
	}
}

func TestReauthorizeAuthorization(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.ReauthorizeAuthorization(testAuthID, &Amount{Total: "200", Currency: "USD"})

	if err == nil {
		t.Errorf("Reauthorization not allowed for this product, 500 error must be returned")
	}
}

func TestGrantNewAccessTokenFromAuthCode(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)

	_, err := c.GrantNewAccessTokenFromAuthCode("123", "http://example.com/myapp/return.php")
	if err == nil {
		t.Errorf("GrantNewAccessTokenFromAuthCode must return error for invalid code")
	}
}

func TestGrantNewAccessTokenFromRefreshToken(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)

	_, err := c.GrantNewAccessTokenFromRefreshToken("123")
	if err == nil {
		t.Errorf("GrantNewAccessTokenFromRefreshToken must return error for invalid refresh token")
	}
}

func TestGetUserInfo(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	u, err := c.GetUserInfo("openid")
	if u.ID != testUserID || err != nil {
		t.Errorf("GetUserInfo must return valid test ID=%s", testUserID)
	}
}

func TestGetOrder(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	o, err := c.GetOrder(testOrderID)
	if err != nil || o.ID != testOrderID {
		t.Errorf("GetOrder failed for ID=%s", testOrderID)
	}

	o, err = c.GetOrder(testFakeOrderID)
	if err == nil {
		t.Errorf("GetOrder must return error for ID=%s", testFakeOrderID)
	}
}

func TestAuthorizeOrder(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.AuthorizeOrder(testOrderID, &Amount{Total: "7.00", Currency: "USD"})
	if err == nil {
		t.Errorf("Order is expired, 400 error must be returned")
	}
}

func TestCaptureOrder(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.CaptureOrder(testOrderID, &Amount{Total: "100", Currency: "USD"}, true, nil)
	if err == nil {
		t.Errorf("Order is expired, 400 error must be returned")
	}
}

func TestVoidOrder(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.VoidOrder(testOrderID)
	if err == nil {
		t.Errorf("Order is expired, 400 error must be returned")
	}
}

func TestCreateDirectPaypalPayment(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.SetLog(os.Stdout)
	c.GetAccessToken()

	amount := Amount{
		Total:    "15.11",
		Currency: "USD",
	}

	p, err := c.CreateDirectPaypalPayment(amount, "http://example.com", "http://example.com", "test payment")

	if err != nil || p.ID == "" {
		t.Errorf("Test paypal payment is not created")
	}
}

func TestCreatePayment(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.SetLog(os.Stdout)
	c.GetAccessToken()

	p := Payment{
		Intent: "sale",
		Payer: &Payer{
			PaymentMethod: "credit_card",
			FundingInstruments: []FundingInstrument{{
				CreditCard: &CreditCard{
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
		Transactions: []Transaction{{
			Amount: &Amount{
				Currency: "USD",
				Total:    "10.00", // total cost including shipping
				Details: Details{
					Shipping: "3.00", // total shipping cost
					Subtotal: "7.00", // total cost without shipping
				},
			},
			Description: "My Payment",
			ItemList: &ItemList{
				Items: []Item{
					Item{
						Quantity: 2,
						Price:    "3.50",
						Currency: "USD",
						Name:     "Product 1",
					},
				},
			},
		}},
		RedirectURLs: &RedirectURLs{
			ReturnURL: "http://..",
			CancelURL: "http://..",
		},
	}
	pr, err := c.CreatePayment(p)
	if err != nil {
		t.Errorf("Error creating payment.")
	}
}

func TestGetPayment(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.GetPayment(testPaymentID)

	if err == nil {
		t.Errorf("404 for this payment ID")
	}
}

func TestGetPayments(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.GetPayments()

	if err != nil {
		t.Errorf("Nil error expected")
	}
}

func TestExecuteApprovedPayment(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.ExecuteApprovedPayment(testPaymentID, testPayerID)

	if err == nil {
		t.Errorf("404 for this payment ID")
	}
}

func TestCreateSinglePayout(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	payout := Payout{
		SenderBatchHeader: &SenderBatchHeader{
			EmailSubject: "Subject will be displayed on PayPal",
		},
		Items: []PayoutItem{
			{
				RecipientType: "EMAIL",
				Receiver:      "single-email-payout@mail.com",
				Amount: &AmountPayout{
					Value:    "15.11",
					Currency: "USD",
				},
				Note:         "Optional note",
				SenderItemID: "Optional Item ID",
			},
		},
	}

	p, err := c.CreateSinglePayout(payout)

	if err != nil || len(p.Items) != 1 {
		t.Errorf("Test single payout is not created")
	}
}

func TestGetSale(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.GetSale(testSaleID)
	if err == nil {
		t.Errorf("404 must be returned for ID=%s", testSaleID)
	}
}

func TestRefundSale(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.RefundSale(testSaleID, nil)
	if err == nil {
		t.Errorf("404 must be returned for ID=%s", testSaleID)
	}

	_, err = c.RefundSale(testSaleID, &Amount{Total: "7.00", Currency: "USD"})
	if err == nil {
		t.Errorf("404 must be returned for ID=%s", testSaleID)
	}
}

func TestGetRefund(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.GetRefund("1")
	if err == nil {
		t.Errorf("404 must be returned for ID=%s", testSaleID)
	}
}

func TestStoreCreditCard(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	r1, e1 := c.StoreCreditCard(CreditCard{})
	if e1 == nil || r1 != nil {
		t.Errorf("Error is expected for invalid CC")
	}

	r2, e2 := c.StoreCreditCard(CreditCard{
		Number:      "4417119669820331",
		Type:        "visa",
		ExpireMonth: "11",
		ExpireYear:  "2020",
		CVV2:        "874",
		FirstName:   "Foo",
		LastName:    "Bar",
	})
	if e2 != nil || r2 == nil {
		t.Errorf("200 code expected for valid CC card. Error: %v", e2)
	}
}

func TestDeleteCreditCard(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	e1 := c.DeleteCreditCard("")
	if e1 == nil {
		t.Errorf("Error is expected for invalid CC ID")
	}
}

func TestGetCreditCard(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	r1, e1 := c.GetCreditCard("BBGGG")
	if e1 == nil || r1 != nil {
		t.Errorf("Error is expected for invalid CC, got CC %v", r1)
	}
}

func TestGetCreditCards(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	r1, e1 := c.GetCreditCards(nil)
	if e1 != nil || r1 == nil {
		t.Errorf("200 code expected. Error: %v", e1)
	}
	if r1.TotalItems < 1 {
		t.Errorf("Expected >0 CCs, got %d", r1.TotalItems)
	}
	if r1.TotalPages < 1 {
		t.Errorf("Expected >0 CCs page")
	}

	r2, e2 := c.GetCreditCards(&CreditCardsFilter{
		Page:     2,
		PageSize: 7,
	})
	if e2 != nil || r2 == nil {
		t.Errorf("200 code expected. Error: %v", e2)
	}
	if r2.TotalItems < 1 {
		t.Errorf("Expected >0 CCs, got %d", r2.TotalItems)
	}
	if r2.TotalPages < 1 {
		t.Errorf("Expected >0 CCs page")
	}
}

func TestPatchCreditCard(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	r1, e1 := c.PatchCreditCard(testCardID, nil)
	if e1 == nil || r1 != nil {
		t.Errorf("Error is expected for empty update info")
	}
}
