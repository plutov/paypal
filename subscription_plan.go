package paypal

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type (
	// SubscriptionDetailResp struct
	SubscriptionPlan struct {
		ID                 string                 `json:"id,omitempty"`
		ProductId          string                 `json:"product_id"`
		Name               string                 `json:"name"`
		Status             SubscriptionPlanStatus `json:"status"`
		Description        string                 `json:"description,omitempty"`
		BillingCycles      []BillingCycle         `json:"billing_cycles"`
		PaymentPreferences *PaymentPreferences    `json:"payment_preferences"`
		Taxes              *Taxes                 `json:"taxes"`
		QuantitySupported  bool                   `json:"quantity_supported"` //Indicates whether you can subscribe to this plan by providing a quantity for the goods or service.
	}

	// Doc https://developer.paypal.com/docs/api/subscriptions/v1/#definition-billing_cycle
	BillingCycle struct {
		PricingScheme PricingScheme `json:"pricing_scheme"` // The active pricing scheme for this billing cycle. A free trial billing cycle does not require a pricing scheme.
		Frequency     Frequency     `json:"frequency"`      // The frequency details for this billing cycle.
		TenureType    TenureType    `json:"tenure_type"`    // The tenure type of the billing cycle. In case of a plan having trial cycle, only 2 trial cycles are allowed per plan. The possible values are:
		Sequence      int           `json:"sequence"`       // The order in which this cycle is to run among other billing cycles. For example, a trial billing cycle has a sequence of 1 while a regular billing cycle has a sequence of 2, so that trial cycle runs before the regular cycle.
		TotalCycles   int           `json:"total_cycles"`   // The number of times this billing cycle gets executed. Trial billing cycles can only be executed a finite number of times (value between 1 and 999 for total_cycles). Regular billing cycles can be executed infinite times (value of 0 for total_cycles) or a finite number of times (value between 1 and 999 for total_cycles).
	}

	// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#definition-payment_preferences
	PaymentPreferences struct {
		AutoBillOutstanding     bool                  `json:"auto_bill_outstanding"`
		SetupFee                *Money                `json:"setup_fee"`
		SetupFeeFailureAction   SetupFeeFailureAction `json:"setup_fee_failure_action"`
		PaymentFailureThreshold int                   `json:"payment_failure_threshold"`
	}

	PricingScheme struct {
		Version    int       `json:"version"`
		FixedPrice Money     `json:"fixed_price"`
		CreateTime time.Time `json:"create_time"`
		UpdateTime time.Time `json:"update_time"`
	}

	PricingSchemeUpdateRequest struct {
		Schemes []PricingSchemeUpdate `json:"pricing_schemes"`
	}

	PricingSchemeUpdate struct {
		BillingCycleSequence int           `json:"billing_cycle_sequence"`
		PricingScheme        PricingScheme `json:"pricing_scheme"`
	}

	//doc: https://developer.paypal.com/docs/api/subscriptions/v1/#definition-frequency
	Frequency struct {
		IntervalUnit  IntervalUnit `json:"interval_unit"`
		IntervalCount int          `json:"interval_count"` //different per unit. check documentation
	}

	Taxes struct {
		Percentage string `json:"percentage"`
		Inclusive  bool   `json:"inclusive"`
	}

	CreateSubscriptionPlanResponse struct {
		SubscriptionPlan
		SharedResponse
	}

	SubscriptionPlanListParameters struct {
		ProductId string `json:"product_id"`
		PlanIds   string `json:"plan_ids"` // Filters the response by list of plan IDs. Filter supports upto 10 plan IDs.
		ListParams
	}

	ListSubscriptionPlansResponse struct {
		Plans []SubscriptionPlan `json:"plans"`
		SharedListResponse
	}
)

func (p *SubscriptionPlan) GetUpdatePatch() []Patch {
	result := []Patch{
		{
			Operation: "replace",
			Path:      "/description",
			Value:     p.Description,
		},
	}

	if p.Taxes != nil {
		result = append(result, Patch{
			Operation: "replace",
			Path:      "/taxes/percentage",
			Value:     p.Taxes.Percentage,
		})
	}

	if p.PaymentPreferences != nil {
		if p.PaymentPreferences.SetupFee != nil {
			result = append(result, Patch{
				Operation: "replace",
				Path:      "/payment_preferences/setup_fee",
				Value:     p.PaymentPreferences.SetupFee,
			},
			)
		}

		result = append(result, []Patch{{
			Operation: "replace",
			Path:      "/payment_preferences/auto_bill_outstanding",
			Value:     p.PaymentPreferences.AutoBillOutstanding,
		},
			{
				Operation: "replace",
				Path:      "/payment_preferences/payment_failure_threshold",
				Value:     p.PaymentPreferences.PaymentFailureThreshold,
			},
			{
				Operation: "replace",
				Path:      "/payment_preferences/setup_fee_failure_action",
				Value:     p.PaymentPreferences.SetupFeeFailureAction,
			}}...)
	}

	return result
}

// CreateSubscriptionPlan creates a subscriptionPlan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_create
// Endpoint: POST /v1/billing/plans
func (c *Client) CreateSubscriptionPlan(ctx context.Context, newPlan SubscriptionPlan) (*CreateSubscriptionPlanResponse, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans"), newPlan)
	response := &CreateSubscriptionPlanResponse{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// UpdateSubscriptionPlan. updates a plan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_patch
// Endpoint: PATCH /v1/billing/plans/:plan_id
func (c *Client) UpdateSubscriptionPlan(ctx context.Context, updatedPlan SubscriptionPlan) error {
	req, err := c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/billing/plans/", updatedPlan.ID), updatedPlan.GetUpdatePatch())
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// UpdateSubscriptionPlan. updates a plan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_get
// Endpoint: GET /v1/billing/plans/:plan_id
func (c *Client) GetSubscriptionPlan(ctx context.Context, planId string) (*SubscriptionPlan, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/billing/plans/", planId), nil)
	response := &SubscriptionPlan{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// List all plans
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_list
// Endpoint: GET /v1/billing/plans
func (c *Client) ListSubscriptionPlans(ctx context.Context, params *SubscriptionPlanListParameters) (*ListSubscriptionPlansResponse, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans"), nil)
	response := &ListSubscriptionPlansResponse{}
	if err != nil {
		return response, err
	}

	if params != nil {
		q := req.URL.Query()
		q.Add("page", params.Page)
		q.Add("page_size", params.PageSize)
		q.Add("total_required", params.TotalRequired)
		q.Add("product_id", params.ProductId)
		q.Add("plan_ids", params.PlanIds)
		req.URL.RawQuery = q.Encode()
	}

	err = c.SendWithAuth(req, response)
	return response, err
}

// Activates a plan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_activate
// Endpoint: POST /v1/billing/plans/{id}/activate
func (c *Client) ActivateSubscriptionPlan(ctx context.Context, planId string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/plans/%s/activate", c.APIBase, planId), nil)
	if err != nil {
		return err
	}

	err = c.SendWithAuth(req, nil)
	return err
}

// Deactivates a plan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_deactivate
// Endpoint: POST /v1/billing/plans/{id}/deactivate
func (c *Client) DeactivateSubscriptionPlans(ctx context.Context, planId string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/plans/%s/deactivate", c.APIBase, planId), nil)
	if err != nil {
		return err
	}

	err = c.SendWithAuth(req, nil)
	return err
}

// Updates pricing for a plan. For example, you can update a regular billing cycle from $5 per month to $7 per month.
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_update-pricing-schemes
// Endpoint: POST /v1/billing/plans/{id}/update-pricing-schemes
func (c *Client) UpdateSubscriptionPlanPricing(ctx context.Context, planId string, pricingSchemes []PricingSchemeUpdate) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/plans/%s/update-pricing-schemes", c.APIBase, planId), PricingSchemeUpdateRequest{
		Schemes: pricingSchemes,
	})
	if err != nil {
		return err
	}

	err = c.SendWithAuth(req, nil)
	return err
}
