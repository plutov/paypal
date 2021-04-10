package paypal

import (
	"context"
	"fmt"
)

// GetCapturedPaymentDetails.
// Endpoint: GET /v1/payments/capture/:id
func (c *Client) GetCapturedPaymentDetails(ctx context.Context, id string) (*Capture, error) {
	res := &Capture{}

	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s%s", c.APIBase, "/v1/payments/capture/", id), nil)
	if err != nil {
		return res, err
	}

	if err = c.SendWithAuth(req, res); err != nil {
		return res, err
	}

	return res, nil
}
