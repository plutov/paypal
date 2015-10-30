package paypalsdk

import (
	"testing"
)

func TestGetAuthorizationCodeURL(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)

	_, err := c.GetAuthorizationCodeURL("", []string{})
	if err == nil {
		t.Errorf("redirectURI is required in GetAuthorizationCodeURL")
	}

	uri, err := c.GetAuthorizationCodeURL("test", []string{})
	if uri != "https://api.sandbox.paypal.com/webapps/auth/protocol/openidconnect/v1/authorize?client_id=clid&response_type=code&scope=profile+email&redirect_uri=test" {
		t.Errorf("GetAuthorizationCodeURL returns incorrect value for redirectURI=test")
	}

	uri, err = c.GetAuthorizationCodeURL("test", []string{"address"})
	if uri != "https://api.sandbox.paypal.com/webapps/auth/protocol/openidconnect/v1/authorize?client_id=clid&response_type=code&scope=address&redirect_uri=test" {
		t.Errorf("GetAuthorizationCodeURL returns incorrect value for redirectURI=test and scope=address")
	}
}

func TestGetAccessToken(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)

	_, err := c.GetAccessToken("123", "")
	if err == nil {
		t.Errorf("redirectURI is required in GetRefreshToken")
	}

	_, err = c.GetAccessToken("", "123")
	if err == nil {
		t.Errorf("authorizationCode is required in GetRefreshToken")
	}

	_, err = c.GetAccessToken("123", "123")
}
