package paypalsdk

import (
    "net/http"
    "time"
)

const (
    // APIBaseSandBox points to the sandbox (for testing) version of the API
    APIBaseSandBox = "https://api.sandbox.paypal.com/v1"

    // APIBaseLive points to the live version of the API
    APIBaseLive = "https://api.paypal.com/v1"
)

type (
    // Client represents a Paypal REST API Client
    Client struct {
        client   *http.Client
        ClientID string
        Secret   string
        APIBase  string
        Token    *TokenResponse
    }

    // TokenResponse maps to the API response for the /oauth2/token endpoint
    TokenResponse struct {
        Scope     string    `json:"scope"`        // "https://api.paypal.com/v1/payments/.* https://api.paypal.com/v1/vault/credit-card https://api.paypal.com/v1/vault/credit-card/.*",
        Token     string    `json:"access_token"` // "EEwJ6tF9x5WCIZDYzyZGaz6Khbw7raYRIBV_WxVvgmsG",
        Type      string    `json:"token_type"`   // "Bearer",
        AppID     string    `json:"app_id"`       // "APP-6XR95014BA15863X",
        ExpiresIn int       `json:"expires_in"`   // 28800
        ExpiresAt time.Time `json:"expires_at"`
    }
)
