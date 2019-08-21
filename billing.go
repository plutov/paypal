package paypal

import (
	"bytes"
	"errors"
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

	// BillingPlanListParams struct
	BillingPlanListParams struct {
		Page          string `json:"page,omitempty"`           //Default: 0.
		Status        string `json:"status,omitempty"`         //Allowed values: CREATED, ACTIVE, INACTIVE, ALL.
		PageSize      string `json:"page_size,omitempty"`      //Default: 10.
		TotalRequired string `json:"total_required,omitempty"` //Default: no.

	}

	//BillingPlanListResp struct
	BillingPlanListResp struct {
		Plans      []BillingPlan `json:"plans,omitempty"`
		TotalItems string        `json:"total_items,omitempty"`
		TotalPages string        `json:"total_pages,omitempty"`
		Links      []Link        `json:"links,omitempty"`
	}
)

// CreateBillingPlan creates a billing plan in Paypal
// Endpoint: POST /v2/payments/billing-plans
func (c *Client) CreateBillingPlan(plan BillingPlan) (*CreateBillingResp, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/billing-plans"), plan)
	response := &CreateBillingResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// ActivatePlan activates a billing plan
// By default, a new plan is not activated
// Endpoint: PATCH /v2/payments/billing-plans/
func (c *Client) ActivatePlan(planID string) error {
	buf := bytes.NewBuffer([]byte(`[{"op":"replace","path":"/","value":{"state":"ACTIVE"}}]`))
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/billing-plans/"+planID), buf)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.ClientID, c.Secret)
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)
	return c.SendWithAuth(req, nil)
}

// CreateBillingAgreement creates an agreement for specified plan
// Endpoint: POST /v2/payments/billing-agreements
func (c *Client) CreateBillingAgreement(a BillingAgreement) (*CreateAgreementResp, error) {
	// PayPal needs only ID, so we will remove all fields except Plan ID
	a.Plan = BillingPlan{
		ID: a.Plan.ID,
	}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/billing-agreements"), a)
	response := &CreateAgreementResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// ExecuteApprovedAgreement - Use this call to execute (complete) a PayPal agreement that has been approved by the payer.
// Endpoint: POST /v2/payments/billing-agreements/token/agreement-execute
func (c *Client) ExecuteApprovedAgreement(token string) (*ExecuteAgreementResponse, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/billing-agreements/"+token+"/agreement-execute"), nil)
	response := &ExecuteAgreementResponse{}

	if err != nil {
		return response, err
	}

	req.SetBasicAuth(c.ClientID, c.Secret)
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	if response.ID == "" {
		return response, errors.New("Unable to execute agreement with token=" + token)
	}

	return response, err
}

// ListBillingPlans lists billing-plans
// Endpoint: GET /v2/payments/billing-plans
func (c *Client) ListBillingPlans(bplp BillingPlanListParams) (*BillingPlanListResp, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/billing-plans"), nil)
	q := req.URL.Query()
	q.Add("page", bplp.Page)
	q.Add("page_size", bplp.PageSize)
	q.Add("status", bplp.Status)
	q.Add("total_required", bplp.TotalRequired)
	req.URL.RawQuery = q.Encode()
	response := &BillingPlanListResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}
