package paypal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/plutov/paypal/v4/limiter"
	"github.com/stretchr/testify/assert"
)

type okResp struct {
	OK bool `json:"ok"`
}

func newTestServer() (*httptest.Server, *int32, *int32) {
	tokenHits := new(int32)
	apiHits := new(int32)
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(tokenHits, 1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "test-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	})
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(apiHits, 1)
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(okResp{OK: true})
	})
	return httptest.NewServer(mux), tokenHits, apiHits
}

func newClientForServer(t *testing.T, s *httptest.Server) *Client {
	t.Helper()
	c, err := NewClient("id", "secret", s.URL)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	c.SetHTTPClient(s.Client())
	return c
}

func TestClientLimiter_BlocksWhenInsufficientPermitsForTokenAndRequest(t *testing.T) {
	s, tokenHits, apiHits := newTestServer()
	defer s.Close()
	c := newClientForServer(t, s)

	// limiter allows only 1 request, but we need 2 (token + api)
	lim := limiter.NewFixedWindowLimiter(limiter.NewMemoryStorage(), limiter.FixedWindowConfig{Window: time.Minute, Limit: 1})
	c.SetLimiter(lim, "test-key")

	req, _ := c.NewRequest(context.Background(), http.MethodGet, s.URL+"/api", nil)
	var out okResp
	err := c.SendWithAuth(req, &out)
	assert.Error(t, err)
	assert.Equal(t, ErrRateLimited, err)
	// should have short-circuited before hitting network
	assert.Equal(t, int32(0), atomic.LoadInt32(tokenHits))
	assert.Equal(t, int32(0), atomic.LoadInt32(apiHits))
}

func TestClientLimiter_AllowsWithSufficientPermitsForTokenAndRequest(t *testing.T) {
	s, tokenHits, apiHits := newTestServer()
	defer s.Close()
	c := newClientForServer(t, s)

	// need 2 permits and we have 2
	lim := limiter.NewFixedWindowLimiter(limiter.NewMemoryStorage(), limiter.FixedWindowConfig{Window: time.Minute, Limit: 2})
	c.SetLimiter(lim, "test-key")

	req, _ := c.NewRequest(context.Background(), http.MethodGet, s.URL+"/api", nil)
	var out okResp
	err := c.SendWithAuth(req, &out)
	assert.NoError(t, err)
	assert.True(t, out.OK)
	assert.Equal(t, int32(1), atomic.LoadInt32(tokenHits))
	assert.Equal(t, int32(1), atomic.LoadInt32(apiHits))
}

func TestClientLimiter_ConsumesOnePermitWhenTokenAlreadySet(t *testing.T) {
	s, tokenHits, apiHits := newTestServer()
	defer s.Close()
	c := newClientForServer(t, s)

	// set token so we don't need refresh
	c.SetAccessToken("preset")
	// allow only 1 permit
	lim := limiter.NewFixedWindowLimiter(limiter.NewMemoryStorage(), limiter.FixedWindowConfig{Window: time.Minute, Limit: 1})
	c.SetLimiter(lim, "test-key")

	req, _ := c.NewRequest(context.Background(), http.MethodGet, s.URL+"/api", nil)
	var out okResp
	err := c.SendWithAuth(req, &out)
	assert.NoError(t, err)
	assert.True(t, out.OK)
	assert.Equal(t, int32(0), atomic.LoadInt32(*&tokenHits))
	assert.Equal(t, int32(1), atomic.LoadInt32(apiHits))
}

func TestClientLimiter_TokenSoonToExpire_ConsumesTwoPermits(t *testing.T) {
	s, tokenHits, apiHits := newTestServer()
	defer s.Close()
	c := newClientForServer(t, s)

	// pre-existing token but expiring soon (< 60s)
	c.Token = &TokenResponse{Token: "old"}
	c.tokenExpiresAt = time.Now().Add(5 * time.Second)

	// need 2 permits; provide exactly 2
	lim := limiter.NewFixedWindowLimiter(limiter.NewMemoryStorage(), limiter.FixedWindowConfig{Window: time.Minute, Limit: 2})
	c.SetLimiter(lim, "test-key")

	req, _ := c.NewRequest(context.Background(), http.MethodGet, s.URL+"/api", nil)
	var out okResp
	err := c.SendWithAuth(req, &out)
	assert.NoError(t, err)
	assert.True(t, out.OK)
	assert.Equal(t, int32(1), atomic.LoadInt32(tokenHits))
	assert.Equal(t, int32(1), atomic.LoadInt32(apiHits))
}
