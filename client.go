package paypalsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

// NewClient returns new Client struct
// APIBase is a base API URL, for testing you can use paypalsdk.APIBaseSandBox
func NewClient(clientID string, secret string, APIBase string) (*Client, error) {
	if clientID == "" || secret == "" || APIBase == "" {
		return nil, errors.New("ClientID, Secret and APIBase are required to create a Client")
	}

	return &Client{
		Client:   &http.Client{},
		ClientID: clientID,
		Secret:   secret,
		APIBase:  APIBase,
	}, nil
}

// GetAccessToken returns struct of TokenResponse
// No need to call SetAccessToken to apply new access token for current Client
// Endpoint: POST /v1/oauth2/token
func (c *Client) GetAccessToken() (*TokenResponse, error) {
	buf := bytes.NewBuffer([]byte("grant_type=client_credentials"))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/oauth2/token"), buf)
	if err != nil {
		return &TokenResponse{}, err
	}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	t := TokenResponse{}
	err = c.SendWithBasicAuth(req, &t)

	// Set Token fur current Client
	if t.Token != "" {
		c.Token = &t
		c.tokenExpiresAt = time.Now().Add(time.Duration(t.ExpiresIn) * time.Second)
	}

	return &t, err
}

// SetHTTPClient sets *http.Client to current client
func (c *Client) SetHTTPClient(client *http.Client) error {
	c.Client = client
	return nil
}

// SetAccessToken sets saved token to current client
func (c *Client) SetAccessToken(token string) error {
	c.Token = &TokenResponse{
		Token: token,
	}
	c.tokenExpiresAt = time.Time{}

	return nil
}

// SetLog will set/change the output destination.
// If log file is set paypalsdk will log all requests and responses to this Writer
func (c *Client) SetLog(log io.Writer) error {
	c.Log = log
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

	resp, err = c.Client.Do(req)
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
		if !c.tokenExpiresAt.IsZero() && c.tokenExpiresAt.Sub(time.Now()) < RequestNewTokenBeforeExpiresIn {
			// c.Token will be updated in GetAccessToken call
			if _, err := c.GetAccessToken(); err != nil {
				return err
			}
		}

		req.Header.Set("Authorization", "Bearer "+c.Token.Token)
	}

	return c.Send(req, v)
}

// SendWithBasicAuth makes a request to the API using clientID:secret basic auth
func (c *Client) SendWithBasicAuth(req *http.Request, v interface{}) error {
	req.SetBasicAuth(c.ClientID, c.Secret)

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
func (c *Client) log(r *http.Request, resp *http.Response) {
	if c.Log != nil {
		var (
			reqDump  string
			respDump []byte
		)

		if r != nil {
			reqDump = fmt.Sprintf("%s %s. Data: %s", r.Method, r.URL.String(), r.Form.Encode())
		}
		if resp != nil {
			respDump, _ = httputil.DumpResponse(resp, true)
		}

		c.Log.Write([]byte(fmt.Sprintf("Request: %s\nResponse: %s\n", reqDump, string(respDump))))
	}
}
