package paypal

import "fmt"

// CreateSubscription - Use this call to create a subscription
// Endpoint: POST /v1/billing/subscriptions
func (c *Client) CreateSubscription(subscription CreateSubscriptionRequest) (*Subscription, error) {
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

// ShowSubscription shows details for a subscription by ID
// Endpoint: GET /v1/billing/subscriptions/{subscription_id}
func (c *Client) ShowSubscription(subscriptionID string, params *ShowSubscriptionRequest) (*Subscription, error) {
	resp := &Subscription{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions/"+subscriptionID), nil)
	q := req.URL.Query()
	q.Add("fields", params.Fields)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ActivatePlan activates subscription by ID
// Endpoint: POST /v1/billing/subscriptions/{subscription_id}/activate
func (c *Client) ActivateSubscription(subscriptionID string, body UpdateSubscriptionStatusRequest) error {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions/"+subscriptionID+"/activate"), body)
	if err != nil {
		return err
	}

	return c.SendWithAuth(req, nil)
}

// CancelSubscription cancels subscription by ID
// Endpoint: POST /v1/billing/subscriptions/{subscription_id}/cancel
func (c *Client) CancelSubscription(subscriptionID string, body UpdateSubscriptionStatusRequest) error {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions/"+subscriptionID+"/cancel"), body)
	if err != nil {
		return err
	}

	return c.SendWithAuth(req, nil)
}

// SuspendSubscription suspends subscription by ID
// Endpoint: POST /v1/billing/subscriptions/{subscription_id}/suspend
func (c *Client) SuspendSubscription(subscriptionID string, body UpdateSubscriptionStatusRequest) error {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions/"+subscriptionID+"/suspend"), body)
	if err != nil {
		return err
	}

	return c.SendWithAuth(req, nil)
}

// UpdateSubscription updates subscription by ID
// You can update the following attributes and objects
//-----------------------------------------------
// | Attribute or Object | Operation			|
// ----------------------------------------------
// | path to the object  | add, replace, remove |
// | (delimiter: /)      |                      |
// ----------------------------------------------
// Endpoint: PATCH /v1/billing/subscriptions/{subscription_id}
func (c *Client) UpdateSubscription(subscriptionID string, body []*PatchObject) error {
	req, err := c.NewRequest("PATCH", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions/"+subscriptionID), body)
	if err != nil {
		return err
	}

	return c.SendWithAuth(req, nil)
}

// CaptureAuthorizedPaymentOnSubscription captures an authorized payment from the subscriber on the subscription
// Endpoint: POST /v1/billing/subscriptions/{subscription_id}/capture
func (c *Client) CaptureAuthorizedPaymentOnSubscription(subscriptionID string, body *CaptureAuthorizedPaymentOnSubscriptionRequest) error {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions/"+subscriptionID+"/capture"), body)
	if err != nil {
		return err
	}

	return c.SendWithAuth(req, nil)
}

// ListTransactionsForSubscription lists transactions for subscription
// Endpoint: GET /v1/billing/subscriptions/{subscription_id}/transactions
func (c *Client) ListTransactionsForSubscription(subscriptionID string, params *ListTransactionsForSubscriptionRequest) (*TransactionsList, error) {
	resp := &TransactionsList{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions/"+subscriptionID+"/transactions"), nil)
	q := req.URL.Query()
	q.Add("start_time", params.StartTime)
	q.Add("end_time", params.EndTime)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ReviseSubscription updates the quantity of the product or service in a subscription. You can also use this method to switch the plan
// and update the shipping_amount, shipping_address values for the subscription. This type of update requires the buyer's consent.
// Endpoint: POST /v1/billing/subscriptions/{subscription_id}/revise
func (c *Client) ReviseSubscription(subscriptionID string, body *ReviseSubscriptionRequest) (*ReviseSubscriptionResponse, error) {
	resp := &ReviseSubscriptionResponse{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions/"+subscriptionID+"/revise"), body)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}