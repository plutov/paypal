package paypalsdk

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
)

// NewClient returns new Client struct
func NewClient(clientID string, secret string, APIBase string) (*Client, error) {
	if clientID == "" || secret == "" || APIBase == "" {
		return &Client{}, errors.New("ClientID, Secret and APIBase are required to create a Client")
	}

	return &Client{
		&http.Client{},
		clientID,
		secret,
		APIBase,
		"",
		nil,
	}, nil
}

// SetLogFile func
func (c *Client) SetLogFile(filepath string) error {
	c.LogFile = filepath

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
	// Set default headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en_US")

	// Default values for headers
	if req.Header.Get("Content-type") == "" {
		req.Header.Set("Content-type", "application/json")
	}

	resp, err := c.client.Do(req)
	c.log(req, resp)

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

// SendWithAuth makes a request to the API and apply OAuth2 header automatically.
// If the access token soon to be expired, it will try to get a new one before
// making the main request
func (c *Client) SendWithAuth(req *http.Request, v interface{}) error {
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)

	return c.Send(req, v)
}

func (c *Client) log(req *http.Request, resp *http.Response) {
	if c.LogFile != "" {
		os.OpenFile(c.LogFile, os.O_CREATE, 0755)

		logFile, err := os.OpenFile(c.LogFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
		if err == nil {
			reqDump, _ := httputil.DumpRequestOut(req, true)
			respDump, _ := httputil.DumpResponse(resp, true)

			logFile.WriteString("Request: " + string(reqDump) + "\nResponse: " + string(respDump) + "\n\n")
		}
	}
}
