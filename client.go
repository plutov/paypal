package paypal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/plutov/paypal/v4/limiter"
)

var ErrRateLimited = errors.New("rate limited")

type Option func(o *Client) *Client

func WithCache(cache Cache) Option {
	return func(c *Client) *Client {
		c.cache = cache
		return c
	}
}

// NewClient returns new Client struct
// APIBase is a base API URL, for testing you can use paypal.APIBaseSandBox
func NewClient(ctx context.Context, clientID string, secret string, APIBase string, options ...Option) (*Client, error) {
	if clientID == "" || secret == "" || APIBase == "" {
		return nil, errors.New("ClientID, Secret and APIBase are required to create a Client")
	}

	cli := &Client{
		Client:   &http.Client{},
		ClientID: clientID,
		Secret:   secret,
		APIBase:  APIBase,
	}

	for _, option := range options {
		cli = option(cli)
	}
	cli.checkCacheToken(ctx)
	return cli, nil
}

func (c *Client) checkCacheToken(ctx context.Context) {
	if c.cache == nil {
		return
	}

	token, err := c.cache.GetToken(ctx, c.ClientID)
	if err != nil {
		return
	}
	if token != nil && token.Token != "" {
		c.SetAccessTokenWithExpiry(token.Token, token.ExpiresAt)
	}
}

// GetAccessToken returns struct of TokenResponse
// No need to call SetAccessToken to apply new access token for current Client
// Endpoint: POST /v1/oauth2/token
func (c *Client) GetAccessToken(ctx context.Context) (*TokenResponse, error) {
	buf := bytes.NewBuffer([]byte("grant_type=client_credentials"))
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/oauth2/token"), buf)
	if err != nil {
		return &TokenResponse{}, err
	}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	response := &TokenResponse{}
	err = c.SendWithBasicAuth(req, response)

	// Set Token fur current Client
	if response.Token != "" {
		c.Token = response
		c.tokenExpiresAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)

		if c.cache != nil {
			token := &CacheToken{
				RefreshToken: response.RefreshToken,
				Token:        response.Token,
				Type:         response.Type,
				ExpiresIn:    int64(response.ExpiresIn),
				ExpiresAt:    c.tokenExpiresAt.Unix(),
			}
			err = c.cache.SetToken(ctx, c.ClientID, token)
			if err != nil {
				return response, err
			}
		}
	}

	return response, err
}

// SetHTTPClient sets *http.Client to current client
func (c *Client) SetHTTPClient(client *http.Client) {
	c.Client = client
}

// SetAccessToken sets saved token to current client
func (c *Client) SetAccessToken(token string) {
	c.Token = &TokenResponse{
		Token: token,
	}
	c.tokenExpiresAt = time.Time{}
}

// SetAccessTokenWithExpiry sets saved token to current client with expiry time
func (c *Client) SetAccessTokenWithExpiry(token string, expiresAt int64) {
	c.Token = &TokenResponse{
		Token: token,
	}
	if expiresAt > 0 {
		c.tokenExpiresAt = time.Unix(expiresAt, 0)
	}
}

// SetLog will set/change the output destination.
// If log file is set paypal will log all requests and responses to this Writer
func (c *Client) SetLog(log io.Writer) {
	c.Log = log
}

// SetReturnRepresentation enables verbose response
// Verbose response: https://developer.paypal.com/docs/api/orders/v2/#orders-authorize-header-parameters
func (c *Client) SetReturnRepresentation() {
	c.returnRepresentation = true
}

// SetLimiter configures an optional rate limiter used by SendWithAuth.
// key is used as the rate limit bucket key; if empty, a default is used.
func (c *Client) SetLimiter(l limiter.RateLimiter, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rateLimiter = l
	c.rateLimiterKey = key
}

// Send makes a request to the API, the response body will be
// unmarshalled into v, or if v is an io.Writer, the response will
// be written to it without decoding
func (c *Client) Send(req *http.Request, v interface{}) (retErr error) {
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
	if c.returnRepresentation {
		req.Header.Set("Prefer", "return=representation")
	}
	if c.Log != nil {
		if reqDump, err := httputil.DumpRequestOut(req, true); err == nil {
			logMsg := fmt.Sprintf("Request: %s\n", string(reqDump))
			if _, logErr := c.Log.Write([]byte(logMsg)); logErr != nil {
				return logErr
			}
		}
	}

	resp, err = c.Client.Do(req)
	if err != nil {
		return err
	}

	if c.Log != nil {
		if respDump, err := httputil.DumpResponse(resp, true); err == nil {
			logMsg := fmt.Sprintf("Response from %s: %s\n", req.URL.RequestURI(), string(respDump))
			if _, logErr := c.Log.Write([]byte(logMsg)); logErr != nil {
				return logErr
			}
		}
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil && retErr == nil {
			retErr = err
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errResp := &ErrorResponse{Response: resp}

		var bodyBuffer bytes.Buffer
		teeReader := io.TeeReader(resp.Body, &bodyBuffer)
		data, err = io.ReadAll(teeReader)
		resp.Body = io.NopCloser(&bodyBuffer)

		if err == nil && len(data) > 0 {
			err := json.Unmarshal(data, errResp)
			if err != nil {
				return err
			}
		}

		return errResp
	}
	if v == nil {
		return nil
	}

	if w, ok := v.(io.Writer); ok {
		_, err := io.Copy(w, resp.Body)
		return err
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// SendWithAuth makes a request to the API and apply OAuth2 header automatically.
// If the access token soon to be expired or already expired, it will try to get a new one before
// making the main request
// client.Token will be updated when changed
func (c *Client) SendWithAuth(req *http.Request, v interface{}) error {
	// c.Lock()
	c.mu.Lock()
	// Note: Here we do not want to `defer c.Unlock()` because we need `c.Send(...)`
	// to happen outside of the locked section.

	// determine if a token request is needed under the lock
	needToken := c.Token == nil || (!c.tokenExpiresAt.IsZero() && time.Until(c.tokenExpiresAt) < RequestNewTokenBeforeExpiresIn)

	// optional rate limiting
	if c.rateLimiter != nil {
		key := c.rateLimiterKey
		if key == "" {
			key = "paypal:client"
		}
		ctx := req.Context()
		permits := 1
		if needToken {
			permits = 2 // one for token request, one for the main API call
		}
		for i := 0; i < permits; i++ {
			dec, err := c.rateLimiter.Allow(ctx, key)
			if err != nil {
				c.mu.Unlock()
				return err
			}
			if !dec.Allowed {
				c.mu.Unlock()
				return ErrRateLimited
			}
		}
	}

	if needToken {
		// c.Token will be updated in GetAccessToken call
		if _, err := c.GetAccessToken(req.Context()); err != nil {
			// c.Unlock()
			c.mu.Unlock()
			return err
		}
	}

	req.Header.Set("Authorization", "Bearer "+c.Token.Token)
	// Unlock the client mutex before sending the request, this allows multiple requests
	// to be in progress at the same time.
	// c.Unlock()
	c.mu.Unlock()
	return c.Send(req, v)
}

// SendWithBasicAuth makes a request to the API using clientID:secret basic auth
func (c *Client) SendWithBasicAuth(req *http.Request, v interface{}) error {
	req.SetBasicAuth(c.ClientID, c.Secret)

	return c.Send(req, v)
}

// NewRequest constructs a request
// Convert payload to a JSON
func (c *Client) NewRequest(ctx context.Context, method, url string, payload interface{}) (*http.Request, error) {
	var buf io.Reader
	if payload != nil {
		b, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}
	return http.NewRequestWithContext(ctx, method, url, buf)
}
