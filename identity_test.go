package paypalsdk

import "testing"

func TestGrantNewAccessTokenFromAuthCode(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)

	_, err := c.GrantNewAccessTokenFromAuthCode("123", "http://example.com/myapp/return.php")
	if err == nil {
		t.Errorf("GrantNewAccessTokenFromAuthCode must return error for invalid code")
	}
}

func TestGrantNewAccessTokenFromRefreshToken(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)

	_, err := c.GrantNewAccessTokenFromRefreshToken("123")
	if err == nil {
		t.Errorf("GrantNewAccessTokenFromRefreshToken must return error for invalid refresh token")
	}
}

func TestGetUserInfo(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	u, err := c.GetUserInfo("openid")
	if u.ID != testUserID || err != nil {
		t.Errorf("GetUserInfo must return valid test ID=%s", testUserID)
	}
}
