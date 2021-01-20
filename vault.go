package paypal

import (
	"context"
	"fmt"
)

// StoreCreditCard func
// Endpoint: POST /v1/vault/credit-cards
func (c *Client) StoreCreditCard(ctx context.Context, cc CreditCard) (*CreditCard, error) {
	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/vault/credit-cards"), cc)
	if err != nil {
		return nil, err
	}

	response := &CreditCard{}

	if err = c.SendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteCreditCard func
// Endpoint: DELETE /v1/vault/credit-cards/credit_card_id
func (c *Client) DeleteCreditCard(ctx context.Context, id string) error {
	req, err := c.NewRequest(ctx, "DELETE", fmt.Sprintf("%s/v1/vault/credit-cards/%s", c.APIBase, id), nil)
	if err != nil {
		return err
	}

	if err = c.SendWithAuth(req, nil); err != nil {
		return err
	}

	return nil
}

// GetCreditCard func
// Endpoint: GET /v1/vault/credit-cards/credit_card_id
func (c *Client) GetCreditCard(ctx context.Context, id string) (*CreditCard, error) {
	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s/v1/vault/credit-cards/%s", c.APIBase, id), nil)
	if err != nil {
		return nil, err
	}

	response := &CreditCard{}

	if err = c.SendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetCreditCards func
// Endpoint: GET /v1/vault/credit-cards
func (c *Client) GetCreditCards(ctx context.Context, ccf *CreditCardsFilter) (*CreditCards, error) {
	page := 1
	if ccf != nil && ccf.Page > 0 {
		page = ccf.Page
	}
	pageSize := 10
	if ccf != nil && ccf.PageSize > 0 {
		pageSize = ccf.PageSize
	}

	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s/v1/vault/credit-cards?page=%d&page_size=%d", c.APIBase, page, pageSize), nil)
	if err != nil {
		return nil, err
	}

	response := &CreditCards{}

	if err = c.SendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

// PatchCreditCard func
// Endpoint: PATCH /v1/vault/credit-cards/credit_card_id
func (c *Client) PatchCreditCard(ctx context.Context, id string, ccf []CreditCardField) (*CreditCard, error) {
	req, err := c.NewRequest(ctx, "PATCH", fmt.Sprintf("%s/v1/vault/credit-cards/%s", c.APIBase, id), ccf)
	if err != nil {
		return nil, err
	}

	response := &CreditCard{}

	if err = c.SendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}
