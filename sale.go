package paypalsdk

import "fmt"

// GetSale returns a sale by ID
func (c *Client) GetSale(saleID string) (*Sale, error) {
	sale := &Sale{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/sale/"+saleID), nil)
	if err != nil {
		return sale, err
	}

	err = c.SendWithAuth(req, sale)
	if err != nil {
		return sale, err
	}

	return sale, nil
}

// RefundSale refunds a completed payment.
// Amount can be sent to make a partial refund only
func (c *Client) RefundSale(saleID string, a *Amount) (*Refund, error) {
	type refundRequest struct {
		Amount *Amount `json:"amount"`
	}

	refund := &Refund{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/sale/"+saleID+"/refund"), &refundRequest{Amount: a})
	if err != nil {
		return refund, err
	}

	err = c.SendWithAuth(req, refund)
	if err != nil {
		return refund, err
	}

	return refund, nil
}
