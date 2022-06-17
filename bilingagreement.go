package paypal

import (
	"fmt"
)

type AgreementRequest struct {
	Note string `json:"note"`
}

// SuspendAgreement returns status code and message
// pauses billing-agreement
// Endpoint: GET /v1/payments/billing-agreements/subID/suspend"
func (c *Client) SuspendAgreement(bAID string, agr AgreementRequest) (*DefaultResponse, error) {

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-agreements/"+bAID+"/suspend"), agr)
	resp := &DefaultResponse{}

	if err != nil {
		return resp, err
	}

	err = c.SendWithBasicAuth(req, resp)
	return resp, err
}

// ReActivateAgreement returns status code and message
// re-activate billing-agreement
// Endpoint: GET /v1/payments/billing-agreements/subID/suspend"
func (c *Client) ReActivateAgreement(bAID string, agr AgreementRequest) (*DefaultResponse, error) {

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-agreements/"+bAID+"/re-activate"), agr)
	resp := &DefaultResponse{}

	if err != nil {
		return resp, err
	}

	err = c.SendWithBasicAuth(req, resp)
	return resp, err
}

// CancelAgreement returns status code and message
// cancel billing-agreement
// Endpoint: GET /v1/payments/billing-agreements/subID/cancel"
func (c *Client) CancelAgreement(bAID string, agr AgreementRequest) (*DefaultResponse, error) {

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-agreements/"+bAID+"/cancel"), agr)
	resp := &DefaultResponse{}

	if err != nil {
		return resp, err
	}

	err = c.SendWithBasicAuth(req, resp)
	return resp, err
}
