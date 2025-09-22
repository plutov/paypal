package paypal

import (
	"context"
	"fmt"
)

// SetTransactionContext
// PayPal partners and merchants can Set transaction context to send additional data about a customer to PayPal before the customer transaction is processed.
// PayPal uses this data to complete a pre-transaction risk management evaluation
// Endpoint: PUT <endpoint>/v1/risk/transaction-contexts/<merchant_id>/<tracking_id>
func (c *Client) SetTransactionContext(ctx context.Context, merchantID, orderID string, tcPayload interface{}) error {

	if merchantID == "" || orderID == "" {
		return fmt.Errorf("paypal: merchantID or orderID is empty")
	}

	url := fmt.Sprintf("%s%s%s", c.APIBase, "/v1/risk/transaction-contexts/", merchantID+"/"+orderID)

	req, err := c.NewRequest(ctx, "PUT", url, tcPayload)

	if err != nil {
		return err
	}

	if err = c.SendWithAuth(req, nil); err != nil {
		return err
	}

	return nil
}
