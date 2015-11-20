package paypalsdk

import (
	"bytes"
	"fmt"
	"net/http"
)

// GetAccessToken returns struct of TokenResponse
func (c *Client) GetAccessToken() (*TokenResponse, error) {
	buf := bytes.NewBuffer([]byte("grant_type=client_credentials"))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/oauth2/token"), buf)
	if err != nil {
		return &TokenResponse{}, err
	}

	req.SetBasicAuth(c.ClientID, c.Secret)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	t := TokenResponse{}
	err = c.Send(req, &t)

	// Set Token fur current Client
	if t.Token != "" {
		c.Token = &t
	}

	return &t, err
}
