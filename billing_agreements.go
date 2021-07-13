package paypal

import (
	"context"
	"fmt"
)

// CreateBillingAgreementToken - Use this call to create a billing agreement
// Endpoint: POST /v1/billing-agreements/agreement-tokens
func (c *Client) CreateBillingAgreementToken(
	ctx context.Context,
	name string,
	description string,
	startDate string,
	payer *Payer,
	plan *BillingPlan,
) (*BillingAgreementToken, error) {
	type createBARequest struct {
		Name        string       `json:"name"`
		Description string       `json:"description"`
		StartDate   string       `json:"start_date"`
		Payer       *Payer       `json:"payer"`
		Plan        *BillingPlan `json:"plan"`
	}

	billingAgreementToken := &BillingAgreementToken{}

	req, err := c.NewRequest(
		ctx,
		"POST",
		fmt.Sprintf("%s%s", c.APIBase, "/v1/billing-agreements/agreement-tokens"),
		createBARequest{Name: name, Description: description, StartDate: startDate, Payer: payer, Plan: plan})
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
