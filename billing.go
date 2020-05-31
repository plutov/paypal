package paypal

import (
	"bytes"
	"encoding/json"
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
		ListParams
		Status string `json:"status,omitempty"` //Allowed values: CREATED, ACTIVE, INACTIVE, ALL.
	}

	//BillingPlanListResp struct
	BillingPlanListResp struct {
		SharedListResponse
		Plans []BillingPlan `json:"plans,omitempty"`
	}
)

// CreateBillingPlan creates a billing plan in Paypal
// Endpoint: POST /v1/payments/billing-plans
func (c *Client) CreateBillingPlan(plan BillingPlan) (*CreateBillingResp, error) {
	req, err := c.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-plans"), plan)
	response := &CreateBillingResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// UpdateBillingPlan updates values inside a billing plan
// Endpoint: PATCH /v1/payments/billing-plans
func (c *Client) UpdateBillingPlan(planId string, pathValues map[string]map[string]interface{}) error {
	patchData := []Patch{}
	for path, data := range pathValues {
		patchData = append(patchData, Patch{
			Operation: "replace",
			Path:      path,
			Value:     data,
		})
	}

	jsonData, err := json.Marshal(patchData)
	buf := bytes.NewBuffer(jsonData)
	req, err := c.NewRequest(http.MethodPatch, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/payments/billing-plans/", planId), buf)
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// ActivatePlan activates a billing plan
// By default, a new plan is not activated
// Endpoint: PATCH /v1/payments/billing-plans/
func (c *Client) ActivatePlan(planID string) error {
	return c.UpdateBillingPlan(planID, map[string]map[string]interface{}{
		"/": {"state": BillingPlanStatusActive},
	})
}

// CreateBillingAgreement creates an agreement for specified plan
// Endpoint: POST /v1/payments/billing-agreements
func (c *Client) CreateBillingAgreement(a BillingAgreement) (*CreateAgreementResp, error) {
	// PayPal needs only ID, so we will remove all fields except Plan ID
	a.Plan = BillingPlan{
		ID: a.Plan.ID,
	}

	req, err := c.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-agreements"), a)
	response := &CreateAgreementResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// ExecuteApprovedAgreement - Use this call to execute (complete) a PayPal agreement that has been approved by the payer.
// Endpoint: POST /v1/payments/billing-agreements/token/agreement-execute
func (c *Client) ExecuteApprovedAgreement(token string) (*ExecuteAgreementResponse, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/payments/billing-agreements/%s/agreement-execute", c.APIBase, token), nil)
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
// Endpoint: GET /v1/payments/billing-plans
func (c *Client) ListBillingPlans(bplp BillingPlanListParams) (*BillingPlanListResp, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-plans"), nil)
	response := &BillingPlanListResp{}
	if err != nil {
		return response, err
	}

	q := req.URL.Query()
	q.Add("page", bplp.Page)
	q.Add("page_size", bplp.PageSize)
	q.Add("status", bplp.Status)
	q.Add("total_required", bplp.TotalRequired)
	req.URL.RawQuery = q.Encode()

	err = c.SendWithAuth(req, response)
	return response, err
}
