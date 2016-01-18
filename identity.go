package paypalsdk

import "fmt"

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
