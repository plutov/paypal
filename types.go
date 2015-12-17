package paypalsdk

import (
	"fmt"
	"net/http"
	"time"
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

	// Authorization maps to the  authorization object
	Authorization struct {
		Amount                    *Amount    `json:"amount,omitempty"`
		CreateTime                *time.Time `json:"create_time,omitempty"`
		UpdateTime                *time.Time `json:"update_time,omitempty"`
		State                     string     `json:"state,omitempty"`
		ParentPayment             string     `json:"parent_payment,omitempty"`
		ID                        string     `json:"id,omitempty"`
		ValidUntil                *time.Time `json:"valid_until,omitempty"`
		Links                     []Links    `json:"links,omitempty"`
		ClearingTime              string     `json:"clearing_time,omitempty"`
		ProtectionEligibility     string     `json:"protection_eligibility,omitempty"`
		ProtectionEligibilityType string     `json:"protection_eligibility_type,omitempty"`
	}

	// Capture maps to the capture object
	Capture struct {
		Amount         *Amount    `json:"amount,omitempty"`
		IsFinalCapture bool       `json:"is_final_capture"`
		CreateTime     *time.Time `json:"create_time,omitempty"`
		UpdateTime     *time.Time `json:"update_time,omitempty"`
		State          string     `json:"state,omitempty"`
		ParentPayment  string     `json:"parent_payment,omitempty"`
		ID             string     `json:"id,omitempty"`
		Links          []Links    `json:"links,omitempty"`
	}

	// Address maps to address object
	Address struct {
		Line1       string `json:"line1"`
		Line2       string `json:"line2,omitempty"`
		City        string `json:"city"`
		CountryCode string `json:"country_code"`
		PostalCode  string `json:"postal_code,omitempty"`
		State       string `json:"state,omitempty"`
		Phone       string `json:"phone,omitempty"`
	}

	// CreditCard maps to credit_card object
	CreditCard struct {
		ID             string   `json:"id,omitempty"`
		PayerID        string   `json:"payer_id,omitempty"`
		Number         string   `json:"number"`
		Type           string   `json:"type"`
		ExpireMonth    string   `json:"expire_month"`
		ExpireYear     string   `json:"expire_year"`
		CVV2           string   `json:"cvv2,omitempty"`
		FirstName      string   `json:"first_name,omitempty"`
		LastName       string   `json:"last_name,omitempty"`
		BillingAddress *Address `json:"billing_address,omitempty"`
		State          string   `json:"state,omitempty"`
		ValidUntil     string   `json:"valid_until,omitempty"`
	}

	// CreditCardToken maps to credit_card_token object
	CreditCardToken struct {
		CreditCardID string `json:"credit_card_id"`
		PayerID      string `json:"payer_id,omitempty"`
		Last4        string `json:"last4,omitempty"`
		ExpireYear   string `json:"expire_year,omitempty"`
		ExpireMonth  string `json:"expire_month,omitempty"`
	}

	// FundingInstrument maps to funding_instrument object
	FundingInstrument struct {
		CreditCard      *CreditCard      `json:"credit_card,omitempty"`
		CreditCardToken *CreditCardToken `json:"credit_card_token,omitempty"`
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

	// Payer maps to payer object
	Payer struct {
		PaymentMethod      string              `json:"payment_method"`
		FundingInstruments []FundingInstrument `json:"funding_instruments,omitempty"`
		PayerInfo          *PayerInfo          `json:"payer_info,omitempty"`
		Status             string              `json:"payer_status,omitempty"`
	}

	// PayerInfo maps to payer_info object
	PayerInfo struct {
		Email           string           `json:"email,omitempty"`
		FirstName       string           `json:"first_name,omitempty"`
		LastName        string           `json:"last_name,omitempty"`
		PayerID         string           `json:"payer_id,omitempty"`
		Phone           string           `json:"phone,omitempty"`
		ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
		TaxIDType       string           `json:"tax_id_type,omitempty"`
		TaxID           string           `json:"tax_id,omitempty"`
	}

	// ShippingAddress maps to shipping_address object
	ShippingAddress struct {
		RecipientName string `json:"recipient_name,omitempty"`
		Type          string `json:"type,omitempty"`
		Line1         string `json:"line1"`
		Line2         string `json:"line2,omitempty"`
		City          string `json:"city"`
		CountryCode   string `json:"country_code"`
		PostalCode    string `json:"postal_code,omitempty"`
		State         string `json:"state,omitempty"`
		Phone         string `json:"phone,omitempty"`
	}

	// RedirectURLs maps to redirect_urls object
	RedirectURLs struct {
		ReturnURL string `json:"return_url,omitempty"`
		CancelURL string `json:"cancel_url,omitempty"`
	}

	// Links maps to links object
	Links struct {
		Href    string `json:"href"`
		Rel     string `json:"rel"`
		Method  string `json:"method"`
		Enctype string `json:"enctype"`
	}

	// Payment maps to payment object
	Payment struct {
		Intent              string        `json:"intent"`
		Payer               *Payer        `json:"payer"`
		Transactions        []Transaction `json:"transactions"`
		RedirectURLs        *RedirectURLs `json:"redirect_urls,omitempty"`
		ID                  string        `json:"id,omitempty"`
		CreateTime          *time.Time    `json:"create_time,omitempty"`
		State               string        `json:"state,omitempty"`
		UpdateTime          *time.Time    `json:"update_time,omitempty"`
		ExperienceProfileID string        `json:"experience_profile_id,omitempty"`
	}

	// Item maps to item object
	Item struct {
		Quantity    int    `json:"quantity"`
		Name        string `json:"name"`
		Price       string `json:"price"`
		Currency    string `json:"currency"`
		SKU         string `json:"sku,omitempty"`
		Description string `json:"description,omitempty"`
		Tax         string `json:"tax,omitempty"`
	}

	// ItemList maps to item_list object
	ItemList struct {
		Items           []Item           `json:"items,omitempty"`
		ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
	}

	// Transaction maps to transaction object
	Transaction struct {
		Amount         *Amount   `json:"amount"`
		Description    string    `json:"description,omitempty"`
		ItemList       *ItemList `json:"item_list,omitempty"`
		InvoiceNumber  string    `json:"invoice_number,omitempty"`
		Custom         string    `json:"custom,omitempty"`
		SoftDescriptor string    `json:"soft_descriptor,omitempty"`
	}

	// PaymentLink structure
	PaymentLink struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	}

	// Amount to pay
	Amount struct {
		Currency string `json:"currency"`
		Total    string `json:"total"`
	}

	// ExecuteResponse structure
	ExecuteResponse struct {
		ID    string        `json:"id"`
		Links []PaymentLink `json:"links"`
		State string        `json:"state"`
	}

	// Sale will be returned by GetSale
	Sale struct {
		ID                        string     `json:"id,omitempty"`
		Amount                    *Amount    `json:"amount,omitempty"`
		Description               string     `json:"description,omitempty"`
		CreateTime                *time.Time `json:"create_time,omitempty"`
		State                     string     `json:"state,omitempty"`
		ParentPayment             string     `json:"parent_payment,omitempty"`
		UpdateTime                *time.Time `json:"update_time,omitempty"`
		PaymentMode               string     `json:"payment_mode,omitempty"`
		PendingReason             string     `json:"pending_reason,omitempty"`
		ReasonCode                string     `json:"reason_code,omitempty"`
		ClearingTime              string     `json:"clearing_time,omitempty"`
		ProtectionEligibility     string     `json:"protection_eligibility,omitempty"`
		ProtectionEligibilityType string     `json:"protection_eligibility_type,omitempty"`
		Links                     []Links    `json:"links,omitempty"`
	}

	// Refund will be returned by RefundSale
	Refund struct {
		ID            string     `json:"id,omitempty"`
		Amount        *Amount    `json:"amount,omitempty"`
		CreateTime    *time.Time `json:"create_time,omitempty"`
		State         string     `json:"state,omitempty"`
		CaptureID     string     `json:"capture_id,omitempty"`
		ParentPayment string     `json:"parent_payment,omitempty"`
		UpdateTime    *time.Time `json:"update_time,omitempty"`
	}
)

// Error method implementation for ErrorResponse struct
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v\nDetails: %v", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, r.Details)
}
