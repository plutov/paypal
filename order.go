package paypalsdk

import "fmt"

// GetOrder retrieves order by ID
// Endpoint: GET /v2/checkout/orders/ID
func (c *Client) GetOrder(orderID string) (*Order, error) {
	order := &Order{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s%s", c.APIBase, "/v2/checkout/orders/", orderID), nil)
	if err != nil {
		return order, err
	}

	if err = c.SendWithAuth(req, order); err != nil {
		return order, err
	}

	return order, nil
}

// AuthorizeOrder - Use this call to authorize an order.
// Endpoint: POST /v2/checkout/orders/ID/authorize
func (c *Client) AuthorizeOrder(orderID string, amount *Amount) (*Authorization, error) {
	type authRequest struct {
		Amount *Amount `json:"amount"`
	}

	auth := &Authorization{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/checkout/orders/"+orderID+"/authorize"), authRequest{Amount: amount})
	if err != nil {
		return auth, err
	}

	if err = c.SendWithAuth(req, auth); err != nil {
		return auth, err
	}

	return auth, nil
}

// CaptureOrder - Use this call to capture a payment on an order. To use this call, an original payment call must specify an intent of order.
// Endpoint: POST /v2/checkout/orders/ID/capture
func (c *Client) CaptureOrder(orderID string, amount *Amount, isFinalCapture bool, currency *Currency) (*Capture, error) {
	type captureRequest struct {
		Amount         *Amount   `json:"amount"`
		IsFinalCapture bool      `json:"is_final_capture"`
		Currency       *Currency `json:"transaction_fee"`
	}

	capture := &Capture{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/checkout/orders/"+orderID+"/capture"), captureRequest{Amount: amount, IsFinalCapture: isFinalCapture, Currency: currency})
	if err != nil {
		return capture, err
	}

	if err = c.SendWithAuth(req, capture); err != nil {
		return capture, err
	}

	return capture, nil
}
