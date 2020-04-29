package paypal

import "fmt"

// ShowCapturedPayments shows details for a captured payment, by ID.
// Endpoint: GET /v2/payments/captures/{capture_id}
func (c *Client) ShowCapturedPayment(captureID string) (*Capture, error) {
	resp := &Capture{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/captures/"+captureID), nil)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithBasicAuth(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// RefundCapturedPayment refunds details for a captured payment, by ID.
// For a full refund, include an empty payload in the JSON request body. For a partial refund,
// include an amount object in the JSON request body.
// Endpoint: POST /v2/payments/captures/{capture_id}/refund
func (c *Client) RefundCapturedPayment(captureID string, body *RefundRequest) (*Refund, error) {
	resp := &Refund{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/captures/"+captureID+"/refund"), body)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithBasicAuth(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
