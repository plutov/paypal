package paypalsdk

import (
	"fmt"
	"testing"
)

func TestGetSale(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.GetSale(testSaleID)
	if err == nil {
		t.Errorf("404 must be returned for ID=%s", testSaleID)
	} else {
		fmt.Println(err.Error())
	}
}

func TestRefundSale(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.RefundSale(testSaleID, nil)
	if err == nil {
		t.Errorf("404 must be returned for ID=%s", testSaleID)
	} else {
		fmt.Println(err.Error())
	}

	_, err = c.RefundSale(testSaleID, &Amount{Total: "7.00", Currency: "USD"})
	if err == nil {
		t.Errorf("404 must be returned for ID=%s", testSaleID)
	} else {
		fmt.Println(err.Error())
	}
}

func TestGetRefund(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.GetRefund("1")
	if err == nil {
		t.Errorf("404 must be returned for ID=%s", testSaleID)
	} else {
		fmt.Println(err.Error())
	}
}
