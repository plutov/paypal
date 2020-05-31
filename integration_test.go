// +build integration

package paypal

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// All test values are defined here
var testClientID = "AXy9orp-CDaHhBZ9C78QHW2BKZpACgroqo85_NIOa9mIfJ9QnSVKzY-X_rivR_fTUUr6aLjcJsj6sDur"
var testSecret = "EBoIiUSkCKeSk49hHSgTem1qnjzzJgRQHDEHvGpzlLEf_nIoJd91xu8rPOBDCdR_UYNKVxJE-UgS2iCw"
var testUserID = "https://www.paypal.com/webapps/auth/identity/user/VBqgHcgZwb1PBs69ybjjXfIW86_Hr93aBvF_Rgbh2II"
var testCardID = "CARD-54E6956910402550WKGRL6EA"

var testProductId = "" // will be fetched in  func TestProduct(t *testing.T)
var testBillingPlan = "" // will be fetched in  func TestSubscriptionPlans(t *testing.T)

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

func TestProduct(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	//create a product
	productData := Product{
		Name:        "Test Product",
		Description: "A Test Product",
		Category:    ProductCategorySoftware,
		Type:        ProductTypeService,
		ImageUrl:    "https://example.com/image.png",
		HomeUrl:     "https://example.com",
	}

	productCreateResponse, err := c.CreateProduct(productData)
	assert.Equal(t, nil, err)

	testProductId = productCreateResponse.ID

	//update the product
	productData.ID = productCreateResponse.ID
	productData.Description = "Updated product"

	err = c.UpdateProduct(productData)
	assert.Equal(t, nil, err)

	//get product data
	productFetched, err := c.GetProduct(productData.ID)
	assert.Equal(t, nil, err)
	assert.Equal(t, productFetched.Description, "Updated product")

	//test that lising products have more than one product
	productList, err := c.ListProducts(nil)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, len(productList.Products), 0)
}

func TestSubscriptionPlans(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	//create a product
	newSubscriptionPlan := SubscriptionPlan{
		ProductId:   testProductId,
		Name:        "Test subscription plan",
		Status:      SubscriptionPlanStatusCreated,
		Description: "Integration test subscription plan",
		BillingCycles: []BillingCycle{
			{
				PricingScheme: PricingScheme{
					Version: 1,
					FixedPrice: Money{
						Currency: "EUR",
						Value:    "5",
					},
					CreateTime: time.Now(),
					UpdateTime: time.Now(),
				},
				Frequency: Frequency{
					IntervalUnit:  IntervalUnitYear,
					IntervalCount: 1,
				},
				TenureType:  TenureTypeRegular,
				Sequence:    1,
				TotalCycles: 0,
			},
		},
		PaymentPreferences: PaymentPreferences{
			AutoBillOutstanding:     false,
			SetupFee:                nil,
			SetupFeeFailureAction:   SetupFeeFailureActionCancel,
			PaymentFailureThreshold: 0,
		},
		Taxes: Taxes{
			Percentage: "19",
			Inclusive:  false,
		},
		QuantitySupported: false,
	}

	//test create new plan
	planCreateResponse, err := c.CreateSubscriptionPlan(newSubscriptionPlan)
	assert.Equal(t, nil, err)
	testBillingPlan = planCreateResponse.ID // for next test

	//test update the newly created plan
	newSubscriptionPlan.ID = planCreateResponse.ID
	newSubscriptionPlan.Description = "updated description"
	err = c.UpdateSubscriptionPlan(newSubscriptionPlan)
	assert.Equal(t, nil, err)

	//test get plan information
	existingPlan, err := c.GetSubscriptionPlan(newSubscriptionPlan.ID)
	assert.Equal(t, nil, err)
	assert.Equal(t, newSubscriptionPlan.Description, existingPlan.Description)

	//test activate plan
	err = c.ActivateSubscriptionPlan(newSubscriptionPlan.ID)
	assert.Equal(t, nil, err)

	//test deactivate plan
	err = c.DeactivateSubscriptionPlans(newSubscriptionPlan.ID)
	assert.Equal(t, nil, err)

	//reactivate this plan for next next (subscription)
	err = c.ActivateSubscriptionPlan(newSubscriptionPlan.ID)
	assert.Equal(t, nil, err)

	//test upadte plan pricing
	err = c.UpdateSubscriptionPlanPricing(newSubscriptionPlan.ID, []PricingSchemeUpdate{
		{
			BillingCycleSequence: 1,
			PricingScheme: PricingScheme{
				Version: 1,
				FixedPrice: Money{
					Currency: "EUR",
					Value:    "6",
				},
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			},
		},
	})
	assert.Equal(t, nil, err)

	//test update pricing scheme
	updatedPricingPlan, err := c.GetSubscriptionPlan(newSubscriptionPlan.ID)
	assert.Equal(t, nil, err)
	assert.Equal(t, "6.0", updatedPricingPlan.BillingCycles[0].PricingScheme.FixedPrice.Value)

}

func TestSubscription(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	newSubscription := SubscriptionBase{
		PlanID:             testBillingPlan,
	}

	//create new subscription
	newSubResponse, err := c.CreateSubscription(newSubscription)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", newSubResponse.ID)

	//get subscription details
	subDetails, err := c.GetSubscriptionDetails(newSubResponse.ID)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", subDetails.ID)

}
