package paypalsdk

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GrantNewAccessTokenFromAuthCode - Use this call to grant a new access token, using the previously obtained authorization code.
// Endpoint: POST /v1/identity/openidconnect/tokenservice
func (c *Client) GrantNewAccessTokenFromAuthCode(code string, redirectURI string) (*TokenResponse, error) {
	token := &TokenResponse{}

	q := url.Values{}
	q.Set("grant_type", "authorization_code")
	q.Set("code", code)
	q.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/identity/openidconnect/tokenservice"), strings.NewReader(q.Encode()))
	if err != nil {
		return token, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err = c.SendWithBasicAuth(req, token); err != nil {
		return token, err
	}

	return token, nil
}

// GrantNewAccessTokenFromRefreshToken - Use this call to grant a new access token, using a refresh token.
// Endpoint: POST /v1/identity/openidconnect/tokenservice
func (c *Client) GrantNewAccessTokenFromRefreshToken(refreshToken string) (*TokenResponse, error) {
	type request struct {
		GrantType    string `json:"grant_type"`
		RefreshToken string `json:"refresh_token"`
	}

	token := &TokenResponse{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/identity/openidconnect/tokenservice"), request{GrantType: "refresh_token", RefreshToken: refreshToken})
	if err != nil {
		return token, err
	}

	if err = c.SendWithAuth(req, token); err != nil {
		return token, err
	}

	return token, nil
}

// GetUserInfo - Use this call to retrieve user profile attributes.
// Endpoint: GET /v1/identity/openidconnect/userinfo/?schema=<Schema>
// Pass the schema that is used to return as per openidconnect protocol. The only supported schema value is openid.
func (c *Client) GetUserInfo(schema string) (*UserInfo, error) {
	u := UserInfo{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", c.APIBase, "/v1/identity/openidconnect/userinfo/?schema=", schema), nil)
	if err != nil {
		return &u, err
	}

	if err = c.SendWithAuth(req, &u); err != nil {
		return &u, err
	}

	return &u, nil
}
