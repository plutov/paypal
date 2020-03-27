package paypal

import (
	"fmt"
	"strconv"
)

// CreatePlan - Use this call to create a plan
// Endpoint: POST /v1/billing/plans
func (c *Client) CreatePlan(plan *CreatePlan) (*Plan, error) {
	resp := &Plan{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans"), plan)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ListAllPlans lists all products
// Endpoint: GET /v1/billing/plans
func (c *Client) ListAllPlans(params *ListPlansParams) (*ListPlansResponse, error) {
	resp := &ListPlansResponse{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans"), nil)
	q := req.URL.Query()
	q.Add("product_id", params.ProductID)
	q.Add("page_size", strconv.FormatUint(params.PageSize, 10))
	q.Add("page", strconv.FormatUint(params.Page, 10))
	q.Add("total_required", strconv.FormatBool(params.TotalRequired))
	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ShowProduct shows details for a product by ID
// Endpoint: GET /v1/billing/plans/{plan_id}
func (c *Client) ShowPlan(planID string) (*Plan, error) {
	resp := &Plan{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans/"+planID), nil)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ActivatePlan shows details for a product by ID
// Endpoint: GET /v1/billing/plans/{plan_id}/activate
func (c *Client) ActivatePlan(planID string) error {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans/"+planID+"/activate"), nil)
	if err != nil {
		return err
	}

	return c.SendWithAuth(req, nil)
}
