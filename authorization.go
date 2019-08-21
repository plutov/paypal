package paypal

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

// GetAuthorization returns an authorization by ID
// Endpoint: GET /v2/payments/authorization/ID
func (c *Client) GetAuthorization(authID string) (*Authorization, error) {
	buf := bytes.NewBuffer([]byte(""))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", c.APIBase, "/v2/payments/authorization/", authID), buf)
	auth := &Authorization{}

	if err != nil {
		return auth, err
	}

	err = c.SendWithAuth(req, auth)
	return auth, err
}

// CaptureAuthorization captures and process an existing authorization.
// To use this method, the original payment must have Intent set to "authorize"
// Endpoint: POST /v2/payments/authorization/ID/capture
func (c *Client) CaptureAuthorization(authID string, a *Amount, isFinalCapture bool) (*Capture, error) {
	isFinalStr := strconv.FormatBool(isFinalCapture)

	buf := bytes.NewBuffer([]byte(`{"amount":{"currency":"` + a.Currency + `,"total":"` + a.Total + `"},"is_final_capture":` + isFinalStr + `}`))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/authorization/"+authID+"/capture"), buf)
	capture := &Capture{}

	if err != nil {
		return capture, err
	}

	err = c.SendWithAuth(req, capture)
	return capture, err
}

// VoidAuthorization voids a previously authorized payment
// Endpoint: POST /v2/payments/authorization/ID/void
func (c *Client) VoidAuthorization(authID string) (*Authorization, error) {
	buf := bytes.NewBuffer([]byte(""))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/authorization/"+authID+"/void"), buf)
	auth := &Authorization{}

	if err != nil {
		return auth, err
	}

	err = c.SendWithAuth(req, auth)
	return auth, err
}

// ReauthorizeAuthorization reauthorize a Paypal account payment.
// PayPal recommends to reauthorize payment after ~3 days
// Endpoint: POST /v2/payments/authorization/ID/reauthorize
func (c *Client) ReauthorizeAuthorization(authID string, a *Amount) (*Authorization, error) {
	buf := bytes.NewBuffer([]byte(`{"amount":{"currency":"` + a.Currency + `","total":"` + a.Total + `"}}`))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/authorization/"+authID+"/reauthorize"), buf)
	auth := &Authorization{}

	if err != nil {
		return auth, err
	}

	err = c.SendWithAuth(req, auth)
	return auth, err
}
