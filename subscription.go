package paypal

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type (
	SubscriptionBase struct {
		PlanID             string              `json:"plan_id"`
		StartTime          *JSONTime           `json:"start_time,omitempty"`
		EffectiveTime      *JSONTime           `json:"effective_time,omitempty"`
		Quantity           string              `json:"quantity,omitempty"`
		ShippingAmount     *Money              `json:"shipping_amount,omitempty"`
		Subscriber         *Subscriber         `json:"subscriber,omitempty"`
		AutoRenewal        bool                `json:"auto_renewal,omitempty"`
		ApplicationContext *ApplicationContext `json:"application_context,omitempty"`
		CustomID           string              `json:"custom_id,omitempty"`
		Plan               *PlanOverride       `json:"plan,omitempty"`
	}

	SubscriptionDetails struct {
		ID                           string             `json:"id,omitempty"`
		SubscriptionStatus           SubscriptionStatus `json:"status,omitempty"`
		SubscriptionStatusChangeNote string             `json:"status_change_note,omitempty"`
		StatusUpdateTime             time.Time          `json:"status_update_time,omitempty"`
	}

	Subscription struct {
		SubscriptionDetailResp
	}

	// SubscriptionDetailResp struct
	SubscriptionDetailResp struct {
		SubscriptionBase
		SubscriptionDetails
		BillingInfo BillingInfo `json:"billing_info,omitempty"` // not found in documentation
		SharedResponse
	}

	SubscriptionCaptureResponse struct {
		Status              SubscriptionTransactionStatus `json:"status"`
		Id                  string                        `json:"id"`
		AmountWithBreakdown AmountWithBreakdown           `json:"amount_with_breakdown"`
		PayerName           Name                          `json:"payer_name"`
		PayerEmail          string                        `json:"payer_email"`
		Time                time.Time                     `json:"time"`
	}

	//Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#definition-amount_with_breakdown
	AmountWithBreakdown struct {
		GrossAmount    Money `json:"gross_amount"`
		FeeAmount      Money `json:"fee_amount"`
		ShippingAmount Money `json:"shipping_amount"`
		TaxAmount      Money `json:"tax_amount"`
		NetAmount      Money `json:"net_amount"`
	}

	SubscriptionTransactionsParams struct {
		SubscriptionId string
		StartTime      time.Time
		EndTime        time.Time
	}

	SubscriptionTransactionsResponse struct {
		Transactions []SubscriptionCaptureResponse `json:"transactions"`
		SharedListResponse
	}

	CaptureReqeust struct {
		Note        string      `json:"note"`
		CaptureType CaptureType `json:"capture_type"`
		Amount      Money       `json:"amount"`
	}

	// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-plan_override
	PlanOverride struct {
		BillingCycles      []BillingCycleOverride       `json:"billing_cycles,omitempty"`
		PaymentPreferences *PaymentPreferencesOverride `json:"payment_preferences,omitempty"`
		Taxes              *TaxesOverride              `json:"taxes,omitempty"`
	}

	// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-payment_preferences_override
	PaymentPreferencesOverride struct {
		AutoBillOutstanding     bool                  `json:"auto_bill_outstanding,omitempty"`
		SetupFee                Money                 `json:"setup_fee,omitempty"`
		SetupFeeFailureAction   SetupFeeFailureAction `json:"setup_fee_failure_action,omitempty"`
		PaymentFailureThreshold int                   `json:"payment_failure_threshold,omitempty"`
	}

	// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-payment_preferences_override
	TaxesOverride struct {
		Percentage string `json:"percentage,omitempty"`
		Inclusive  *bool  `json:"inclusive,omitempty"`
	}

	// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-billing_cycle_override
	BillingCycleOverride struct {
		PricingScheme PricingScheme `json:"pricing_scheme,omitempty"`
		Sequence      *int          `json:"sequence,omitempty"`
		TotalCycles   *int          `json:"total_cycles,omitempty"`
	}
)

func (self *Subscription) GetUpdatePatch() []Patch {
	result := []Patch{
		{
			Operation: "replace",
			Path:      "/billing_info/outstanding_balance",
			Value:     self.BillingInfo.OutstandingBalance,
		},
	}
	return result
}

// CreateSubscriptionPlan creates a subscriptionPlan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_create
// Endpoint: POST /v1/billing/subscriptions
func (c *Client) CreateSubscription(ctx context.Context, newSubscription SubscriptionBase) (*SubscriptionDetailResp, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions"), newSubscription)
	req.Header.Add("Prefer", "return=representation")
	response := &SubscriptionDetailResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// UpdateSubscriptionPlan. updates a plan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_patch
// Endpoint: PATCH /v1/billing/subscriptions/:subscription_id
func (c *Client) UpdateSubscription(ctx context.Context, updatedSubscription Subscription) error {
	req, err := c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/billing/subscriptions/", updatedSubscription.ID), updatedSubscription.GetUpdatePatch())
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// GetSubscriptionDetails shows details for a subscription, by ID.
// Endpoint: GET /v1/billing/subscriptions/
func (c *Client) GetSubscriptionDetails(ctx context.Context, subscriptionID string) (*SubscriptionDetailResp, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/v1/billing/subscriptions/%s", c.APIBase, subscriptionID), nil)
	response := &SubscriptionDetailResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// Activates the subscription.
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_activate
// Endpoint: POST /v1/billing/subscriptions/{id}/activate
func (c *Client) ActivateSubscription(ctx context.Context, subscriptionId, activateReason string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/subscriptions/%s/activate", c.APIBase, subscriptionId), map[string]string{"reason": activateReason})
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Cancels the subscription.
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_cancel
// Endpoint: POST /v1/billing/subscriptions/{id}/cancel
func (c *Client) CancelSubscription(ctx context.Context, subscriptionId, cancelReason string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/subscriptions/%s/cancel", c.APIBase, subscriptionId), map[string]string{"reason": cancelReason})
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Captures an authorized payment from the subscriber on the subscription.
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_capture
// Endpoint: POST /v1/billing/subscriptions/{id}/capture
func (c *Client) CaptureSubscription(ctx context.Context, subscriptionId string, request CaptureReqeust) (*SubscriptionCaptureResponse, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/subscriptions/%s/capture", c.APIBase, subscriptionId), request)
	response := &SubscriptionCaptureResponse{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// Suspends the subscription.
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_suspend
// Endpoint: POST /v1/billing/subscriptions/{id}/suspend
func (c *Client) SuspendSubscription(ctx context.Context, subscriptionId, reason string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/subscriptions/%s/suspend", c.APIBase, subscriptionId), map[string]string{"reason": reason})
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Lists transactions for a subscription.
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_transactions
// Endpoint: GET /v1/billing/subscriptions/{id}/transactions
func (c *Client) GetSubscriptionTransactions(ctx context.Context, requestParams SubscriptionTransactionsParams) (*SubscriptionTransactionsResponse, error) {
	startTime := requestParams.StartTime.Format("2006-01-02T15:04:05Z")
	endTime := requestParams.EndTime.Format("2006-01-02T15:04:05Z")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/v1/billing/subscriptions/%s/transactions?start_time=%s&end_time=%s", c.APIBase, requestParams.SubscriptionId, startTime, endTime), nil)
	response := &SubscriptionTransactionsResponse{}
	if err != nil {
		return response, err
	}

	err = c.SendWithAuth(req, response)
	return response, err
}

// Revise plan or quantity of subscription
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_revise
// Endpoint: POST /v1/billing/subscriptions/{id}/revise
func (c *Client) ReviseSubscription(ctx context.Context, subscriptionId string, reviseSubscription SubscriptionBase) (*SubscriptionDetailResp, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/subscriptions/%s/revise", c.APIBase, subscriptionId), reviseSubscription)
	response := &SubscriptionDetailResp{}
	if err != nil {
		return response, err
	}

	req.Header.Add("Content-Type", "application/json")
	err = c.SendWithAuth(req, response)

	return response, err
}
