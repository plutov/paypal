// +build integration

package paypal

import (
	"testing"
)

// All test values are defined here
var testClientID = "AXy9orp-CDaHhBZ9C78QHW2BKZpACgroqo85_NIOa9mIfJ9QnSVKzY-X_rivR_fTUUr6aLjcJsj6sDur"
var testSecret = "EBoIiUSkCKeSk49hHSgTem1qnjzzJgRQHDEHvGpzlLEf_nIoJd91xu8rPOBDCdR_UYNKVxJE-UgS2iCw"
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
			SenderBatchID: "Payouts_2018_100007",
			EmailSubject:  "You have a payout!",
			EmailMessage:  "You have received a payout! Thanks for using our service!",
		},
		Items: []PayoutItem{
			{
				RecipientType: "EMAIL",
				Receiver:      "receiver@example.com",
				Amount: &AmountPayout{
					Value:    "9.87",
					Currency: "USD",
				},
				Note:         "Thanks for your patronage!",
				SenderItemID: "201403140001",
			},
		},
	}

	c.CreateSinglePayout(payout)
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

// Creates, gets, and deletes single webhook
func TestCreateAndGetWebhook(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	payload := &CreateWebhookRequest{
		URL: "https://example.com/paypal_webhooks",
		EventTypes: []WebhookEventType{
			WebhookEventType{
				Name: "PAYMENT.AUTHORIZATION.CREATED",
			},
		},
	}

	createdWebhook, err := c.CreateWebhook(payload)
	if err != nil {
		t.Errorf("Webhook couldn't be created, error %v", err)
	}

	_, err = c.GetWebhook(createdWebhook.ID)
	if err != nil {
		t.Errorf("An error occurred while getting webhook, error %v", err)
	}

	err = c.DeleteWebhook(createdWebhook.ID)
	if err != nil {
		t.Errorf("An error occurred while webhooks deletion, error %v", err)
	}
}

// Creates, updates, and deletes single webhook
func TestCreateAndUpdateWebhook(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	creationPayload := &CreateWebhookRequest{
		URL: "https://example.com/paypal_webhooks",
		EventTypes: []WebhookEventType{
			WebhookEventType{
				Name: "PAYMENT.AUTHORIZATION.CREATED",
			},
		},
	}

	createdWebhook, err := c.CreateWebhook(creationPayload)
	if err != nil {
		t.Errorf("Webhook couldn't be created, error %v", err)
	}

	updatePayload := []WebhookField{
		WebhookField{
			Operation: "replace",
			Path:      "/event_types",
			Value: []interface{}{
				map[string]interface{}{
					"name": "PAYMENT.SALE.REFUNDED",
				},
			},
		},
	}

	_, err = c.UpdateWebhook(createdWebhook.ID, updatePayload)
	if err != nil {
		t.Errorf("Couldn't update webhook, error %v", err)
	}

	err = c.DeleteWebhook(createdWebhook.ID)
	if err != nil {
		t.Errorf("An error occurred while webhooks deletion, error %v", err)
	}
}

func TestListWebhooks(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.ListWebhooks(AncorTypeApplication)
	if err != nil {
		t.Errorf("Cannot registered list webhooks, error %v", err)
	}
}
