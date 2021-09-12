package paypal

import (
	"context"
	"fmt"
	"net/http"
)

// CreateWebProfile creates a new web experience profile in Paypal
//
// Allows for the customisation of the payment experience
//
// Endpoint: POST /v1/payment-experience/web-profiles
func (c *Client) CreateWebProfile(ctx context.Context, wp WebProfile) (*WebProfile, error) {
	url := fmt.Sprintf("%s%s", c.APIBase, "/v1/payment-experience/web-profiles")
	req, err := c.NewRequest(ctx, "POST", url, wp)
	response := &WebProfile{}

	if err != nil {
		return response, err
	}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	return response, nil
}

// GetWebProfile gets an exists payment experience from Paypal
//
// Endpoint: GET /v1/payment-experience/web-profiles/<profile-id>
func (c *Client) GetWebProfile(ctx context.Context, profileID string) (*WebProfile, error) {
	var wp WebProfile

	url := fmt.Sprintf("%s%s%s", c.APIBase, "/v1/payment-experience/web-profiles/", profileID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return &wp, err
	}

	if err = c.SendWithAuth(req, &wp); err != nil {
		return &wp, err
	}

	if wp.ID == "" {
		return &wp, fmt.Errorf("paypal: unable to get web profile with ID = %s", profileID)
	}

	return &wp, nil
}

// GetWebProfiles retrieves web experience profiles from Paypal
//
// Endpoint: GET /v1/payment-experience/web-profiles
func (c *Client) GetWebProfiles(ctx context.Context) ([]WebProfile, error) {
	var wps []WebProfile

	url := fmt.Sprintf("%s%s", c.APIBase, "/v1/payment-experience/web-profiles")
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return wps, err
	}

	if err = c.SendWithAuth(req, &wps); err != nil {
		return wps, err
	}

	return wps, nil
}

// SetWebProfile sets a web experience profile in Paypal with given id
//
// Endpoint: PUT /v1/payment-experience/web-profiles
func (c *Client) SetWebProfile(ctx context.Context, wp WebProfile) error {

	if wp.ID == "" {
		return fmt.Errorf("paypal: no ID specified for WebProfile")
	}

	url := fmt.Sprintf("%s%s%s", c.APIBase, "/v1/payment-experience/web-profiles/", wp.ID)

	req, err := c.NewRequest(ctx, "PUT", url, wp)

	if err != nil {
		return err
	}

	if err = c.SendWithAuth(req, nil); err != nil {
		return err
	}

	return nil
}

// DeleteWebProfile deletes a web experience profile from Paypal with given id
//
// Endpoint: DELETE /v1/payment-experience/web-profiles
func (c *Client) DeleteWebProfile(ctx context.Context, profileID string) error {

	url := fmt.Sprintf("%s%s%s", c.APIBase, "/v1/payment-experience/web-profiles/", profileID)

	req, err := c.NewRequest(ctx, "DELETE", url, nil)

	if err != nil {
		return err
	}

	if err = c.SendWithAuth(req, nil); err != nil {
		return err
	}

	return nil
}
