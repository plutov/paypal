package paypalsdk

import (
	"fmt"
	"testing"
)

func TestCreateDirectPaypalPayment(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.Token = &TokenResponse{
		Token: "invalidtoken",
	}

	amount := Amount{
		Total:    "15.1111",
		Currency: "USD",
	}

	_, err := c.CreateDirectPaypalPayment(amount, "http://example.com", "http://example.com", "test payment")

	if err == nil {
		t.Errorf("Error must be returned for invalid token")
	} else {
		fmt.Println(err.Error())
	}
}

func TestGetPayment(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.Token = &TokenResponse{
		Token: "invalidtoken",
	}

	_, err := c.GetPayment("PAY-TEST-123")

	if err == nil {
		t.Errorf("Error must be returned for invalid ID")
	} else {
		fmt.Println(err.Error())
	}
}

func TestGetPayments(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.Token = &TokenResponse{
		Token: "invalidtoken",
	}

	payments, _ := c.GetPayments()

	if len(payments) != 0 {
		t.Errorf("0 payments must be returned for unautgorized request")
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
	} else {
		fmt.Println(err.Error())
	}
}
