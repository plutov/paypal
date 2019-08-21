// +build integration

package paypal

import (
	"testing"
)

// All test values are defined here
var testClientID = "AQzSx89isj-yV7BhuN_TY1s4phiQXlcUEwFPUYD7tWFxts-bf2Zf6f_S0K7J_suOkiZuIKSkNnB1rem-"
var testSecret = "EAW_tyBnkTLxC7RB8CHT39QYZYfT7LwyxPsWle0834O60KGo0A351iMLOFdQBQ5q95DbZM1hOlT9w8Yg"
var testUserID = "https://www.paypal.com/webapps/auth/identity/user/VBqgHcgZwb1PBs69ybjjXfIW86_Hr93aBvF_Rgbh2II"
var testCardID = "CARD-54E6956910402550WKGRL6EA"

func TestGetAccessToken(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	token, err := c.GetAccessToken()
	if err != nil {
		t.Errorf("Not expected error for GetAccessToken(), got %s", err.Error())
	}
	if token.Token == "" {
		t.Errorf("Expected non-empty token for GetAccessToken()")
	}
}

func TestGetUserInfo(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	u, err := c.GetUserInfo("openid")
	if u.ID != testUserID || err != nil {
		t.Errorf("GetUserInfo must return valid test ID %s, got %s, error: %v", testUserID, u.ID, err)
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

	payoutRes, err := c.CreateSinglePayout(payout)

	if err != nil {
		t.Errorf("test single payout is not created, error: %v, payout: %v", err, payoutRes)
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

	r2, e2 := c.GetCreditCards(&CreditCardsFilter{
		Page:     2,
		PageSize: 7,
	})
	if e2 != nil || r2 == nil {
		t.Errorf("200 code expected. Error: %v", e2)
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
