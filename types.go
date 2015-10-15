package paypalsdk

import (
    "net/http"
    "time"
    "fmt"
)


const (
    // APIBaseSandBox points to the sandbox (for testing) version of the API
    APIBaseSandBox = "https://api.sandbox.paypal.com"

    // APIBaseLive points to the live version of the API
    APIBaseLive = "https://api.paypal.com"
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

    // ErrorResponse is used when a response has errors
    ErrorResponse struct {
        Response        *http.Response `json:"-"`
        Name            string        `json:"name"`
        DebugID         string        `json:"debug_id"`
        Message         string        `json:"message"`
        InformationLink string        `json:"information_link"`
        Details         []ErrorDetail `json:"details"`
    }

    // ErrorDetails map to error_details object
    ErrorDetail struct {
        Field string `json:"field"`
        Issue string `json:"issue"`
    }
)

// Error method implementation for ErrorResponse struct
func (r *ErrorResponse) Error() string {
    return fmt.Sprintf("%v %v: %d %v\nDetails: %v", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, r.Details)
}
