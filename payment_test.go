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

	_, err := c.CreateDirectPaypalPayment(amount, "http://example.com")

	if err == nil {
		t.Errorf("Error must be returned for invalid token")
	}
}

func TestExecuteApprovedPayment(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.Token = &TokenResponse{
		Token: "invalidtoken",
	}

	_, err := c.ExecuteApprovedPayment("PAY-6RV70583SB702805EKEYSZ6Y", "7E7MGXCWTTKK2")

	if err == nil {
		t.Errorf("Error must be returned for invalid token")
	}
}
