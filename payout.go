package paypalsdk

import (
	"fmt"
)

// CreateSinglePayout submits a payout with an asynchronous API call, which immediately returns the results of a PayPal payment.
// For email payout set RecipientType: "EMAIL" and receiver email into Receiver
// Endpoint: POST /v1/payments/payouts
func (c *Client) CreateSinglePayout(p Payout) (*PayoutResponse, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payouts"), p)
	if err != nil {
		return &PayoutResponse{}, err
	}

	response := &PayoutResponse{}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	return response, nil
}

// GetPayout shows the latest status of a batch payout along with the transaction status and other data for individual items.
// Also, returns IDs for the individual payout items. You can use these item IDs in other calls.
// Endpoint: GET /v1/payments/payouts/ID
func (c *Client) GetPayout(payoutBatchID string) (*PayoutResponse, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payouts/"+payoutBatchID), nil)

	if err != nil {
		return &PayoutResponse{}, err
	}

	response := &PayoutResponse{}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	return response, nil
}

// GetPayoutItem shows the details for a payout item.
// Use this call to review the current status of a previously unclaimed, or pending, payout item.
// Endpoint: GET /v1/payments/payouts-item/ID
func (c *Client) GetPayoutItem(payoutItemID string) (*PayoutItemResponse, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payouts-item/"+payoutItemID), nil)

	if err != nil {
		return &PayoutItemResponse{}, err
	}

	response := &PayoutItemResponse{}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	return response, nil
}

// CancelPayoutItem cancels an unclaimed Payout Item. If no one claims the unclaimed item within 30 days,
// the funds are automatically returned to the sender. Use this call to cancel the unclaimed item before the automatic 30-day refund.
// Endpoint: POST /v1/payments/payouts-item/ID/cancel
func (c *Client) CancelPayoutItem(payoutItemID string) (*PayoutItemResponse, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payouts-item/"+payoutItemID+"/cancel"), nil)

	if err != nil {
		return &PayoutItemResponse{}, err
	}

	response := &PayoutItemResponse{}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	return response, nil
}
