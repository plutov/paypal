package paypalsdk

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// CreateDirectPaypalPayment sends request with payment
func (c *Client) CreateDirectPaypalPayment(amount Amount, redirectURI string) (*PaymentResponse, error) {
	buf := bytes.NewBuffer([]byte("{\"intent\":\"sale\",\"payer\":{\"payment_method\":\"paypal\"}," +
		"\"transactions\":[{\"amount\":{\"total\":\"" + strconv.FormatFloat(amount.Total, 'f', 2, 64) +
		"\",\"currency\":\"" + amount.Currency + "\"}}],\"redirect_urls\":{\"return_url\":\"" + redirectURI + "\"}}"))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payment"), buf)
	if err != nil {
		return &PaymentResponse{}, err
	}

	req.SetBasicAuth(c.ClientID, c.Secret)
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)

	p := PaymentResponse{}
	err = c.Send(req, &p)

	if p.ID == "" {
		return &p, errors.New("Unable to create payment with this access token")
	}

	return &p, err
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
	err = c.Send(req, &e)

	if e.ID == "" {
		return &e, errors.New("Unable to execute payment with paymentID=" + paymentID)
	}

	return &e, err
}
