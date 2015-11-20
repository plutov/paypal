package paypalsdk

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
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

func (c *Client) log(request *http.Request, response *http.Response) {
	if c.LogFile != "" {
		os.OpenFile(c.LogFile, os.O_CREATE, 0755)

		logFile, err := os.OpenFile(c.LogFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
		if err == nil {
			logFile.WriteString("URL: " + request.RequestURI + "\n\n")
		}
	}
}
