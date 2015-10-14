package paypalsdk

import (
    "net/http"
)

// NewClient returns new Client struct
func NewClient(clientID string, secret string, APIBase string) (*Client, error) {
    return &Client{
        &http.Client{},
        clientID,
        secret,
        APIBase,
        nil,
    }, nil
}
