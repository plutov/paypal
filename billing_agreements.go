package paypal

import (
	"context"
	"fmt"
)

// CreateBillingAgreementToken - Use this call to create a billing agreement
// Endpoint: POST /v1/billing-agreements/agreement-tokens
func (c *Client) CreateBillingAgreementToken(
	ctx context.Context,
	description string,
	shippingAddress *ShippingAddress,
	payer *Payer,
	plan *BillingPlan,
) (*BillingAgreementToken, error) {
	type createBARequest struct {
		Description     string           `json:"description,omitempty"`
		ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
		Payer           *Payer           `json:"payer"`
		Plan            *BillingPlan     `json:"plan"`
	}

	billingAgreementToken := &BillingAgreementToken{}

	req, err := c.NewRequest(
		ctx,
		"POST",
		fmt.Sprintf("%s%s", c.APIBase, "/v1/billing-agreements/agreement-tokens"),
		createBARequest{Description: description, ShippingAddress: shippingAddress, Payer: payer, Plan: plan})
	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, billingAgreementToken); err != nil {
		return billingAgreementToken, err
	}

	return billingAgreementToken, nil
}

// CreateBillingAgreementFromToken - Use this call to create a billing agreement
// Endpoint: POST /v1/billing-agreements/agreements
func (c *Client) CreateBillingAgreementFromToken(
	ctx context.Context,
	tokenID string,
) (*BillingAgreement, error) {
	type createBARequest struct {
		TokenID string `json:"token_id"`
	}

	billingAgreement := &BillingAgreement{}

	req, err := c.NewRequest(
		ctx,
		"POST",
		fmt.Sprintf("%s%s", c.APIBase, "/v1/billing-agreements/agreements"),
		createBARequest{TokenID: tokenID})
	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, billingAgreement); err != nil {
		return billingAgreement, err
	}

	return billingAgreement, nil
}
