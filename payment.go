package paypalsdk

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// CreateDirectPaypalPayment sends request with payment
func (c *Client) CreateDirectPaypalPayment(amount Amount) (*PaymentResponse, error) {
	buf := bytes.NewBuffer([]byte("{\"intent\":\"sale\",\"payer\":{\"payment_method\":\"paypal\"},\"transactions\":[{\"amount\":{\"total\":\"" + strconv.FormatFloat(amount.Total, 'f', 2, 64) + "\",\"currency\":\"" + amount.Currency + "\"}}]}"))
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
