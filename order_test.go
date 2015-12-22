package paypalsdk

import (
	"fmt"
	"testing"
)

func TestGetOrder(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	o, err := c.GetOrder(testOrderID)

	if err != nil || o.ID != testOrderID {
		t.Errorf("GetOrder failed for ID=" + testOrderID)
	}
}

func TestAuthorizeOrder(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.AuthorizeOrder(testOrderID, &Amount{Total: "7.00", Currency: "USD"})
	if err == nil {
		t.Errorf("Order is expired, 400 error must be returned")
	} else {
		fmt.Println(err.Error())
	}
}

func TestCaptureOrder(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.CaptureOrder(testOrderID, &Amount{Total: "100", Currency: "USD"}, true, nil)
	if err == nil {
		t.Errorf("Order is expired, 400 error must be returned")
	} else {
		fmt.Println(err.Error())
	}
}

func TestVoidOrder(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.VoidOrder(testOrderID)
	if err == nil {
		t.Errorf("Order is expired, 400 error must be returned")
	} else {
		fmt.Println(err.Error())
	}
}
