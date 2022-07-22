package paypal

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetOrder retrieves order by ID
// Endpoint: GET /v2/checkout/orders/ID
func (c *Client) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	order := &Order{}

	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s%s", c.APIBase, "/v2/checkout/orders/", orderID), nil)
	if err != nil {
		return order, err
	}

	if err = c.SendWithAuth(req, order); err != nil {
		return order, err
	}

	return order, nil
}

// CreateOrder - Use this call to create an order
// Endpoint: POST /v2/checkout/orders
func (c *Client) CreateOrder(ctx context.Context, intent string, purchaseUnits []PurchaseUnitRequest, payer *CreateOrderPayer, appContext *ApplicationContext) (*Order, error) {
	return c.CreateOrderWithPaypalRequestID(ctx, intent, purchaseUnits, payer, appContext, "")
}

// CreateOrderWithPaypalRequestID - Use this call to create an order with idempotency
// Endpoint: POST /v2/checkout/orders
func (c *Client) CreateOrderWithPaypalRequestID(ctx context.Context,
	intent string,
	purchaseUnits []PurchaseUnitRequest,
	payer *CreateOrderPayer,
	appContext *ApplicationContext,
	requestID string,
) (*Order, error) {
	type createOrderRequest struct {
		Intent             string                `json:"intent"`
		Payer              *CreateOrderPayer     `json:"payer,omitempty"`
		PurchaseUnits      []PurchaseUnitRequest `json:"purchase_units"`
		ApplicationContext *ApplicationContext   `json:"application_context,omitempty"`
	}

	order := &Order{}

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/checkout/orders"), createOrderRequest{Intent: intent, PurchaseUnits: purchaseUnits, Payer: payer, ApplicationContext: appContext})
	if err != nil {
		return order, err
	}

	if requestID != "" {
		req.Header.Set("PayPal-Request-Id", requestID)
	}

	if err = c.SendWithAuth(req, order); err != nil {
		return order, err
	}

	return order, nil
}

// UpdateOrder updates the order by ID
// Endpoint: PATCH /v2/checkout/orders/ID
func (c *Client) UpdateOrder(ctx context.Context, orderID string, op string, path string, value map[string]string) error {

	type patchRequest struct {
		Op    string            `json:"op"`
		Path  string            `json:"path"`
		Value map[string]string `json:"value"`
	}

	req, err := c.NewRequest(ctx, "PATCH", fmt.Sprintf("%s%s%s", c.APIBase, "/v2/checkout/orders/", orderID), []patchRequest{
		{
			Op:    op,
			Path:  path,
			Value: value,
		},
	})
	if err != nil {
		return err
	}

	if err = c.SendWithAuth(req, nil); err != nil {
		return err
	}
	return nil
}

// AuthorizeOrder - https://developer.paypal.com/docs/api/orders/v2/#orders_authorize
// Endpoint: POST /v2/checkout/orders/ID/authorize
func (c *Client) AuthorizeOrder(ctx context.Context, orderID string, authorizeOrderRequest AuthorizeOrderRequest) (*AuthorizeOrderResponse, error) {
	auth := &AuthorizeOrderResponse{}

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/checkout/orders/"+orderID+"/authorize"), authorizeOrderRequest)
	if err != nil {
		return auth, err
	}

	if err = c.SendWithAuth(req, auth); err != nil {
		return auth, err
	}

	return auth, nil
}

// CaptureOrder - https://developer.paypal.com/docs/api/orders/v2/#orders_capture
// Endpoint: POST /v2/checkout/orders/ID/capture
func (c *Client) CaptureOrder(ctx context.Context, orderID string, captureOrderRequest CaptureOrderRequest) (*CaptureOrderResponse, error) {
	return c.CaptureOrderWithPaypalRequestId(ctx, orderID, captureOrderRequest, "", nil)
}

// CaptureOrder with idempotency - https://developer.paypal.com/docs/api/orders/v2/#orders_capture
// Endpoint: POST /v2/checkout/orders/ID/capture
// https://developer.paypal.com/docs/api/reference/api-requests/#http-request-headers
func (c *Client) CaptureOrderWithPaypalRequestId(ctx context.Context,
	orderID string,
	captureOrderRequest CaptureOrderRequest,
	requestID string,
	mockResponse *CaptureOrderMockResponse,
) (*CaptureOrderResponse, error) {
	capture := &CaptureOrderResponse{}

	c.SetReturnRepresentation()
	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/checkout/orders/"+orderID+"/capture"), captureOrderRequest)
	if err != nil {
		return capture, err
	}

	if requestID != "" {
		req.Header.Set("PayPal-Request-Id", requestID)
	}

	if mockResponse != nil {
		mock, err := json.Marshal(mockResponse)
		if err != nil {
			return nil, err
		}

		req.Header.Set("PayPal-Mock-Response", string(mock))
	}

	if err = c.SendWithAuth(req, capture); err != nil {
		return capture, err
	}

	return capture, nil
}

// RefundCapture - https://developer.paypal.com/docs/api/payments/v2/#captures_refund
// Endpoint: POST /v2/payments/captures/ID/refund
func (c *Client) RefundCapture(ctx context.Context, captureID string, refundCaptureRequest RefundCaptureRequest) (*RefundResponse, error) {
	return c.RefundCaptureWithPaypalRequestId(ctx, captureID, refundCaptureRequest, "")
}

// RefundCapture with idempotency - https://developer.paypal.com/docs/api/payments/v2/#captures_refund
// Endpoint: POST /v2/payments/captures/ID/refund
func (c *Client) RefundCaptureWithPaypalRequestId(ctx context.Context,
	captureID string,
	refundCaptureRequest RefundCaptureRequest,
	requestID string,
) (*RefundResponse, error) {
	refund := &RefundResponse{}

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/captures/"+captureID+"/refund"), refundCaptureRequest)
	if err != nil {
		return refund, err
	}

	if requestID != "" {
		req.Header.Set("PayPal-Request-Id", requestID)
	}

	if err = c.SendWithAuth(req, refund); err != nil {
		return refund, err
	}
	return refund, nil
}

// CapturedDetail - https://developer.paypal.com/docs/api/payments/v2/#captures_get
// Endpoint: GET /v2/payments/captures/ID
func (c *Client) CapturedDetail(ctx context.Context, captureID string) (*CaptureDetailsResponse, error) {
	response := &CaptureDetailsResponse{}

	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/captures/"+captureID), nil)
	if err != nil {
		return response, err
	}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}
	return response, nil
}
