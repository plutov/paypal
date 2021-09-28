package paypal

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

// GetAuthorization returns an authorization by ID
// Endpoint: GET /v2/payments/authorizations/ID
func (c *Client) GetAuthorization(ctx context.Context, authID string) (*Authorization, error) {
	buf := bytes.NewBuffer([]byte(""))
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s%s", c.APIBase, "/v2/payments/authorizations/", authID), buf)
	auth := &Authorization{}

	if err != nil {
		return auth, err
	}

	err = c.SendWithAuth(req, auth)
	return auth, err
}

// CaptureAuthorization captures and process an existing authorization.
// To use this method, the original payment must have Intent set to "authorize"
// Endpoint: POST /v2/payments/authorizations/ID/capture
func (c *Client) CaptureAuthorization(ctx context.Context, authID string, paymentCaptureRequest *PaymentCaptureRequest) (*PaymentCaptureResponse, error) {
	return c.CaptureAuthorizationWithPaypalRequestId(ctx, authID, paymentCaptureRequest, "")
}

// CaptureAuthorization captures and process an existing authorization with idempotency.
// To use this method, the original payment must have Intent set to "authorize"
// Endpoint: POST /v2/payments/authorizations/ID/capture
func (c *Client) CaptureAuthorizationWithPaypalRequestId(ctx context.Context,
	authID string,
	paymentCaptureRequest *PaymentCaptureRequest,
	requestID string,
) (*PaymentCaptureResponse, error) {
	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/authorizations/"+authID+"/capture"), paymentCaptureRequest)
	paymentCaptureResponse := &PaymentCaptureResponse{}

	if err != nil {
		return paymentCaptureResponse, err
	}

	if requestID != "" {
		req.Header.Set("PayPal-Request-Id", requestID)
	}

	err = c.SendWithAuth(req, paymentCaptureResponse)
	return paymentCaptureResponse, err
}

// VoidAuthorization voids a previously authorized payment
// Endpoint: POST /v2/payments/authorizations/ID/void
func (c *Client) VoidAuthorization(ctx context.Context, authID string) (*Authorization, error) {
	buf := bytes.NewBuffer([]byte(""))
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/authorizations/"+authID+"/void"), buf)
	auth := &Authorization{}

	if err != nil {
		return auth, err
	}

	err = c.SendWithAuth(req, auth)
	return auth, err
}

// ReauthorizeAuthorization reauthorize a Paypal account payment.
// PayPal recommends reauthorizing payment after ~3 days
// Endpoint: POST /v2/payments/authorizations/ID/reauthorize
func (c *Client) ReauthorizeAuthorization(ctx context.Context, authID string, a *Amount) (*Authorization, error) {
	buf := bytes.NewBuffer([]byte(`{"amount":{"currency_code":"` + a.Currency + `","value":"` + a.Total + `"}}`))
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/authorizations/"+authID+"/reauthorize"), buf)
	auth := &Authorization{}

	if err != nil {
		return auth, err
	}

	err = c.SendWithAuth(req, auth)
	return auth, err
}
