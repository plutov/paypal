package paypalsdk

import (
	"fmt"
	"testing"
)

func TestGetOrder(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.GetOrder("1")

	if err == nil {
		t.Errorf("GetOrder must be failed")
	} else {
		fmt.Println(err.Error())
	}
}

func TestAuthorizeOrder(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.AuthorizeOrder("1", &Amount{Total: "100", Currency: "USD"})
	if err == nil {
		t.Errorf("AuthorizeOrder must be failed")
	} else {
		fmt.Println(err.Error())
	}
}

func TestCaptureOrder(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.CaptureOrder("1", &Amount{Total: "100", Currency: "USD"}, true, nil)
	if err == nil {
		t.Errorf("CaptureOrder must be failed")
	} else {
		fmt.Println(err.Error())
	}
}

func TestVoidOrder(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.VoidOrder("1")
	if err == nil {
		t.Errorf("VoidOrder must be failed")
	} else {
		fmt.Println(err.Error())
	}
}
