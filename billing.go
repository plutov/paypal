package paypalsdk

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type (
	// CreateBillingResp struct
	CreateBillingResp struct {
		ID                  string              `json:"id,omitempty"`
		State               string              `json:"state,omitempty"`
		PaymentDefinitions  []PaymentDefinition `json:"payment_definitions,omitempty"`
		MerchantPreferences MerchantPreferences `json:"merchant_preferences,omitempty"`
		CreateTime          time.Time           `json:"create_time,omitempty"`
		UpdateTime          time.Time           `json:"update_time,omitempty"`
		Links               []Link              `json:"links,omitempty"`
	}

	// CreateAgreementResp struct
	CreateAgreementResp struct {
		Name        string      `json:"name,omitempty"`
		Description string      `json:"description,omitempty"`
		Plan        BillingPlan `json:"plan,omitempty"`
		Links       []Link      `json:"links,omitempty"`
		StartTime   time.Time   `json:"start_time,omitempty"`
	}
)

// CreateBillingPlan creates a billing plan in Paypal
// Endpoint: POST /v1/payments/billing-plans
func (c *Client) CreateBillingPlan(plan BillingPlan) (*CreateBillingResp, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-plans"), plan)
	response := &CreateBillingResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// Activates a billing plan
// By default, a new plan is not activated
// Endpoint: PATCH /v1/payments/billing-plans/
func (c *Client) ActivatePlan(planID string) error {
	buf := bytes.NewBuffer([]byte("[{\"op\":\"replace\",\"path\":\"/\",\"value\":{\"state\":\"ACTIVE\"}}]"))
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-plans/"+planID), buf)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.ClientID, c.Secret)
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)
	return c.SendWithAuth(req, nil)
}

// Creates an agreement for specified plan
// Endpoint: POST /v1/payments/billing-agreements
func (c *Client) CreateBillingAgreement(a BillingAgreement) (*CreateAgreementResp, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-agreements"), a)
	response := &CreateAgreementResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}
