package paypalsdk

import (
	"fmt"
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	token, err := c.GetAccessToken()
	if err != nil || token.Token == "" {
		t.Errorf("Token is not returned by GetAccessToken")
	}
}

func TestGetAuthorization(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	a, err := c.GetAuthorization(testAuthID)
	if err != nil || a.ID != testAuthID {
		t.Errorf("GetAuthorization failed for ID=%s", testAuthID)
	}

	a, err = c.GetAuthorization(testFakeAuthID)
	if err == nil {
		t.Errorf("GetAuthorization must return error for ID=%s", testFakeAuthID)
	}
}

func TestCaptureAuthorization(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.CaptureAuthorization(testAuthID, &Amount{Total: "200", Currency: "USD"}, true)

	if err == nil {
		t.Errorf("Auth is expired, 400 error must be returned")
	} else {
		fmt.Println(err.Error())
	}
}

func TestVoidAuthorization(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.VoidAuthorization(testAuthID)

	if err == nil {
		t.Errorf("Auth is expired, 400 error must be returned")
	} else {
		fmt.Println(err.Error())
	}
}

func TestReauthorizeAuthorization(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	_, err := c.ReauthorizeAuthorization(testAuthID, &Amount{Total: "200", Currency: "USD"})

	if err == nil {
		t.Errorf("Reauthorization not allowed for this product, 500 error must be returned")
	} else {
		fmt.Println(err.Error())
	}
}
