package paypalsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

// NewClient returns new Client struct
// APIBase is a base API URL, for testing you can use paypalsdk.APIBaseSandBox
func NewClient(clientID string, secret string, APIBase string) (*Client, error) {
	if clientID == "" || secret == "" || APIBase == "" {
		return &Client{}, errors.New("ClientID, Secret and APIBase are required to create a Client")
	}

	return &Client{
		&http.Client{},
		clientID,
		secret,
		APIBase,
		nil,
		nil,
	}, nil
}

// SetLog will set/change the output destination.
// If log file is set paypalsdk will log all requests and responses to this Writer
func (c *Client) SetLog(log io.Writer) error {
	c.Log = log
	return nil
}

// SetAccessToken sets saved token to current client
func (c *Client) SetAccessToken(token string) error {
	c.Token = &TokenResponse{
		Token: token,
	}

	return nil
}

// Send makes a request to the API, the response body will be
// unmarshaled into v, or if v is an io.Writer, the response will
// be written to it without decoding
func (c *Client) Send(req *http.Request, v interface{}) error {
	var (
		err  error
		resp *http.Response
		data []byte
	)

	// Set default headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en_US")

	// Default values for headers
	if req.Header.Get("Content-type") == "" {
		req.Header.Set("Content-type", "application/json")
	}

	resp, err = c.client.Do(req)
	c.log(req, resp)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errResp := &ErrorResponse{Response: resp}
		data, err = ioutil.ReadAll(resp.Body)

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

// SendWithAuth makes a request to the API and apply OAuth2 header automatically.
// If the access token soon to be expired or already expired, it will try to get a new one before
// making the main request
// client.Token will be updated when changed
func (c *Client) SendWithAuth(req *http.Request, v interface{}) error {
	if c.Token != nil {
		if c.Token.ExpiresIn < RequestNewTokenBeforeExpiresIn {
			// c.Token willbe updated in GetAccessToken call
			_, err := c.GetAccessToken()
			if err != nil {
				return err
			}
		}

		req.Header.Set("Authorization", "Bearer "+c.Token.Token)
	}

	return c.Send(req, v)
}

// NewRequest constructs a request
// Convert payload to a JSON
func (c *Client) NewRequest(method, url string, payload interface{}) (*http.Request, error) {
	var buf io.Reader
	if payload != nil {
		var b []byte
		b, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}
	return http.NewRequest(method, url, buf)
}

// log will dump request and response to the log file
func (c *Client) log(req *http.Request, resp *http.Response) {
	if c.Log != nil {
		reqDump, _ := httputil.DumpRequestOut(req, true)
		respDump, _ := httputil.DumpResponse(resp, true)

		c.Log.Write([]byte("Request: " + string(reqDump) + "\nResponse: " + string(respDump) + "\n\n"))

	}
}
