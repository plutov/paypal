package paypal

import "fmt"

// ShowRefund shows details for a refund by ID
// Endpoint: GET /v2/payments/refunds/{refund_id}
func (c *Client) ShowRefund(refundID string) (*Refund, error) {
	resp := &Refund{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/refunds/"+refundID), nil)
	if err != nil {
		return nil, err
	}

	err = c.SendWithBasicAuth(req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
