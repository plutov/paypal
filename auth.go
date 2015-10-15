package paypalsdk

import (
    "strings"
    "net/url"
    "errors"
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
