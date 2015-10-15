package paypalsdk

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "bytes"
    "time"
    "fmt"
    "io"
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

// GetAcessToken request a new access token from Paypal
func (c *Client) GetAccessToken() (*TokenResponse, error) {
    buf := bytes.NewBuffer([]byte("grant_type=client_credentials"))
    req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/oauth2/token"), buf)
    if err != nil {
        return nil, err
    }

    req.SetBasicAuth(c.ClientID, c.Secret)
    req.Header.Set("Content-type", "application/x-www-form-urlencoded")

    t := TokenResponse{}
    err = c.Send(req, &t)
    if err == nil {
        t.ExpiresAt = time.Now().Add(time.Duration(t.ExpiresIn/2) * time.Second)
    }

    return &t, err
}

// Send makes a request to the API, the response body will be
// unmarshaled into v, or if v is an io.Writer, the response will
// be written to it without decoding
func (c *Client) Send(req *http.Request, v interface{}) error {
    // Set default headers
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Accept-Language", "en_US")

    // Default values for headers
    if req.Header.Get("Content-type") == "" {
        req.Header.Set("Content-type", "application/json")
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode > 299 {
        errResp := &ErrorResponse{Response: resp}
        data, err := ioutil.ReadAll(resp.Body)

        if err == nil && len(data) > 0 {
            json.Unmarshal(data, errResp)
        }

        return errResp
    }

    if v != nil {
        if w, ok := v.(io.Writer); ok {
            io.Copy(w, resp.Body)
        } else {
            err = json.NewDecoder(resp.Body).Decode(v)
            if err != nil {
                return err
            }
        }
    }

    return nil
}
