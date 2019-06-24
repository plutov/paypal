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

// Create Order - Use this call to create an order
// Endpoint: POST /v2/checkout/orders
func (c *Client) CreateOrder(intent string, purchaseUnits []PurchaseUnitRequest, payer *PayerInfo, appContext *ApplicationContext) (*Order, error) {
	type createOrderRequest struct {
		Intent             string                `json:"intent"`
		Payer              *PayerInfo            `json:"payer,omitempty"`
		PurchaseUnits      []PurchaseUnitRequest `json:"purchase_units"`
		ApplicationContext *ApplicationContext   `json:"application_context,omitempty"`
	}

	order := &Order{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/checkout/orders"), createOrderRequest{Intent: intent, PurchaseUnits: purchaseUnits, Payer: payer, ApplicationContext: appContext})

	if err = c.SendWithAuth(req, order); err != nil {
		return order, err
	}

	return order, nil
}

// AuthorizeOrder - Use this call to authorize an order.
// Endpoint: POST /v2/checkout/orders/ID/authorize
func (c *Client) AuthorizeOrder(orderID string, paymentSource PaymentSource) (*Authorization, error) {
	auth := &Authorization{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/checkout/orders/"+orderID+"/authorize"), paymentSource)
	if err != nil {
		return auth, err
	}

	if err = c.SendWithAuth(req, auth); err != nil {
		return auth, err
	}

	return auth, nil
}

// CaptureOrder - Use this call to capture a payment on an order.
// Endpoint: POST /v2/checkout/orders/ID/capture
func (c *Client) CaptureOrder(orderID string, paymentSource PaymentSource) (*Capture, error) {
	capture := &Capture{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/checkout/orders/"+orderID+"/capture"), paymentSource)
	if err != nil {
		return capture, err
	}

	if err = c.SendWithAuth(req, capture); err != nil {
		return capture, err
	}

	return capture, nil
}
