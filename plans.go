package paypal

import "fmt"

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