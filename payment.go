package paypalsdk

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

// ListPaymentsResp slice of payments
type ListPaymentsResp struct {
	Payments []Payment `json:"payments"`
}

// CreatePaymentResp contains Payment Info and Links slice
type CreatePaymentResp struct {
	*Payment
	Links []Link `json:"links"`
}

// CreateDirectPaypalPayment sends request to create a payment with payment_method=paypal
// CreatePayment is more common function for any kind of payment
// Endpoint: POST /v1/payments/payment
func (c *Client) CreateDirectPaypalPayment(amount Amount, redirectURI string, cancelURI string, description string) (*PaymentResponse, error) {
	buf := bytes.NewBuffer([]byte("{\"intent\":\"sale\",\"payer\":{\"payment_method\":\"paypal\"}," +
		"\"transactions\":[{\"amount\":{\"total\":\"" + amount.Total +
		"\",\"currency\":\"" + amount.Currency + "\"},\"description\":\"" + description + "\"}],\"redirect_urls\":{\"return_url\":\"" +
		redirectURI + "\",\"cancel_url\":\"" + cancelURI + "\"}}"))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payment"), buf)
	if err != nil {
		return &PaymentResponse{}, err
	}

	req.SetBasicAuth(c.ClientID, c.Secret)
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)

	p := PaymentResponse{}

	if err = c.SendWithAuth(req, &p); err != nil {
		return &p, err
	}

	if p.ID == "" {
		return &p, errors.New("Unable to create payment with this access token")
	}

	return &p, err
}

// CreatePayment creates a payment in Paypal
// Depending on the payment_method and the funding_instrument, you can use the payment resource for direct credit card payments, stored credit card payments, or PayPal account payments.
// Endpoint: POST /v1/payments/payment
func (c *Client) CreatePayment(p Payment) (*CreatePaymentResp, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payment"), p)
	if err != nil {
		return &CreatePaymentResp{}, err
	}

	response := &CreatePaymentResp{}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	return response, nil
}

// ExecuteApprovedPayment - Use this call to execute (complete) a PayPal payment that has been approved by the payer. You can optionally update transaction information when executing the payment by passing in one or more transactions.
// Endpoint: POST /v1/payments/payment/paymentID/execute
func (c *Client) ExecuteApprovedPayment(paymentID string, payerID string) (*ExecuteResponse, error) {
	buf := bytes.NewBuffer([]byte("{\"payer_id\":\"" + payerID + "\"}"))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payment/"+paymentID+"/execute"), buf)
	if err != nil {
		return &ExecuteResponse{}, err
	}

	req.SetBasicAuth(c.ClientID, c.Secret)
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)

	e := ExecuteResponse{}

	if err = c.SendWithAuth(req, &e); err != nil {
		return &e, err
	}

	if e.ID == "" {
		return &e, errors.New("Unable to execute payment with paymentID=" + paymentID)
	}

	return &e, err
}

// GetPayment gets a payment from PayPal
// Endpoint: GET /v1/payments/payment/ID
func (c *Client) GetPayment(paymentID string) (*Payment, error) {
	p := Payment{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payment/"+paymentID), nil)
	if err != nil {
		return &p, err
	}

	if err = c.SendWithAuth(req, &p); err != nil {
		return &p, err
	}

	if p.ID == "" {
		return &p, errors.New("Unable to get payment with paymentID=" + paymentID)
	}

	return &p, nil
}

// PatchPayment modifies some fields of a payment prior to execution
// Endpoint: PATCH /v1/payments/payment/ID
func (c *Client) PatchPayment(paymentID string, p []PaymentPatch) (*Payment, error) {
	req, err := c.NewRequest("PATCH", fmt.Sprintf("%s/v1/payments/payment/%s", c.APIBase, paymentID), p)
	if err != nil {
		return nil, err
	}

	response := Payment{}

	if err = c.SendWithAuth(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetPayments retrieve payments resources from Paypal
// Endpoint: GET /v1/payments/payment/
func (c *Client) GetPayments() ([]Payment, error) {
	var p ListPaymentsResp

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payment/"), nil)
	if err != nil {
		return p.Payments, err
	}

	if err = c.SendWithAuth(req, &p); err != nil {
		return p.Payments, err
	}

	return p.Payments, nil
}
