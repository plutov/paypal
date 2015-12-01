package paypalsdk

import (
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.GetAccessToken()
}

func TestGetAuthorization(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.GetAuthorization("123")

	if err == nil {
		t.Errorf("Error must be returned for invalid Auth ID")
	}
}

func TestCaptureAuthorization(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.CaptureAuthorization("123", &Amount{Total: "200", Currency: "USD"}, true)

	if err == nil {
		t.Errorf("Error must be returned for invalid Auth ID")
	}
}

func TestVoidAuthorization(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.VoidAuthorization("123")

	if err == nil {
		t.Errorf("Error must be returned for invalid Auth ID")
	}
}

func TestReauthorizeAuthorization(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.ReauthorizeAuthorization("123", &Amount{Total: "200", Currency: "USD"})

	if err == nil {
		t.Errorf("Error must be returned for invalid Auth ID")
	}
}
