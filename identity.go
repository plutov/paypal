package paypal

import (
	"context"
	"fmt"
	"net/http"
)

// GetUserInfo - Use this call to retrieve user profile attributes.
// Endpoint: GET /v1/identity/openidconnect/userinfo/?schema=<Schema>
func (c *Client) GetUserInfo(ctx context.Context, schema string) (*UserInfo, error) {
	u := &UserInfo{}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s%s", c.APIBase, "/v1/identity/openidconnect/userinfo/?schema=", schema), nil)
	if err != nil {
		return u, err
	}

	if err = c.SendWithAuth(req, u); err != nil {
		return u, err
	}

	return u, nil
}
