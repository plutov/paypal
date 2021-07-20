package paypal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	// APIBaseSandBox points to the sandbox (for testing) version of the API
	APIBaseSandBox = "https://api.sandbox.paypal.com"

	// APIBaseLive points to the live version of the API
	APIBaseLive = "https://api.paypal.com"

	// RequestNewTokenBeforeExpiresIn is used by SendWithAuth and try to get new Token when it's about to expire
	RequestNewTokenBeforeExpiresIn = time.Duration(60) * time.Second
)

// Possible values for `no_shipping` in InputFields
//
// https://developer.paypal.com/docs/api/payment-experience/#definition-input_fields
const (
	NoShippingDisplay      uint = 0
	NoShippingHide         uint = 1
	NoShippingBuyerAccount uint = 2
)

// Possible values for `address_override` in InputFields
//
// https://developer.paypal.com/docs/api/payment-experience/#definition-input_fields
const (
	AddrOverrideFromFile uint = 0
	AddrOverrideFromCall uint = 1
)

// Possible values for `landing_page_type` in FlowConfig
//
// https://developer.paypal.com/docs/api/payment-experience/#definition-flow_config
const (
	LandingPageTypeBilling string = "Billing"
	LandingPageTypeLogin   string = "Login"
)

// Possible value for `allowed_payment_method` in PaymentOptions
//
// https://developer.paypal.com/docs/api/payments/#definition-payment_options
const (
	AllowedPaymentUnrestricted         string = "UNRESTRICTED"
	AllowedPaymentInstantFundingSource string = "INSTANT_FUNDING_SOURCE"
	AllowedPaymentImmediatePay         string = "IMMEDIATE_PAY"
)

// Possible value for `intent` in CreateOrder
//
// https://developer.paypal.com/docs/api/orders/v2/#orders_create
const (
	OrderIntentCapture   string = "CAPTURE"
	OrderIntentAuthorize string = "AUTHORIZE"
)

// Possible value for `status` in GetOrder
//
// https://developer.paypal.com/docs/api/orders/v2/#orders-get-response
const (
	OrderStatusCreated   string = "CREATED"
	OrderStatusSaved     string = "SAVED"
	OrderStatusApproved  string = "APPROVED"
	OrderStatusVoided    string = "VOIDED"
	OrderStatusCompleted string = "COMPLETED"
)

// Possible values for `category` in Item
//
// https://developer.paypal.com/docs/api/orders/v2/#definition-item
const (
	ItemCategoryDigitalGood  string = "DIGITAL_GOODS"
	ItemCategoryPhysicalGood string = "PHYSICAL_GOODS"
)

// Possible values for `shipping_preference` in ApplicationContext
//
// https://developer.paypal.com/docs/api/orders/v2/#definition-application_context

const (
	EventCheckoutOrderApproved         string = "CHECKOUT.ORDER.APPROVED"
	EventPaymentCaptureCompleted       string = "PAYMENT.CAPTURE.COMPLETED"
	EventPaymentCaptureDenied          string = "PAYMENT.CAPTURE.DENIED"
	EventPaymentCaptureRefunded        string = "PAYMENT.CAPTURE.REFUNDED"
	EventMerchantOnboardingCompleted   string = "MERCHANT.ONBOARDING.COMPLETED"
	EventMerchantPartnerConsentRevoked string = "MERCHANT.PARTNER-CONSENT.REVOKED"
)

const (
	OperationAPIIntegration   string = "API_INTEGRATION"
	ProductExpressCheckout    string = "EXPRESS_CHECKOUT"
	IntegrationMethodPayPal   string = "PAYPAL"
	IntegrationTypeThirdParty string = "THIRD_PARTY"
	ConsentShareData          string = "SHARE_DATA_CONSENT"
)

const (
	FeaturePayment               string = "PAYMENT"
	FeatureRefund                string = "REFUND"
	FeatureFuturePayment         string = "FUTURE_PAYMENT"
	FeatureDirectPayment         string = "DIRECT_PAYMENT"
	FeaturePartnerFee            string = "PARTNER_FEE"
	FeatureDelayFunds            string = "DELAY_FUNDS_DISBURSEMENT"
	FeatureReadSellerDispute     string = "READ_SELLER_DISPUTE"
	FeatureUpdateSellerDispute   string = "UPDATE_SELLER_DISPUTE"
	FeatureDisputeReadBuyer      string = "DISPUTE_READ_BUYER"
	FeatureUpdateCustomerDispute string = "UPDATE_CUSTOMER_DISPUTES"
)

// https://developer.paypal.com/docs/api/payments.payouts-batch/v1/?mark=recipient_type#definition-recipient_type
const (
	EmailRecipientType    string = "EMAIL"     // An unencrypted email — string of up to 127 single-byte characters.
	PaypalIdRecipientType string = "PAYPAL_ID" // An encrypted PayPal account number.
	PhoneRecipientType    string = "PHONE"     // An unencrypted phone number.
	// Note: The PayPal sandbox doesn't support type PHONE
)

// https://developer.paypal.com/docs/api/payments.payouts-batch/v1/?mark=recipient_wallet#definition-recipient_wallet
const (
	PaypalRecipientWallet string = "PAYPAL"
	VenmoRecipientWallet  string = "VENMO"
)

// Possible value for `batch_status` in GetPayout
//
// https://developer.paypal.com/docs/api/payments.payouts-batch/v1/#definition-batch_status
const (
	BatchStatusDenied     string = "DENIED"
	BatchStatusPending    string = "PENDING"
	BatchStatusProcessing string = "PROCESSING"
	BatchStatusSuccess    string = "SUCCESS"
	BatchStatusCanceled   string = "CANCELED"
)

const (
	LinkRelSelf      string = "self"
	LinkRelActionURL string = "action_url"
)

const (
	AncorTypeApplication string = "APPLICATION"
	AncorTypeAccount     string = "ACCOUNT"
)

type (
	// JSONTime overrides MarshalJson method to format in ISO8601
	JSONTime time.Time

	// Address struct
	Address struct {
		Line1       string `json:"line1,omitempty"`
		Line2       string `json:"line2,omitempty"`
		City        string `json:"city,omitempty"`
		CountryCode string `json:"country_code,omitempty"`
		PostalCode  string `json:"postal_code,omitempty"`
		State       string `json:"state,omitempty"`
		Phone       string `json:"phone,omitempty"`
	}

	// AgreementDetails struct
	AgreementDetails struct {
		OutstandingBalance AmountPayout `json:"outstanding_balance"`
		CyclesRemaining    int          `json:"cycles_remaining,string"`
		CyclesCompleted    int          `json:"cycles_completed,string"`
		NextBillingDate    time.Time    `json:"next_billing_date"`
		LastPaymentDate    time.Time    `json:"last_payment_date"`
		LastPaymentAmount  AmountPayout `json:"last_payment_amount"`
		FinalPaymentDate   time.Time    `json:"final_payment_date"`
		FailedPaymentCount int          `json:"failed_payment_count,string"`
	}

	// Amount struct
	Amount struct {
		Currency string  `json:"currency"`
		Total    string  `json:"total"`
		Details  Details `json:"details,omitempty"`
	}

	// AmountPayout struct
	AmountPayout struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	}

	// ApplicationContext struct
	//Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#definition-application_context
	ApplicationContext struct {
		BrandName          string             `json:"brand_name,omitempty"`
		Locale             string             `json:"locale,omitempty"`
		ShippingPreference ShippingPreference `json:"shipping_preference,omitempty"`
		UserAction         UserAction         `json:"user_action,omitempty"`
		//LandingPage        string `json:"landing_page,omitempty"` // not found in documentation
		ReturnURL string `json:"return_url,omitempty"`
		CancelURL string `json:"cancel_url,omitempty"`
	}

	// Authorization struct
	Authorization struct {
		ID               string                `json:"id,omitempty"`
		CustomID         string                `json:"custom_id,omitempty"`
		InvoiceID        string                `json:"invoice_id,omitempty"`
		Status           string                `json:"status,omitempty"`
		StatusDetails    *CaptureStatusDetails `json:"status_details,omitempty"`
		Amount           *PurchaseUnitAmount   `json:"amount,omitempty"`
		SellerProtection *SellerProtection     `json:"seller_protection,omitempty"`
		CreateTime       *time.Time            `json:"create_time,omitempty"`
		UpdateTime       *time.Time            `json:"update_time,omitempty"`
		ExpirationTime   *time.Time            `json:"expiration_time,omitempty"`
		Links            []Link                `json:"links,omitempty"`
	}

	// AuthorizeOrderResponse .
	AuthorizeOrderResponse struct {
		CreateTime    *time.Time             `json:"create_time,omitempty"`
		UpdateTime    *time.Time             `json:"update_time,omitempty"`
		ID            string                 `json:"id,omitempty"`
		Status        string                 `json:"status,omitempty"`
		Intent        string                 `json:"intent,omitempty"`
		PurchaseUnits []PurchaseUnitRequest  `json:"purchase_units,omitempty"`
		Payer         *PayerWithNameAndPhone `json:"payer,omitempty"`
	}

	// AuthorizeOrderRequest - https://developer.paypal.com/docs/api/orders/v2/#orders_authorize
	AuthorizeOrderRequest struct {
		PaymentSource      *PaymentSource     `json:"payment_source,omitempty"`
		ApplicationContext ApplicationContext `json:"application_context,omitempty"`
	}

	// https://developer.paypal.com/docs/api/payments/v2/#definition-platform_fee
	PlatformFee struct {
		Amount *Money          `json:"amount,omitempty"`
		Payee  *PayeeForOrders `json:"payee,omitempty"`
	}

	// https://developer.paypal.com/docs/api/payments/v2/#definition-payment_instruction
	PaymentInstruction struct {
		PlatformFees     []PlatformFee `json:"platform_fees,omitempty"`
		DisbursementMode string        `json:"disbursement_mode,omitempty"`
	}

	// https://developer.paypal.com/docs/api/payments/v2/#authorizations_capture
	PaymentCaptureRequest struct {
		InvoiceID      string `json:"invoice_id,omitempty"`
		NoteToPayer    string `json:"note_to_payer,omitempty"`
		SoftDescriptor string `json:"soft_descriptor,omitempty"`
		Amount         *Money `json:"amount,omitempty"`
		FinalCapture   bool   `json:"final_capture,omitempty"`
	}

	SellerProtection struct {
		Status            string   `json:"status,omitempty"`
		DisputeCategories []string `json:"dispute_categories,omitempty"`
	}

	// https://developer.paypal.com/docs/api/payments/v2/#definition-capture_status_details
	CaptureStatusDetails struct {
		Reason string `json:"reason,omitempty"`
	}

	PaymentCaptureResponse struct {
		Status           string                `json:"status,omitempty"`
		StatusDetails    *CaptureStatusDetails `json:"status_details,omitempty"`
		ID               string                `json:"id,omitempty"`
		Amount           *Money                `json:"amount,omitempty"`
		InvoiceID        string                `json:"invoice_id,omitempty"`
		FinalCapture     bool                  `json:"final_capture,omitempty"`
		DisbursementMode string                `json:"disbursement_mode,omitempty"`
		Links            []Link                `json:"links,omitempty"`
	}

	//https://developer.paypal.com/docs/api/payments/v2/#captures_get
	CaptureDetailsResponse struct {
		Status                    string                     `json:"status,omitempty"`
		StatusDetails             *CaptureStatusDetails      `json:"status_details,omitempty"`
		ID                        string                     `json:"id,omitempty"`
		Amount                    *Money                     `json:"amount,omitempty"`
		InvoiceID                 string                     `json:"invoice_id,omitempty"`
		CustomID                  string                     `json:"custom_id,omitempty"`
		SellerProtection          *SellerProtection          `json:"seller_protection,omitempty"`
		FinalCapture              bool                       `json:"final_capture,omitempty"`
		SellerReceivableBreakdown *SellerReceivableBreakdown `json:"seller_receivable_breakdown,omitempty"`
		DisbursementMode          string                     `json:"disbursement_mode,omitempty"`
		Links                     []Link                     `json:"links,omitempty"`
		UpdateTime                *time.Time                 `json:"update_time,omitempty"`
		CreateTime                *time.Time                 `json:"create_time,omitempty"`
	}

	// CaptureOrderRequest - https://developer.paypal.com/docs/api/orders/v2/#orders_capture
	CaptureOrderRequest struct {
		PaymentSource *PaymentSource `json:"payment_source"`
	}

	// RefundOrderRequest - https://developer.paypal.com/docs/api/payments/v2/#captures_refund
	RefundCaptureRequest struct {
		Amount      *Money `json:"amount,omitempty"`
		InvoiceID   string `json:"invoice_id,omitempty"`
		NoteToPayer string `json:"note_to_payer,omitempty"`
	}

	// BatchHeader struct
	BatchHeader struct {
		Amount            *AmountPayout      `json:"amount,omitempty"`
		Fees              *AmountPayout      `json:"fees,omitempty"`
		PayoutBatchID     string             `json:"payout_batch_id,omitempty"`
		BatchStatus       string             `json:"batch_status,omitempty"`
		TimeCreated       *time.Time         `json:"time_created,omitempty"`
		TimeCompleted     *time.Time         `json:"time_completed,omitempty"`
		SenderBatchHeader *SenderBatchHeader `json:"sender_batch_header,omitempty"`
	}

	// BillingAgreement struct
	BillingAgreement struct {
		Name                        string               `json:"name,omitempty"`
		Description                 string               `json:"description,omitempty"`
		StartDate                   JSONTime             `json:"start_date,omitempty"`
		Plan                        BillingPlan          `json:"plan,omitempty"`
		Payer                       Payer                `json:"payer,omitempty"`
		ShippingAddress             *ShippingAddress     `json:"shipping_address,omitempty"`
		OverrideMerchantPreferences *MerchantPreferences `json:"override_merchant_preferences,omitempty"`
	}

	// BillingInfo struct
	BillingInfo struct {
		OutstandingBalance  AmountPayout      `json:"outstanding_balance,omitempty"`
		CycleExecutions     []CycleExecutions `json:"cycle_executions,omitempty"`
		LastPayment         LastPayment       `json:"last_payment,omitempty"`
		NextBillingTime     time.Time         `json:"next_billing_time,omitempty"`
		FailedPaymentsCount int               `json:"failed_payments_count,omitempty"`
	}

	// BillingPlan struct
	BillingPlan struct {
		ID                  string               `json:"id,omitempty"`
		Name                string               `json:"name,omitempty"`
		Description         string               `json:"description,omitempty"`
		Type                string               `json:"type,omitempty"`
		PaymentDefinitions  []PaymentDefinition  `json:"payment_definitions,omitempty"`
		MerchantPreferences *MerchantPreferences `json:"merchant_preferences,omitempty"`
	}

	// Capture struct
	Capture struct {
		ID             string     `json:"id,omitempty"`
		Amount         *Amount    `json:"amount,omitempty"`
		State          string     `json:"state,omitempty"`
		ParentPayment  string     `json:"parent_payment,omitempty"`
		TransactionFee string     `json:"transaction_fee,omitempty"`
		IsFinalCapture bool       `json:"is_final_capture"`
		CreateTime     *time.Time `json:"create_time,omitempty"`
		UpdateTime     *time.Time `json:"update_time,omitempty"`
		Links          []Link     `json:"links,omitempty"`
	}

	// ChargeModel struct
	ChargeModel struct {
		Type   string       `json:"type,omitempty"`
		Amount AmountPayout `json:"amount,omitempty"`
	}

	// Client represents a Paypal REST API Client
	Client struct {
		sync.Mutex
		Client               *http.Client
		ClientID             string
		Secret               string
		APIBase              string
		Log                  io.Writer // If user set log file name all requests will be logged there
		Token                *TokenResponse
		tokenExpiresAt       time.Time
		returnRepresentation bool
	}

	// CreditCard struct
	CreditCard struct {
		ID                 string   `json:"id,omitempty"`
		PayerID            string   `json:"payer_id,omitempty"`
		ExternalCustomerID string   `json:"external_customer_id,omitempty"`
		Number             string   `json:"number"`
		Type               string   `json:"type"`
		ExpireMonth        string   `json:"expire_month"`
		ExpireYear         string   `json:"expire_year"`
		CVV2               string   `json:"cvv2,omitempty"`
		FirstName          string   `json:"first_name,omitempty"`
		LastName           string   `json:"last_name,omitempty"`
		BillingAddress     *Address `json:"billing_address,omitempty"`
		State              string   `json:"state,omitempty"`
		ValidUntil         string   `json:"valid_until,omitempty"`
	}

	// CreditCards GET /v1/vault/credit-cards
	CreditCards struct {
		Items []CreditCard `json:"items"`
		SharedListResponse
	}

	// CreditCardToken struct
	CreditCardToken struct {
		CreditCardID string `json:"credit_card_id"`
		PayerID      string `json:"payer_id,omitempty"`
		Last4        string `json:"last4,omitempty"`
		ExpireYear   string `json:"expire_year,omitempty"`
		ExpireMonth  string `json:"expire_month,omitempty"`
	}

	// CreditCardsFilter struct
	CreditCardsFilter struct {
		PageSize int
		Page     int
	}

	// CreditCardField PATCH /v1/vault/credit-cards/credit_card_id
	CreditCardField struct {
		Operation string `json:"op"`
		Path      string `json:"path"`
		Value     string `json:"value"`
	}

	// Currency struct
	Currency struct {
		Currency string `json:"currency,omitempty"`
		Value    string `json:"value,omitempty"`
	}

	// CycleExecutions struct
	CycleExecutions struct {
		TenureType      string `json:"tenure_type,omitempty"`
		Sequence        int    `json:"sequence,omitempty"`
		CyclesCompleted int    `json:"cycles_completed,omitempty"`
		CyclesRemaining int    `json:"cycles_remaining,omitempty"`
		TotalCycles     int    `json:"total_cycles,omitempty"`
	}

	// LastPayment struct
	LastPayment struct {
		Amount Money     `json:"amount,omitempty"`
		Time   time.Time `json:"time,omitempty"`
	}

	// Details structure used in Amount structures as optional value
	Details struct {
		Subtotal         string `json:"subtotal,omitempty"`
		Shipping         string `json:"shipping,omitempty"`
		Tax              string `json:"tax,omitempty"`
		HandlingFee      string `json:"handling_fee,omitempty"`
		ShippingDiscount string `json:"shipping_discount,omitempty"`
		Insurance        string `json:"insurance,omitempty"`
		GiftWrap         string `json:"gift_wrap,omitempty"`
	}

	// ErrorResponseDetail struct
	ErrorResponseDetail struct {
		Field string `json:"field"`
		Issue string `json:"issue"`
		Links []Link `json:"link"`
	}

	// ErrorResponse https://developer.paypal.com/docs/api/errors/
	ErrorResponse struct {
		Response        *http.Response        `json:"-"`
		Name            string                `json:"name"`
		DebugID         string                `json:"debug_id"`
		Message         string                `json:"message"`
		InformationLink string                `json:"information_link"`
		Details         []ErrorResponseDetail `json:"details"`
	}

	// ExecuteAgreementResponse struct
	ExecuteAgreementResponse struct {
		ID               string           `json:"id"`
		State            string           `json:"state"`
		Description      string           `json:"description,omitempty"`
		Payer            Payer            `json:"payer"`
		Plan             BillingPlan      `json:"plan"`
		StartDate        time.Time        `json:"start_date"`
		ShippingAddress  ShippingAddress  `json:"shipping_address"`
		AgreementDetails AgreementDetails `json:"agreement_details"`
		Links            []Link           `json:"links"`
	}

	// ExecuteResponse struct
	ExecuteResponse struct {
		ID           string        `json:"id"`
		Links        []Link        `json:"links"`
		State        string        `json:"state"`
		Payer        PaymentPayer  `json:"payer"`
		Transactions []Transaction `json:"transactions,omitempty"`
	}

	// FundingInstrument struct
	FundingInstrument struct {
		CreditCard      *CreditCard      `json:"credit_card,omitempty"`
		CreditCardToken *CreditCardToken `json:"credit_card_token,omitempty"`
	}

	// Item struct
	Item struct {
		Name        string `json:"name"`
		UnitAmount  *Money `json:"unit_amount,omitempty"`
		Tax         *Money `json:"tax,omitempty"`
		Quantity    string `json:"quantity"`
		Description string `json:"description,omitempty"`
		SKU         string `json:"sku,omitempty"`
		Category    string `json:"category,omitempty"`
	}

	// ItemList struct
	ItemList struct {
		Items           []Item           `json:"items,omitempty"`
		ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
	}

	// Link struct
	Link struct {
		Href        string `json:"href"`
		Rel         string `json:"rel,omitempty"`
		Method      string `json:"method,omitempty"`
		Description string `json:"description,omitempty"`
		Enctype     string `json:"enctype,omitempty"`
	}

	// PurchaseUnitAmount struct
	PurchaseUnitAmount struct {
		Currency  string                       `json:"currency_code"`
		Value     string                       `json:"value"`
		Breakdown *PurchaseUnitAmountBreakdown `json:"breakdown,omitempty"`
	}

	// PurchaseUnitAmountBreakdown struct
	PurchaseUnitAmountBreakdown struct {
		ItemTotal        *Money `json:"item_total,omitempty"`
		Shipping         *Money `json:"shipping,omitempty"`
		Handling         *Money `json:"handling,omitempty"`
		TaxTotal         *Money `json:"tax_total,omitempty"`
		Insurance        *Money `json:"insurance,omitempty"`
		ShippingDiscount *Money `json:"shipping_discount,omitempty"`
		Discount         *Money `json:"discount,omitempty"`
	}

	// Money struct
	//
	// https://developer.paypal.com/docs/api/orders/v2/#definition-money
	Money struct {
		Currency string `json:"currency_code"`
		Value    string `json:"value"`
	}

	// PurchaseUnit struct
	PurchaseUnit struct {
		ReferenceID        string              `json:"reference_id"`
		Amount             *PurchaseUnitAmount `json:"amount,omitempty"`
		Payments           *CapturedPayments   `json:"payments,omitempty"`
		PaymentInstruction *PaymentInstruction `json:"payment_instruction,omitempty"`
	}

	// TaxInfo used for orders.
	TaxInfo struct {
		TaxID     string `json:"tax_id,omitempty"`
		TaxIDType string `json:"tax_id_type,omitempty"`
	}

	// PhoneWithTypeNumber struct for PhoneWithType
	PhoneWithTypeNumber struct {
		NationalNumber string `json:"national_number,omitempty"`
	}

	// PhoneWithType struct used for orders
	PhoneWithType struct {
		PhoneType   string               `json:"phone_type,omitempty"`
		PhoneNumber *PhoneWithTypeNumber `json:"phone_number,omitempty"`
	}

	// CreateOrderPayerName create order payer name
	CreateOrderPayerName struct {
		GivenName string `json:"given_name,omitempty"`
		Surname   string `json:"surname,omitempty"`
	}

	// CreateOrderPayer used with create order requests
	CreateOrderPayer struct {
		Name         *CreateOrderPayerName          `json:"name,omitempty"`
		EmailAddress string                         `json:"email_address,omitempty"`
		PayerID      string                         `json:"payer_id,omitempty"`
		Phone        *PhoneWithType                 `json:"phone,omitempty"`
		BirthDate    string                         `json:"birth_date,omitempty"`
		TaxInfo      *TaxInfo                       `json:"tax_info,omitempty"`
		Address      *ShippingDetailAddressPortable `json:"address,omitempty"`
	}

	// PurchaseUnitRequest struct
	PurchaseUnitRequest struct {
		ReferenceID        string              `json:"reference_id,omitempty"`
		Amount             *PurchaseUnitAmount `json:"amount"`
		Payee              *PayeeForOrders     `json:"payee,omitempty"`
		Description        string              `json:"description,omitempty"`
		CustomID           string              `json:"custom_id,omitempty"`
		InvoiceID          string              `json:"invoice_id,omitempty"`
		SoftDescriptor     string              `json:"soft_descriptor,omitempty"`
		Items              []Item              `json:"items,omitempty"`
		Shipping           *ShippingDetail     `json:"shipping,omitempty"`
		PaymentInstruction *PaymentInstruction `json:"payment_instruction,omitempty"`
	}

	// MerchantPreferences struct
	MerchantPreferences struct {
		SetupFee                *AmountPayout `json:"setup_fee,omitempty"`
		ReturnURL               string        `json:"return_url,omitempty"`
		CancelURL               string        `json:"cancel_url,omitempty"`
		AutoBillAmount          string        `json:"auto_bill_amount,omitempty"`
		InitialFailAmountAction string        `json:"initial_fail_amount_action,omitempty"`
		MaxFailAttempts         string        `json:"max_fail_attempts,omitempty"`
	}

	// Order struct
	Order struct {
		ID            string                 `json:"id,omitempty"`
		Status        string                 `json:"status,omitempty"`
		Intent        string                 `json:"intent,omitempty"`
		Payer         *PayerWithNameAndPhone `json:"payer,omitempty"`
		PurchaseUnits []PurchaseUnit         `json:"purchase_units,omitempty"`
		Links         []Link                 `json:"links,omitempty"`
		CreateTime    *time.Time             `json:"create_time,omitempty"`
		UpdateTime    *time.Time             `json:"update_time,omitempty"`
	}

	// ExchangeRate struct
	//
	// https://developer.paypal.com/docs/api/orders/v2/#definition-exchange_rate
	ExchangeRate struct {
		SourceCurrency string `json:"source_currency"`
		TargetCurrency string `json:"target_currency"`
		Value          string `json:"value"`
	}

	// SellerReceivableBreakdown has the detailed breakdown of the capture activity.
	SellerReceivableBreakdown struct {
		GrossAmount                   *Money        `json:"gross_amount,omitempty"`
		PaypalFee                     *Money        `json:"paypal_fee,omitempty"`
		PaypalFeeInReceivableCurrency *Money        `json:"paypal_fee_in_receivable_currency,omitempty"`
		NetAmount                     *Money        `json:"net_amount,omitempty"`
		ReceivableAmount              *Money        `json:"receivable_amount,omitempty"`
		ExchangeRate                  *ExchangeRate `json:"exchange_rate,omitempty"`
		PlatformFees                  []PlatformFee `json:"platform_fees,omitempty"`
	}

	// CaptureAmount struct
	CaptureAmount struct {
		ID                        string                     `json:"id,omitempty"`
		CustomID                  string                     `json:"custom_id,omitempty"`
		Amount                    *PurchaseUnitAmount        `json:"amount,omitempty"`
		SellerProtection          *SellerProtection          `json:"seller_protection,omitempty"`
		SellerReceivableBreakdown *SellerReceivableBreakdown `json:"seller_receivable_breakdown,omitempty"`
	}

	// CapturedPayments has the amounts for a captured order
	CapturedPayments struct {
		Captures []CaptureAmount `json:"captures,omitempty"`
	}

	// CapturedPurchaseItem are items for a captured order
	CapturedPurchaseItem struct {
		Quantity    string `json:"quantity"`
		Name        string `json:"name"`
		SKU         string `json:"sku,omitempty"`
		Description string `json:"description,omitempty"`
	}

	// CapturedPurchaseUnit are purchase units for a captured order
	CapturedPurchaseUnit struct {
		Items       []CapturedPurchaseItem       `json:"items,omitempty"`
		ReferenceID string                       `json:"reference_id"`
		Shipping    CapturedPurchaseUnitShipping `json:"shipping,omitempty"`
		Payments    *CapturedPayments            `json:"payments,omitempty"`
	}

	CapturedPurchaseUnitShipping struct {
		Address ShippingDetailAddressPortable `json:"address,omitempty"`
	}

	// PayerWithNameAndPhone struct
	PayerWithNameAndPhone struct {
		Name         *CreateOrderPayerName `json:"name,omitempty"`
		EmailAddress string                `json:"email_address,omitempty"`
		Phone        *PhoneWithType        `json:"phone,omitempty"`
		PayerID      string                `json:"payer_id,omitempty"`
		Address      Address               `json:"address,omitempty"`
	}

	// CaptureOrderResponse is the response for capture order
	CaptureOrderResponse struct {
		ID            string                 `json:"id,omitempty"`
		Status        string                 `json:"status,omitempty"`
		Payer         *PayerWithNameAndPhone `json:"payer,omitempty"`
		Address       *Address               `json:"address,omitempty"`
		PurchaseUnits []CapturedPurchaseUnit `json:"purchase_units,omitempty"`
	}

	// Payer struct
	Payer struct {
		PaymentMethod      string              `json:"payment_method"`
		FundingInstruments []FundingInstrument `json:"funding_instruments,omitempty"`
		PayerInfo          *PayerInfo          `json:"payer_info,omitempty"`
		Status             string              `json:"payer_status,omitempty"`
	}

	// PayerInfo struct
	PayerInfo struct {
		Email           string           `json:"email,omitempty"`
		FirstName       string           `json:"first_name,omitempty"`
		LastName        string           `json:"last_name,omitempty"`
		PayerID         string           `json:"payer_id,omitempty"`
		Phone           string           `json:"phone,omitempty"`
		ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
		TaxIDType       string           `json:"tax_id_type,omitempty"`
		TaxID           string           `json:"tax_id,omitempty"`
		CountryCode     string           `json:"country_code"`
	}

	// PaymentDefinition struct
	PaymentDefinition struct {
		ID                string        `json:"id,omitempty"`
		Name              string        `json:"name,omitempty"`
		Type              string        `json:"type,omitempty"`
		Frequency         string        `json:"frequency,omitempty"`
		FrequencyInterval string        `json:"frequency_interval,omitempty"`
		Amount            AmountPayout  `json:"amount,omitempty"`
		Cycles            string        `json:"cycles,omitempty"`
		ChargeModels      []ChargeModel `json:"charge_models,omitempty"`
	}

	// PaymentOptions struct
	PaymentOptions struct {
		AllowedPaymentMethod string `json:"allowed_payment_method,omitempty"`
	}

	// PaymentPatch PATCH /v2/payments/payment/{payment_id)
	PaymentPatch struct {
		Operation string      `json:"op"`
		Path      string      `json:"path"`
		Value     interface{} `json:"value"`
	}

	// PaymentPayer struct
	PaymentPayer struct {
		PaymentMethod string     `json:"payment_method"`
		Status        string     `json:"status,omitempty"`
		PayerInfo     *PayerInfo `json:"payer_info,omitempty"`
	}

	// PaymentResponse structure
	PaymentResponse struct {
		ID           string        `json:"id"`
		State        string        `json:"state"`
		Intent       string        `json:"intent"`
		Payer        Payer         `json:"payer"`
		Transactions []Transaction `json:"transactions"`
		Links        []Link        `json:"links"`
	}

	// PaymentSource structure
	PaymentSource struct {
		Card  *PaymentSourceCard  `json:"card"`
		Token *PaymentSourceToken `json:"token"`
	}

	// PaymentSourceCard structure
	PaymentSourceCard struct {
		ID             string              `json:"id"`
		Name           string              `json:"name"`
		Number         string              `json:"number"`
		Expiry         string              `json:"expiry"`
		SecurityCode   string              `json:"security_code"`
		LastDigits     string              `json:"last_digits"`
		CardType       string              `json:"card_type"`
		BillingAddress *CardBillingAddress `json:"billing_address"`
	}

	// CardBillingAddress structure
	CardBillingAddress struct {
		AddressLine1 string `json:"address_line_1"`
		AddressLine2 string `json:"address_line_2"`
		AdminArea2   string `json:"admin_area_2"`
		AdminArea1   string `json:"admin_area_1"`
		PostalCode   string `json:"postal_code"`
		CountryCode  string `json:"country_code"`
	}

	// PaymentSourceToken structure
	PaymentSourceToken struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}

	// Payout struct
	Payout struct {
		SenderBatchHeader *SenderBatchHeader `json:"sender_batch_header"`
		Items             []PayoutItem       `json:"items"`
	}

	// PayoutItem struct
	PayoutItem struct {
		RecipientType   string        `json:"recipient_type"`
		RecipientWallet string        `json:"recipient_wallet"`
		Receiver        string        `json:"receiver"`
		Amount          *AmountPayout `json:"amount"`
		Note            string        `json:"note,omitempty"`
		SenderItemID    string        `json:"sender_item_id,omitempty"`
	}

	// PayoutItemResponse struct
	PayoutItemResponse struct {
		PayoutItemID      string        `json:"payout_item_id"`
		TransactionID     string        `json:"transaction_id"`
		TransactionStatus string        `json:"transaction_status"`
		PayoutBatchID     string        `json:"payout_batch_id,omitempty"`
		PayoutItemFee     *AmountPayout `json:"payout_item_fee,omitempty"`
		PayoutItem        *PayoutItem   `json:"payout_item"`
		TimeProcessed     *time.Time    `json:"time_processed,omitempty"`
		Links             []Link        `json:"links"`
		Error             ErrorResponse `json:"errors,omitempty"`
	}

	// PayoutResponse struct
	PayoutResponse struct {
		BatchHeader *BatchHeader         `json:"batch_header"`
		Items       []PayoutItemResponse `json:"items"`
		Links       []Link               `json:"links"`
	}

	// RedirectURLs struct
	RedirectURLs struct {
		ReturnURL string `json:"return_url,omitempty"`
		CancelURL string `json:"cancel_url,omitempty"`
	}

	// Refund struct
	Refund struct {
		ID            string     `json:"id,omitempty"`
		Amount        *Amount    `json:"amount,omitempty"`
		CreateTime    *time.Time `json:"create_time,omitempty"`
		State         string     `json:"state,omitempty"`
		CaptureID     string     `json:"capture_id,omitempty"`
		ParentPayment string     `json:"parent_payment,omitempty"`
		UpdateTime    *time.Time `json:"update_time,omitempty"`
	}

	// RefundResponse .
	RefundResponse struct {
		ID     string              `json:"id,omitempty"`
		Amount *PurchaseUnitAmount `json:"amount,omitempty"`
		Status string              `json:"status,omitempty"`
	}

	// Related struct
	Related struct {
		Sale          *Sale          `json:"sale,omitempty"`
		Authorization *Authorization `json:"authorization,omitempty"`
		Order         *Order         `json:"order,omitempty"`
		Capture       *Capture       `json:"capture,omitempty"`
		Refund        *Refund        `json:"refund,omitempty"`
	}

	// Sale struct
	Sale struct {
		ID                        string     `json:"id,omitempty"`
		Amount                    *Amount    `json:"amount,omitempty"`
		TransactionFee            *Currency  `json:"transaction_fee,omitempty"`
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
		Links                     []Link     `json:"links,omitempty"`
	}

	// SenderBatchHeader struct
	SenderBatchHeader struct {
		EmailSubject  string `json:"email_subject"`
		EmailMessage  string `json:"email_message"`
		SenderBatchID string `json:"sender_batch_id,omitempty"`
	}

	//ShippingAmount struct
	ShippingAmount struct {
		Money
	}

	// ShippingAddress struct
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

	// ShippingDetailAddressPortable used with create orders
	ShippingDetailAddressPortable struct {
		AddressLine1 string `json:"address_line_1,omitempty"`
		AddressLine2 string `json:"address_line_2,omitempty"`
		AdminArea1   string `json:"admin_area_1,omitempty"`
		AdminArea2   string `json:"admin_area_2,omitempty"`
		PostalCode   string `json:"postal_code,omitempty"`
		CountryCode  string `json:"country_code,omitempty"`
	}

	// Name struct
	//Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#definition-name
	Name struct {
		FullName   string `json:"full_name,omitempty"`
		Suffix     string `json:"suffix,omitempty"`
		Prefix     string `json:"prefix,omitempty"`
		GivenName  string `json:"given_name,omitempty"`
		Surname    string `json:"surname,omitempty"`
		MiddleName string `json:"middle_name,omitempty"`
	}

	// ShippingDetail struct
	ShippingDetail struct {
		Name    *Name                          `json:"name,omitempty"`
		Address *ShippingDetailAddressPortable `json:"address,omitempty"`
	}

	// Subscriber struct
	Subscriber struct {
		ShippingAddress ShippingDetail       `json:"shipping_address,omitempty"`
		Name            CreateOrderPayerName `json:"name,omitempty"`
		EmailAddress    string               `json:"email_address,omitempty"`
	}

	expirationTime int64

	// TokenResponse is for API response for the /oauth2/token endpoint
	TokenResponse struct {
		RefreshToken string         `json:"refresh_token"`
		Token        string         `json:"access_token"`
		Type         string         `json:"token_type"`
		ExpiresIn    expirationTime `json:"expires_in"`
	}

	// Transaction struct
	Transaction struct {
		Amount           *Amount         `json:"amount"`
		Description      string          `json:"description,omitempty"`
		ItemList         *ItemList       `json:"item_list,omitempty"`
		InvoiceNumber    string          `json:"invoice_number,omitempty"`
		Custom           string          `json:"custom,omitempty"`
		SoftDescriptor   string          `json:"soft_descriptor,omitempty"`
		RelatedResources []Related       `json:"related_resources,omitempty"`
		PaymentOptions   *PaymentOptions `json:"payment_options,omitempty"`
		NotifyURL        string          `json:"notify_url,omitempty"`
		OrderURL         string          `json:"order_url,omitempty"`
		Payee            *Payee          `json:"payee,omitempty"`
	}

	//Payee struct
	Payee struct {
		Email string `json:"email"`
	}

	// PayeeForOrders struct
	PayeeForOrders struct {
		EmailAddress string `json:"email_address,omitempty"`
		MerchantID   string `json:"merchant_id,omitempty"`
	}

	// UserInfo struct
	UserInfo struct {
		ID              string   `json:"user_id"`
		Name            string   `json:"name"`
		GivenName       string   `json:"given_name"`
		FamilyName      string   `json:"family_name"`
		Email           string   `json:"email"`
		Verified        bool     `json:"verified,omitempty,string"`
		Gender          string   `json:"gender,omitempty"`
		BirthDate       string   `json:"birthdate,omitempty"`
		ZoneInfo        string   `json:"zoneinfo,omitempty"`
		Locale          string   `json:"locale,omitempty"`
		Phone           string   `json:"phone_number,omitempty"`
		Address         *Address `json:"address,omitempty"`
		VerifiedAccount bool     `json:"verified_account,omitempty,string"`
		AccountType     string   `json:"account_type,omitempty"`
		AgeRange        string   `json:"age_range,omitempty"`
		PayerID         string   `json:"payer_id,omitempty"`
	}

	// WebProfile represents the configuration of the payment web payment experience
	//
	// https://developer.paypal.com/docs/api/payment-experience/
	WebProfile struct {
		ID           string       `json:"id,omitempty"`
		Name         string       `json:"name"`
		Presentation Presentation `json:"presentation,omitempty"`
		InputFields  InputFields  `json:"input_fields,omitempty"`
		FlowConfig   FlowConfig   `json:"flow_config,omitempty"`
	}

	// Presentation represents the branding and locale that a customer sees on
	// redirect payments
	//
	// https://developer.paypal.com/docs/api/payment-experience/#definition-presentation
	Presentation struct {
		BrandName  string `json:"brand_name,omitempty"`
		LogoImage  string `json:"logo_image,omitempty"`
		LocaleCode string `json:"locale_code,omitempty"`
	}

	// InputFields represents the fields that are displayed to a customer on
	// redirect payments
	//
	// https://developer.paypal.com/docs/api/payment-experience/#definition-input_fields
	InputFields struct {
		AllowNote       bool `json:"allow_note,omitempty"`
		NoShipping      uint `json:"no_shipping,omitempty"`
		AddressOverride uint `json:"address_override,omitempty"`
	}

	// FlowConfig represents the general behaviour of redirect payment pages
	//
	// https://developer.paypal.com/docs/api/payment-experience/#definition-flow_config
	FlowConfig struct {
		LandingPageType   string `json:"landing_page_type,omitempty"`
		BankTXNPendingURL string `json:"bank_txn_pending_url,omitempty"`
		UserAction        string `json:"user_action,omitempty"`
	}

	// VerifyWebhookResponse struct
	VerifyWebhookResponse struct {
		VerificationStatus string `json:"verification_status,omitempty"`
	}

	WebhookEventTypesResponse struct {
		EventTypes []WebhookEventType `json:"event_types"`
	}

	// Webhook struct
	Webhook struct {
		ID         string             `json:"id"`
		URL        string             `json:"url"`
		EventTypes []WebhookEventType `json:"event_types"`
		Links      []Link             `json:"links"`
	}

	// Event struct.
	//
	// The basic webhook event data type. This struct is intended to be
	// embedded into resource type specific event structs.
	Event struct {
		ID              string    `json:"id"`
		CreateTime      time.Time `json:"create_time"`
		ResourceType    string    `json:"resource_type"`
		EventType       string    `json:"event_type"`
		Summary         string    `json:"summary,omitempty"`
		Links           []Link    `json:"links"`
		EventVersion    string    `json:"event_version,omitempty"`
		ResourceVersion string    `json:"resource_version,omitempty"`
	}

	AnyEvent struct {
		Event
		Resource json.RawMessage `json:"resource"`
	}

	// WebhookEventType struct
	WebhookEventType struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status,omitempty"`
	}

	// CreateWebhookRequest struct
	CreateWebhookRequest struct {
		URL        string             `json:"url"`
		EventTypes []WebhookEventType `json:"event_types"`
	}

	ListWebhookResponse struct {
		Webhooks []Webhook `json:"webhooks"`
	}

	WebhookField struct {
		Operation string      `json:"op"`
		Path      string      `json:"path"`
		Value     interface{} `json:"value"`
	}

	// Resource is a mix of fields from several webhook resource types.
	//
	// Deprecated: Add implementation of specific resource types in your own
	// code and don't use this catch all struct, you show know which resource
	// type you are expecting and handle that type only.
	//
	// Every resource struct type should be unique for every combination of
	// "resource_type"/"resource_version" combination of the Event type /
	// webhook message.
	Resource struct {
		ID                        string                     `json:"id,omitempty"`
		Status                    string                     `json:"status,omitempty"`
		StatusDetails             *CaptureStatusDetails      `json:"status_details,omitempty"`
		Amount                    *PurchaseUnitAmount        `json:"amount,omitempty"`
		UpdateTime                string                     `json:"update_time,omitempty"`
		CreateTime                string                     `json:"create_time,omitempty"`
		ExpirationTime            string                     `json:"expiration_time,omitempty"`
		SellerProtection          *SellerProtection          `json:"seller_protection,omitempty"`
		FinalCapture              bool                       `json:"final_capture,omitempty"`
		SellerPayableBreakdown    *CaptureSellerBreakdown    `json:"seller_payable_breakdown,omitempty"`
		SellerReceivableBreakdown *SellerReceivableBreakdown `json:"seller_receivable_breakdown,omitempty"`
		NoteToPayer               string                     `json:"note_to_payer,omitempty"`
		CustomID                  string                     `json:"custom_id,omitempty"`
		PartnerClientID           string                     `json:"partner_client_id,omitempty"`
		MerchantID                string                     `json:"merchant_id,omitempty"`
		Intent                    string                     `json:"intent,omitempty"`
		BillingAgreementID        *string                    `json:"billing_agreement_id,omitempty"`
		PurchaseUnits             []*PurchaseUnitRequest     `json:"purchase_units,omitempty"`
		Payer                     *PayerWithNameAndPhone     `json:"payer,omitempty"`
		Links                     []Link                     `json:"links,omitempty"`
	}

	CaptureSellerBreakdown struct {
		GrossAmount         PurchaseUnitAmount  `json:"gross_amount"`
		PayPalFee           PurchaseUnitAmount  `json:"paypal_fee"`
		NetAmount           PurchaseUnitAmount  `json:"net_amount"`
		TotalRefundedAmount *PurchaseUnitAmount `json:"total_refunded_amount,omitempty"`
	}

	ReferralRequest struct {
		TrackingID            string                 `json:"tracking_id"`
		PartnerConfigOverride *PartnerConfigOverride `json:"partner_config_override,omitempty"`
		Operations            []Operation            `json:"operations,omitempty"`
		Products              []string               `json:"products,omitempty"`
		LegalConsents         []Consent              `json:"legal_consents,omitempty"`
	}

	PartnerConfigOverride struct {
		PartnerLogoURL       string `json:"partner_logo_url,omitempty"`
		ReturnURL            string `json:"return_url,omitempty"`
		ReturnURLDescription string `json:"return_url_description,omitempty"`
		ActionRenewalURL     string `json:"action_renewal_url,omitempty"`
		ShowAddCreditCard    *bool  `json:"show_add_credit_card,omitempty"`
	}

	Operation struct {
		Operation                string              `json:"operation"`
		APIIntegrationPreference *IntegrationDetails `json:"api_integration_preference,omitempty"`
	}

	IntegrationDetails struct {
		RestAPIIntegration *RestAPIIntegration `json:"rest_api_integration,omitempty"`
	}

	RestAPIIntegration struct {
		IntegrationMethod string            `json:"integration_method"`
		IntegrationType   string            `json:"integration_type"`
		ThirdPartyDetails ThirdPartyDetails `json:"third_party_details"`
	}

	ThirdPartyDetails struct {
		Features []string `json:"features"`
	}

	Consent struct {
		Type    string `json:"type"`
		Granted bool   `json:"granted"`
	}

	SearchItemDetails struct {
		ItemCode            string                 `json:"item_code"`
		ItemName            string                 `json:"item_name"`
		ItemDescription     string                 `json:"item_description"`
		ItemOptions         string                 `json:"item_options"`
		ItemQuantity        string                 `json:"item_quantity"`
		ItemUnitPrice       Money                  `json:"item_unit_price"`
		ItemAmount          Money                  `json:"item_amount"`
		DiscountAmount      *Money                 `json:"discount_amount"`
		AdjustmentAmount    *Money                 `json:"adjustment_amount"`
		GiftWrapAmount      *Money                 `json:"gift_wrap_amount"`
		TaxPercentage       string                 `json:"tax_percentage"`
		TaxAmounts          []SearchTaxAmount      `json:"tax_amounts"`
		BasicShippingAmount *Money                 `json:"basic_shipping_amount"`
		ExtraShippingAmount *Money                 `json:"extra_shipping_amount"`
		HandlingAmount      *Money                 `json:"handling_amount"`
		InsuranceAmount     *Money                 `json:"insurance_amount"`
		TotalItemAmount     Money                  `json:"total_item_amount"`
		InvoiceNumber       string                 `json:"invoice_number"`
		CheckoutOptions     []SearchCheckoutOption `json:"checkout_options"`
	}

	SearchCheckoutOption struct {
		CheckoutOptionName  string `json:"checkout_option_name"`
		CheckoutOptionValue string `json:"checkout_option_value"`
	}

	SearchCartInfo struct {
		ItemDetails     []SearchItemDetails `json:"item_details"`
		TaxInclusive    *bool               `json:"tax_inclusive"`
		PayPalInvoiceID string              `json:"paypal_invoice_id"`
	}

	SearchShippingInfo struct {
		Name                     string   `json:"name"`
		Method                   string   `json:"method"`
		Address                  Address  `json:"address"`
		SecondaryShippingAddress *Address `json:"secondary_shipping_address"`
	}

	SearchPayerName struct {
		GivenName string `json:"given_name"`
		Surname   string `json:"surname"`
	}

	SearchPayerInfo struct {
		AccountID     string               `json:"account_id"`
		EmailAddress  string               `json:"email_address"`
		PhoneNumber   *PhoneWithTypeNumber `json:"phone_number"`
		AddressStatus string               `json:"address_status"`
		PayerStatus   string               `json:"payer_status"`
		PayerName     SearchPayerName      `json:"payer_name"`
		CountryCode   string               `json:"country_code"`
		Address       *Address             `json:"address"`
	}

	SearchTaxAmount struct {
		TaxAmount Money `json:"tax_amount"`
	}

	SearchTransactionInfo struct {
		PayPalAccountID           string   `json:"paypal_account_id"`
		TransactionID             string   `json:"transaction_id"`
		PayPalReferenceID         string   `json:"paypal_reference_id"`
		PayPalReferenceIDType     string   `json:"paypal_reference_id_type"`
		TransactionEventCode      string   `json:"transaction_event_code"`
		TransactionInitiationDate JSONTime `json:"transaction_initiation_date"`
		TransactionUpdatedDate    JSONTime `json:"transaction_updated_date"`
		TransactionAmount         Money    `json:"transaction_amount"`
		FeeAmount                 *Money   `json:"fee_amount"`
		InsuranceAmount           *Money   `json:"insurance_amount"`
		ShippingAmount            *Money   `json:"shipping_amount"`
		ShippingDiscountAmount    *Money   `json:"shipping_discount_amount"`
		ShippingTaxAmount         *Money   `json:"shipping_tax_amount"`
		OtherAmount               *Money   `json:"other_amount"`
		TipAmount                 *Money   `json:"tip_amount"`
		TransactionStatus         string   `json:"transaction_status"`
		TransactionSubject        string   `json:"transaction_subject"`
		PaymentTrackingID         string   `json:"payment_tracking_id"`
		BankReferenceID           string   `json:"bank_reference_id"`
		TransactionNote           string   `json:"transaction_note"`
		EndingBalance             *Money   `json:"ending_balance"`
		AvailableBalance          *Money   `json:"available_balance"`
		InvoiceID                 string   `json:"invoice_id"`
		CustomField               string   `json:"custom_field"`
		ProtectionEligibility     string   `json:"protection_eligibility"`
		CreditTerm                string   `json:"credit_term"`
		CreditTransactionalFee    *Money   `json:"credit_transactional_fee"`
		CreditPromotionalFee      *Money   `json:"credit_promotional_fee"`
		AnnualPercentageRate      string   `json:"annual_percentage_rate"`
		PaymentMethodType         string   `json:"payment_method_type"`
	}

	SearchTransactionDetails struct {
		TransactionInfo SearchTransactionInfo `json:"transaction_info"`
		PayerInfo       *SearchPayerInfo      `json:"payer_info"`
		ShippingInfo    *SearchShippingInfo   `json:"shipping_info"`
		CartInfo        *SearchCartInfo       `json:"cart_info"`
	}

	SharedResponse struct {
		CreateTime string `json:"create_time"`
		UpdateTime string `json:"update_time"`
		Links      []Link `json:"links"`
	}

	ListParams struct {
		Page          string `json:"page,omitempty"`           //Default: 0.
		PageSize      string `json:"page_size,omitempty"`      //Default: 10.
		TotalRequired string `json:"total_required,omitempty"` //Default: no.
	}

	SharedListResponse struct {
		TotalItems int    `json:"total_items,omitempty"`
		TotalPages int    `json:"total_pages,omitempty"`
		Links      []Link `json:"links,omitempty"`
	}
)

// Error method implementation for ErrorResponse struct
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %s, %+v", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, r.Details)
}

// MarshalJSON for JSONTime
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf(`"%s"`, time.Time(t).UTC().Format(time.RFC3339))
	return []byte(stamp), nil
}

// UnmarshalJSON for JSONTime, timezone offset is missing a colon ':"
func (t *JSONTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse("2006-01-02T15:04:05Z0700", s)
	*t = JSONTime(nt)
	return err
}

func (e *expirationTime) UnmarshalJSON(b []byte) error {
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	i, err := n.Int64()
	if err != nil {
		return err
	}
	*e = expirationTime(i)
	return nil
}
