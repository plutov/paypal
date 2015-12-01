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

// CreatePaymentResp returned by CreatePayment
type CreatePaymentResp struct {
	*Payment
	Links []Links `json:"links"`
}

// CreateDirectPaypalPayment sends request with payment
// CreatePayment is more common function for any kind of payment
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
	err = c.SendWithAuth(req, &p)
	if err != nil {
		return &p, err
	}

	if p.ID == "" {
		return &p, errors.New("Unable to create payment with this access token")
	}

	return &p, err
}

// CreatePayment creates a payment in Paypal
func (c *Client) CreatePayment(p Payment) (*CreatePaymentResp, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payment"), p)
	if err != nil {
		return &CreatePaymentResp{}, err
	}

	response := &CreatePaymentResp{}

	err = c.SendWithAuth(req, response)
	if err != nil {
		return response, err
	}

	return response, nil
}

// ExecuteApprovedPayment executes approved payment
func (c *Client) ExecuteApprovedPayment(paymentID string, payerID string) (*ExecuteResponse, error) {
	buf := bytes.NewBuffer([]byte("{\"payer_id\":\"" + payerID + "\"}"))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payment/"+paymentID+"/execute"), buf)
	if err != nil {
		return &ExecuteResponse{}, err
	}

	req.SetBasicAuth(c.ClientID, c.Secret)
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)

	e := ExecuteResponse{}
	err = c.SendWithAuth(req, &e)
	if err != nil {
		return &e, err
	}

	if e.ID == "" {
		return &e, errors.New("Unable to execute payment with paymentID=" + paymentID)
	}

	return &e, err
}

// GetPayment gets a payment from PayPal
func (c *Client) GetPayment(paymentID string) (*Payment, error) {
	p := Payment{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payment/"+paymentID), nil)
	if err != nil {
		return &p, err
	}

	err = c.SendWithAuth(req, &p)
	if err != nil {
		return &p, err
	}

	if p.ID == "" {
		return &p, errors.New("Unable to get payment with paymentID=" + paymentID)
	}

	return &p, nil
}

// GetPayments retrieve payments resources from Paypal
func (c *Client) GetPayments() ([]Payment, error) {
	var p ListPaymentsResp

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payment/"), nil)
	if err != nil {
		return p.Payments, err
	}

	err = c.SendWithAuth(req, &p)
	if err != nil {
		return p.Payments, err
	}

	return p.Payments, nil
}
