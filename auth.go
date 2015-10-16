package paypalsdk

import (
    "net/http"
    "strings"
    "net/url"
    "errors"
    "bytes"
    "fmt"
)

// GetAuthorizationCodeURL returns URL where we need to redirect user
// After signin in PayPal get authorization_code on redirectURI
func (c *Client) GetAuthorizationCodeURL(redirectURI string, scopes []string) (string, error) {
    if redirectURI == "" {
        return "", errors.New("redirectURI cannot be empty")
    }

    if len(scopes) == 0 {
        scopes = []string{"profile", "email"}
    }

    return c.APIBase + "/webapps/auth/protocol/openidconnect/v1/authorize?client_id=" +
        url.QueryEscape(c.ClientID) + "&response_type=code&scope=" + strings.Join(scopes, "+") +
        "&redirect_uri=" + url.QueryEscape(redirectURI), nil
}

// GetRefreshToken returns struct of RefreshTokenResponse
// Client must to get an authorization code before
// redirectURI must match with redirectURI sent to GetAuthorizationCodeURL
func (c *Client) GetRefreshToken(authorizationCode string, redirectURI string) (*RefreshTokenResponse, error) {
    if authorizationCode == "" {
        return &RefreshTokenResponse{}, errors.New("authorizationCode cannot be empty")
    }
    if redirectURI == "" {
        return &RefreshTokenResponse{}, errors.New("redirectURI cannot be empty")
    }

    buf := bytes.NewBuffer([]byte("grant_type=authorization_code&code=" + url.QueryEscape(authorizationCode) + "&redirect_uri=" + url.QueryEscape(redirectURI)))
    req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/identity/openidconnect/tokenservice"), buf)
    if err != nil {
        return &RefreshTokenResponse{}, err
    }

    req.SetBasicAuth(c.ClientID, c.Secret)
    req.Header.Set("Content-type", "application/x-www-form-urlencoded")

    r := RefreshTokenResponse{}
    err = c.Send(req, &r)

    return &r, err
}
