package paypalsdk

import "testing"

func TestGetSale(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.GetSale("1")
	if err == nil {
		t.Errorf("GetSale must be failed")
	}
}

func TestRefundSale(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.RefundSale("1", nil)
	if err == nil {
		t.Errorf("RefundSale must be failed")
	}

	_, err = c.RefundSale("1", &Amount{Total: "100", Currency: "USD"})
	if err == nil {
		t.Errorf("RefundSale must be failed")
	}
}

func TestGetRefund(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.GetRefund("1")
	if err == nil {
		t.Errorf("GetRefund must be failed")
	}
}
