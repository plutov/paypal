package paypal

type SubscriptionPlanStatus string

const (
	SubscriptionPlanStatusCreated  SubscriptionPlanStatus = "CREATED"
	SubscriptionPlanStatusInactive SubscriptionPlanStatus = "INACTIVE"
	SubscriptionPlanStatusActive   SubscriptionPlanStatus = "ACTIVE"
)

type BillingPlanStatus string

const (
	BillingPlanStatusActive   BillingPlanStatus = "ACTIVE"
)

type IntervalUnit string

const (
	IntervalUnitDay   IntervalUnit = "DAY"
	IntervalUnitWeek  IntervalUnit = "WEEK"
	IntervalUnitMonth IntervalUnit = "MONTH"
	IntervalUnitYear  IntervalUnit = "YEAR"
)

type TenureType string

const (
	TenureTypeRegular TenureType = "REGULAR"
	TenureTypeTrial   TenureType = "TRIAL"
)

type SetupFeeFailureAction string

const (
	SetupFeeFailureActionContinue SetupFeeFailureAction = "CONTINUE"
	SetupFeeFailureActionCancel   SetupFeeFailureAction = "CANCEL"
)

type ShippingPreference string

const (
	ShippingPreferenceGetFromFile        ShippingPreference = "GET_FROM_FILE"
	ShippingPreferenceNoShipping         ShippingPreference = "NO_SHIPPING"
	ShippingPreferenceSetProvidedAddress ShippingPreference = "SET_PROVIDED_ADDRESS"
)

type UserAction string

const (
	UserActionContinue     UserAction = "CONTINUE"
	UserActionSubscribeNow UserAction = "SUBSCRIBE_NOW"
)

type SubscriptionStatus string

const (
	SubscriptionStatusApprovalPending SubscriptionStatus = "APPROVAL_PENDING"
	SubscriptionStatusApproved        SubscriptionStatus = "APPROVED"
	SubscriptionStatusActive          SubscriptionStatus = "ACTIVE"
	SubscriptionStatusSuspended       SubscriptionStatus = "SUSPENDED"
	SubscriptionStatusCancelled       SubscriptionStatus = "CANCELLED"
	SubscriptionStatusExpired         SubscriptionStatus = "EXPIRED"
)

//Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#definition-transaction
type SubscriptionTransactionStatus string
const (
	SubscriptionCaptureStatusCompleted         SubscriptionTransactionStatus = "COMPLETED"
	SubscriptionCaptureStatusDeclined          SubscriptionTransactionStatus = "DECLINED"
	SubscriptionCaptureStatusPartiallyRefunded SubscriptionTransactionStatus = "PARTIALLY_REFUNDED"
	SubscriptionCaptureStatusPending           SubscriptionTransactionStatus = "PENDING"
	SubscriptionCaptureStatusRefunded          SubscriptionTransactionStatus = "REFUNDED"
)

type CaptureType string
const (
	CaptureTypeOutstandingBalance CaptureType = "OUTSTANDING_BALANCE"
)