package paypal

import (
	"fmt"
	"strconv"
)

// CreatePlan creates plan
// Endpoint: POST /v1/billing/plans
func (c *Client) CreatePlan(plan *CreatePlan) (*Plan, error) {
	resp := &Plan{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans"), plan)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithBasicAuth(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ListAllPlans lists all plans
// Endpoint: GET /v1/billing/plans
func (c *Client) ListAllPlans(params *ListPlansParams) (*ListPlansResponse, error) {
	resp := &ListPlansResponse{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans"), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("product_id", params.ProductID)
	q.Add("page_size", strconv.FormatUint(params.PageSize, 10))
	q.Add("page", strconv.FormatUint(params.Page, 10))
	q.Add("total_required", strconv.FormatBool(params.TotalRequired))

	if err = c.SendWithBasicAuth(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ShowProduct shows details for a plan by ID
// Endpoint: GET /v1/billing/plans/{plan_id}
func (c *Client) ShowPlan(planID string) (*Plan, error) {
	resp := &Plan{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans/"+planID), nil)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithBasicAuth(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ActivatePlan updates plan status to active
// Endpoint: POST /v1/billing/plans/{plan_id}/activate
func (c *Client) ActivatePlan(planID string) error {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans/"+planID+"/activate"), nil)
	if err != nil {
		return err
	}

	return c.SendWithBasicAuth(req, nil)
}

// DeactivatePlan updates plan status to inactive
// Endpoint: POST /v1/billing/plans/{plan_id}/deactivate
func (c *Client) DeactivatePlan(planID string) error {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans/"+planID+"/deactivate"), nil)
	if err != nil {
		return err
	}

	return c.SendWithBasicAuth(req, nil)
}

// UpdatePlan updates plan by ID.
// Updates a plan with the CREATED or ACTIVE status. For an INACTIVE plan, you can make only status updates.
// You can update the following attributes and objects
//--------------------------------------------------------------------------
// | Attribute or Object 							| Operation			   |
// -------------------------------------------------------------------------
// | /description          							| add, replace, remove |
// | /payment_preferences/auto_bill_outstanding 	| add, replace, remove |
// | /taxes/percentage								| add, replace, remove |
// | /payment_preferences/payment_failure_threshold	| add, replace, remove |
// | /payment_preferences/setup_fee			        | add, replace, remove |
// | /payment_preferences/setup_fee_failure_action  | add, replace, remove |
// -------------------------------------------------------------------------
// Endpoint: PATCH /v1/billing/plan/{plan_id}
func (c *Client) UpdatePlan(planID string, patchObject []*PatchObject) error {
	req, err := c.NewRequest("PATCH", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plan/"+planID), patchObject)
	if err != nil {
		return err
	}

	return c.SendWithBasicAuth(req, nil)
}

// UpdatePricing updates pricing for a plan
// Endpoint: PATCH /v1/billing/plan/{plan_id}/update-pricing-schemas
func (c *Client) UpdatePricing(planID string,updatePricing UpdatePricingSchemasListRequest) error {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plan/"+planID+"/update-pricing-schemas"), updatePricing)
	if err != nil {
		return err
	}

	return c.SendWithBasicAuth(req, nil)
}