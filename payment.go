package paypalsdk

import (
	"bytes"
	"fmt"
	"net/http"
)

// CreateDirectPaypalPayment sends request with payment
func (c *Client) CreateDirectPaypalPayment(payment PaypalPaymentRequest) (*PaymentResponse, error) {
	buf := bytes.NewBuffer([]byte(""))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payment"), buf)
	if err != nil {
		return &PaymentResponse{}, err
	}

	req.SetBasicAuth(c.ClientID, c.Secret)
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	p := PaymentResponse{}
	err = c.Send(req, &p)

	return &p, nil
}
