package paypal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// CreateWebhook - Subscribes your webhook listener to events.
// Endpoint: POST /v1/notifications/webhooks
func (c *Client) CreateWebhook(createWebhookRequest *CreateWebhookRequest) (*CreateWebhookResponse, error) {
	req, err := c.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/webhooks"), createWebhookRequest)
	res := &CreateWebhookResponse{}
	if err != nil {
		return res, err
	}

	err = c.SendWithAuth(req, res)
	return res, err
}

// ListWebhooks - Lists webhooks for an app.
// Endpoint: GET /v1/notifications/webhooks
func (c *Client) ListWebhooks(anchorType string) (*ListWebhookResponse, error) {
	if len(anchorType) == 0 {
		anchorType = AncorTypeApplication
	}
	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/webhooks"), nil)
	q := req.URL.Query()
	q.Add("anchor_type", anchorType)
	req.URL.RawQuery = q.Encode()
	resp := &ListWebhookResponse{}
	if err != nil {
		return resp, err
	}

	err = c.SendWithAuth(req, resp)
	return resp, err
}

// VerifyWebhookSignature - Use this to verify the signature of a webhook recieved from paypal.
// Endpoint: POST /v1/notifications/verify-webhook-signature
func (c *Client) VerifyWebhookSignature(httpReq *http.Request, webhookID string) (*VerifyWebhookResponse, error) {
	type verifyWebhookSignatureRequest struct {
		AuthAlgo         string          `json:"auth_algo,omitempty"`
		CertURL          string          `json:"cert_url,omitempty"`
		TransmissionID   string          `json:"transmission_id,omitempty"`
		TransmissionSig  string          `json:"transmission_sig,omitempty"`
		TransmissionTime string          `json:"transmission_time,omitempty"`
		WebhookID        string          `json:"webhook_id,omitempty"`
		WebhookEvent     json.RawMessage `json:"webhook_event"`
	}

	// Read the content
	var bodyBytes []byte
	if httpReq.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(httpReq.Body)
	}
	// Restore the io.ReadCloser to its original state
	httpReq.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	verifyRequest := verifyWebhookSignatureRequest{
		AuthAlgo:         httpReq.Header.Get("PAYPAL-AUTH-ALGO"),
		CertURL:          httpReq.Header.Get("PAYPAL-CERT-URL"),
		TransmissionID:   httpReq.Header.Get("PAYPAL-TRANSMISSION-ID"),
		TransmissionSig:  httpReq.Header.Get("PAYPAL-TRANSMISSION-SIG"),
		TransmissionTime: httpReq.Header.Get("PAYPAL-TRANSMISSION-TIME"),
		WebhookID:        webhookID,
		WebhookEvent:     json.RawMessage(bodyBytes),
	}

	response := &VerifyWebhookResponse{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/verify-webhook-signature"), verifyRequest)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}
