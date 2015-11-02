package paypalsdk

import (
	"testing"
)

func TestCreateDirectPaypalPayment(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.Token = &TokenResponse{
		Token: "invalidtoken",
	}

	amount := Amount{
		Total:    15.1111,
		Currency: "USD",
	}

	_, err := c.CreateDirectPaypalPayment(amount)

	if err == nil {
		t.Errorf("Error must be returned for invalid token")
	}
}
