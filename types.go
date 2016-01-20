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
	// Address https://developer.paypal.com/webapps/developer/docs/api/#address-object
	Address struct {
		Line1       string `json:"line1"`
		Line2       string `json:"line2,omitempty"`
		City        string `json:"city"`
		CountryCode string `json:"country_code"`
		PostalCode  string `json:"postal_code,omitempty"`
		State       string `json:"state,omitempty"`
		Phone       string `json:"phone,omitempty"`
	}

	// Amount https://developer.paypal.com/webapps/developer/docs/api/#amount-object
	Amount struct {
		Currency string `json:"currency"`
		Total    string `json:"total"`
	}

	// Authorization rhttps://developer.paypal.com/webapps/developer/docs/api/#authorization-object
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

	// Capture https://developer.paypal.com/webapps/developer/docs/api/#capture-object
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

	// Client represents a Paypal REST API Client
	Client struct {
		client   *http.Client
		ClientID string
		Secret   string
		APIBase  string
		LogFile  string // If user set log file name all requests will be logged there
		Token    *TokenResponse
	}

	// CreditCard https://developer.paypal.com/webapps/developer/docs/api/#creditcard-object
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

	// CreditCardToken https://developer.paypal.com/webapps/developer/docs/api/#creditcardtoken-object
	CreditCardToken struct {
		CreditCardID string `json:"credit_card_id"`
		PayerID      string `json:"payer_id,omitempty"`
		Last4        string `json:"last4,omitempty"`
		ExpireYear   string `json:"expire_year,omitempty"`
		ExpireMonth  string `json:"expire_month,omitempty"`
	}

	// Currency https://developer.paypal.com/webapps/developer/docs/api/#currency-object
	Currency struct {
		Currency string `json:"currency,omitempty"`
		Value    string `json:"value,omitempty"`
	}

	// ErrorResponse https://developer.paypal.com/webapps/developer/docs/api/#error-object
	ErrorResponse struct {
		Response        *http.Response `json:"-"`
		Name            string         `json:"name"`
		DebugID         string         `json:"debug_id"`
		Message         string         `json:"message"`
		InformationLink string         `json:"information_link"`
		Details         []ErrorDetail  `json:"details"`
	}

	// ErrorDetail https://developer.paypal.com/webapps/developer/docs/api/#errordetails-object
	ErrorDetail struct {
		Field string `json:"field"`
		Issue string `json:"issue"`
	}

	// ExecuteResponse structure
	ExecuteResponse struct {
		ID    string        `json:"id"`
		Links []PaymentLink `json:"links"`
		State string        `json:"state"`
	}

	// FundingInstrument https://developer.paypal.com/webapps/developer/docs/api/#fundinginstrument-object
	FundingInstrument struct {
		CreditCard      *CreditCard      `json:"credit_card,omitempty"`
		CreditCardToken *CreditCardToken `json:"credit_card_token,omitempty"`
	}

	// Item https://developer.paypal.com/webapps/developer/docs/api/#item-object
	Item struct {
		Quantity    int    `json:"quantity"`
		Name        string `json:"name"`
		Price       string `json:"price"`
		Currency    string `json:"currency"`
		SKU         string `json:"sku,omitempty"`
		Description string `json:"description,omitempty"`
		Tax         string `json:"tax,omitempty"`
	}

	// ItemList struct
	ItemList struct {
		Items           []Item           `json:"items,omitempty"`
		ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
	}

	// Links https://developer.paypal.com/webapps/developer/docs/api/#itemlist-object
	Links struct {
		Href    string `json:"href"`
		Rel     string `json:"rel"`
		Method  string `json:"method"`
		Enctype string `json:"enctype"`
	}

	// Order https://developer.paypal.com/webapps/developer/docs/api/#order-object
	Order struct {
		ID            string     `json:"id,omitempty"`
		CreateTime    *time.Time `json:"create_time,omitempty"`
		UpdateTime    *time.Time `json:"update_time,omitempty"`
		State         string     `json:"state,omitempty"`
		Amount        *Amount    `json:"amount,omitempty"`
		PendingReason string     `json:"pending_reason,omitempty"`
		ParentPayment string     `json:"parent_payment,omitempty"`
		Links         []Links    `json:"links,omitempty"`
	}

	// Payer https://developer.paypal.com/webapps/developer/docs/api/#payer-object
	Payer struct {
		PaymentMethod      string              `json:"payment_method"`
		FundingInstruments []FundingInstrument `json:"funding_instruments,omitempty"`
		PayerInfo          *PayerInfo          `json:"payer_info,omitempty"`
		Status             string              `json:"payer_status,omitempty"`
	}

	// PayerInfo https://developer.paypal.com/webapps/developer/docs/api/#itemlist-object
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

	// Payment https://developer.paypal.com/webapps/developer/docs/api/#payment-object
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

	// PaymentLink https://developer.paypal.com/webapps/developer/docs/api/#paymentlink-object
	PaymentLink struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	}

	// PaymentResponse structure
	PaymentResponse struct {
		ID    string        `json:"id"`
		Links []PaymentLink `json:"links"`
	}

	// RedirectURLs https://developer.paypal.com/webapps/developer/docs/api/#redirecturls-object
	RedirectURLs struct {
		ReturnURL string `json:"return_url,omitempty"`
		CancelURL string `json:"cancel_url,omitempty"`
	}

	// Refund https://developer.paypal.com/webapps/developer/docs/api/#refund-object
	Refund struct {
		ID            string     `json:"id,omitempty"`
		Amount        *Amount    `json:"amount,omitempty"`
		CreateTime    *time.Time `json:"create_time,omitempty"`
		State         string     `json:"state,omitempty"`
		CaptureID     string     `json:"capture_id,omitempty"`
		ParentPayment string     `json:"parent_payment,omitempty"`
		UpdateTime    *time.Time `json:"update_time,omitempty"`
	}

	// Sale https://developer.paypal.com/webapps/developer/docs/api/#sale-object
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

	// ShippingAddress https://developer.paypal.com/webapps/developer/docs/api/#shippingaddredd-object
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

	// TokenResponse is for API response for the /oauth2/token endpoint
	TokenResponse struct {
		RefreshToken string `json:"refresh_token"`
		Token        string `json:"access_token"`
		Type         string `json:"token_type"`
	}

	// Transaction https://developer.paypal.com/webapps/developer/docs/api/#transaction-object
	Transaction struct {
		Amount         *Amount   `json:"amount"`
		Description    string    `json:"description,omitempty"`
		ItemList       *ItemList `json:"item_list,omitempty"`
		InvoiceNumber  string    `json:"invoice_number,omitempty"`
		Custom         string    `json:"custom,omitempty"`
		SoftDescriptor string    `json:"soft_descriptor,omitempty"`
	}

	// UserInfo https://developer.paypal.com/webapps/developer/docs/api/#userinfo-object
	UserInfo struct {
		ID              string   `json:"user_id"`
		Name            string   `json:"name"`
		GivenName       string   `json:"given_name"`
		FamilyName      string   `json:"family_name"`
		Email           string   `json:"email"`
		Verified        bool     `json:"verified,omitempty"`
		Gender          string   `json:"gender,omitempty"`
		BirthDate       string   `json:"birthdate,omitempty"`
		ZoneInfo        string   `json:"zoneinfo,omitempty"`
		Locale          string   `json:"locale,omitempty"`
		Phone           string   `json:"phone_number,omitempty"`
		Address         *Address `json:"address,omitempty"`
		VerifiedAccount bool     `json:"verified_account,omitempty"`
		AccountType     string   `json:"account_type,omitempty"`
		AgeRange        string   `json:"age_range,omitempty"`
		PayerID         string   `json:"payer_id,omitempty"`
	}
)

// Error method implementation for ErrorResponse struct
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v\nDetails: %v", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, r.Details)
}
