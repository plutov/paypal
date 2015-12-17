package paypalsdk

import "fmt"

// GetOrder retreives order
func (c *Client) GetOrder(orderID string) (*Order, error) {
	order := &Order{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/orders/"+orderID), nil)
	if err != nil {
		return order, err
	}

	err = c.SendWithAuth(req, order)
	if err != nil {
		return order, err
	}

	return order, nil
}

// AuthorizeOrder POST /v1/payments/orders/<Order-Id>/authorize
func (c *Client) AuthorizeOrder(orderID string, amount *Amount) (*Authorization, error) {
	type authRequest struct {
		Amount *Amount `json:"amount"`
	}

	auth := &Authorization{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/orders/"+orderID+"/authorize"), authRequest{Amount: amount})
	if err != nil {
		return auth, err
	}

	err = c.SendWithAuth(req, auth)
	if err != nil {
		return auth, err
	}

	return auth, nil
}

// CaptureOrder POST /v1/payments/orders/<Order-Id>/capture
func (c *Client) CaptureOrder(orderID string, amount *Amount, isFinalCapture bool, currency *Currency) (*Capture, error) {
	type captureRequest struct {
		Amount         *Amount   `json:"amount"`
		IsFinalCapture bool      `json:"is_final_capture"`
		Currency       *Currency `json:"transaction_fee"`
	}

	capture := &Capture{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/orders/"+orderID+"/capture"), captureRequest{Amount: amount, IsFinalCapture: isFinalCapture, Currency: currency})
	if err != nil {
		return capture, err
	}

	err = c.SendWithAuth(req, capture)
	if err != nil {
		return capture, err
	}

	return capture, nil
}

// VoidOrder POST /v1/payments/orders/<Order-Id>/do-void
func (c *Client) VoidOrder(orderID string) (*Order, error) {
	order := &Order{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/orders/"+orderID+"/do-void"), nil)
	if err != nil {
		return order, err
	}

	err = c.SendWithAuth(req, order)
	if err != nil {
		return order, err
	}

	return order, nil
}
