package paypal

import (
	"context"
	"fmt"
)

// GenerateInvoiceNumber: generates the next invoice number that is available to the merchant.
// Endpoint: POST /v2/invoicing/generate-next-invoice-number
func (c *Client) GenerateInvoiceNumber(ctx context.Context) (*InvoiceNumber, error) {

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/invoicing/generate-next-invoice-number"), nil)
	nextInvoiceNumber := &InvoiceNumber{}
	if err != nil {
		return nextInvoiceNumber, err
	}

	if err = c.SendWithAuth(req, nextInvoiceNumber); err != nil {
		return nextInvoiceNumber, err
	}

	return nextInvoiceNumber, nil
}

// GetInvoiceDetails: show invoice details for a particular invoice by ID.
// Endpoint: GET /v2/invoicing/invoices/{invoice_id}
func (c *Client) GetInvoiceDetails(ctx context.Context, invoiceID string) (*Invoice, error) {
	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s%s", c.APIBase, "/v2/invoicing/invoices/", invoiceID), nil)
	invoice := &Invoice{}
	if err != nil {
		return invoice, err
	}

	if err = c.SendWithAuth(req, invoice); err != nil {
		return invoice, err
	}
	return invoice, nil
}
