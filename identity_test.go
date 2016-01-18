package paypalsdk

import "testing"

func TestGrantNewAccessTokenFromAuthCode(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)

	token, _ := c.GrantNewAccessTokenFromAuthCode("123", "http://example.com/myapp/return.php")
	if token.Token != "" {
		t.Errorf("Empty token must be returned for invalid code")
	}
}

func TestGrantNewAccessTokenFromRefreshToken(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)

	token, _ := c.GrantNewAccessTokenFromRefreshToken("123")
	if token.Token != "" {
		t.Errorf("Empty token must be returned for invalid refresh token")
	}
}
