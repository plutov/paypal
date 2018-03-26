package paypalsdk

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type (
	// CreateAgreementResp struct
	CreateAgreementResp struct {
		ID          string      `json:"id,omitempty"`
		Name        string      `json:"name,omitempty"`
		Description string      `json:"description,omitempty"`
		Plan        BillingPlan `json:"plan,omitempty"`
		Links       []Link      `json:"links,omitempty"`
		StartTime   time.Time   `json:"start_time,omitempty"`
	}

	// ListBillingPlansResp struct
	ListBillingPlansResp struct {
		TotalItems string        `json:"total_items,omitempty"`
		TotalPages string        `json:"total_pages,omitempty"`
		Plans      []BillingPlan `json:"plans,omitempty"`
		Links      []Link        `json:"links,omitempty"`
	}
)

// CreateBillingPlan creates a billing plan in Paypal.
// Endpoint: POST /v1/payments/billing-plans
func (c *Client) CreateBillingPlan(plan BillingPlan) (*BillingPlan, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-plans"), plan)
	response := &BillingPlan{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// ActivatePlan activates a billing plan.
// By default, a new plan is not activated.
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

// CreateBillingAgreement creates an agreement for specified plan.
// Endpoint: POST /v1/payments/billing-agreements
func (c *Client) CreateBillingAgreement(a BillingAgreement) (*CreateAgreementResp, error) {
	// PayPal needs only ID, so we will remove all fields except Plan ID
	a.Plan = BillingPlan{
		ID: a.Plan.ID,
	}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-agreements"), a)
	response := &CreateAgreementResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// DeletePlan deletes a billing plan.
// Endpoint: PATCH /v1/payments/billing-plans/
func (c *Client) DeletePlan(planID string) error {
	buf := bytes.NewBuffer([]byte("[{\"op\":\"replace\",\"path\":\"/\",\"value\":{\"state\":\"DELETED\"}}]"))
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-plans/"+planID), buf)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.ClientID, c.Secret)
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)
	return c.SendWithAuth(req, nil)
}

// ExecuteApprovedAgreement - Use this call to execute (complete) a PayPal agreement that has been approved by the payer.
// Endpoint: POST /v1/payments/billing-agreements/token/agreement-execute
func (c *Client) ExecuteApprovedAgreement(token string) (*ExecuteAgreementResponse, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-agreements/"+token+"/agreement-execute"), nil)
	if err != nil {
		return &ExecuteAgreementResponse{}, err
	}

	req.SetBasicAuth(c.ClientID, c.Secret)
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)

	e := ExecuteAgreementResponse{}

	if err = c.SendWithAuth(req, &e); err != nil {
		return &e, err
	}

	if e.ID == "" {
		return &e, errors.New("Unable to execute agreement with token=" + token)
	}

	return &e, err
}

// ListBillingPlans - Lists billing plans.
// Valid values for status: "CREATED", "ACTIVE", "INACTIVE".
// Endpoint: GET /v1/payments/billing-plans/
func (c *Client) ListBillingPlans(status interface{}, page interface{}) (*ListBillingPlansResp, error) {
	if status == nil {
		status = "CREATED"
	}
	if page == nil {
		page = "0"
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-plans?total_required=yes&status="+status.(string)+"&page="+page.(string)), nil)
	if err != nil {
		return &ListBillingPlansResp{}, err
	}

	req.SetBasicAuth(c.ClientID, c.Secret)
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)

	l := ListBillingPlansResp{}

	err = c.SendWithAuth(req, &l)

	return &l, err
}