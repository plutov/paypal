package paypal

import "fmt"

// CreateSubscription - Use this call to create a subscription
// Endpoint: POST /v1/billing/subscriptions
func (c *Client) CreateSubscription(subscription CreateSubscription) (*Subscription, error) {
	resp := &Subscription{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions"), subscription)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
