package paypalsdk

import (
	"fmt"
	"net/http"
)

// GrantNewAccessTokenFromAuthCode - Use this call to grant a new access token, using the previously obtained authorization code.
// Endpoint: POST /v1/identity/openidconnect/tokenservice
func (c *Client) GrantNewAccessTokenFromAuthCode(code string, redirectURI string) (*TokenResponse, error) {
	type request struct {
		GrantType   string `json:"grant_type"`
		Code        string `json:"code"`
		RedirectURI string `json:"redirect_uri"`
	}

	token := &TokenResponse{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/identity/openidconnect/tokenservice"), request{GrantType: "authorization_code", Code: code, RedirectURI: redirectURI})
	if err != nil {
		return token, err
	}

	err = c.Send(req, token)
	if err != nil {
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

	err = c.Send(req, token)
	if err != nil {
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

	err = c.SendWithAuth(req, &u)
	if err != nil {
		return &u, err
	}

	return &u, nil
}
