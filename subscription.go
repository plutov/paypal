package paypal

import (
	"fmt"
	"net/http"
	"time"
)

type (
	// SubscriptionDetailResp struct
	SubscriptionDetailResp struct {
		ID               string         `json:"id,omitempty"`
		PlanID           string         `json:"plan_id,omitempty"`
		StartTime        time.Time      `json:"start_time,omitempty"`
		Quantity         string         `json:"quantity,omitempty"`
		ShippingAmount   ShippingAmount `json:"shipping_amount,omitempty"`
		Subscriber       Subscriber     `json:"subscriber,omitempty"`
		BillingInfo      BillingInfo    `json:"billing_info,omitempty"`
		CreateTime       time.Time      `json:"create_time,omitempty"`
		UpdateTime       time.Time      `json:"update_time,omitempty"`
		Links            []Link         `json:"links,omitempty"`
		Status           string         `json:"status,omitempty"`
		StatusUpdateTime time.Time      `json:"status_update_time,omitempty"`
	}
)

// GetSubscriptionDetails shows details for a subscription, by ID.
// Endpoint: GET /v1/billing/subscriptions/
func (c *Client) GetSubscriptionDetails(subscriptionID string) (*SubscriptionDetailResp, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/billing/subscriptions/%s", c.APIBase, subscriptionID), nil)
	response := &SubscriptionDetailResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}
