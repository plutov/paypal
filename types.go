package paypalsdk

import (
	"fmt"
	"net/http"
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
		LogFile  string // If user set log file name all requests will be logged there
		Token    *TokenResponse
	}

	// TokenResponse maps to the API response for the /oauth2/token endpoint
	TokenResponse struct {
		RefreshToken string `json:"refresh_token"`
		Token        string `json:"access_token"`
		Type         string `json:"token_type"`
	}

	// ErrorResponse is used when a response has errors
	ErrorResponse struct {
		Response        *http.Response `json:"-"`
		Name            string         `json:"name"`
		DebugID         string         `json:"debug_id"`
		Message         string         `json:"message"`
		InformationLink string         `json:"information_link"`
		Details         []ErrorDetail  `json:"details"`
	}

	// ErrorDetail map to error_details object
	ErrorDetail struct {
		Field string `json:"field"`
		Issue string `json:"issue"`
	}

	// PaymentResponse structure
	PaymentResponse struct {
		ID    string        `json:"id"`
		Links []PaymentLink `json:"links"`
	}

	// PaymentLink structure
	PaymentLink struct {
		Href string `json:"href"`
	}

	// Amount to pay
	Amount struct {
		Currency string
		Total    float64
	}

	// ExecuteResponse structure
	ExecuteResponse struct {
		ID    string        `json:"id"`
		Links []PaymentLink `json:"links"`
		State string        `json:"state"`
	}
)

// Error method implementation for ErrorResponse struct
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v\nDetails: %v", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, r.Details)
}
