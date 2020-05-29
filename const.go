package paypal

type SubscriptionPlanStatus string

const (
	SUBSCRIPTION_PLAN_STATUS_CREATED  SubscriptionPlanStatus = "CREATED"
	SUBSCRIPTION_PLAN_STATUS_INACTIVE SubscriptionPlanStatus = "INACTIVE"
	SUBSCRIPTION_PLAN_STATUS_ACTIVE   SubscriptionPlanStatus = "ACTIVE"
)

type IntervalUnit string

const (
	INTERVAL_UNIT_DAY   IntervalUnit = "DAY"
	INTERVAL_UNIT_WEEK  IntervalUnit = "WEEK"
	INTERVAL_UNIT_MONTH IntervalUnit = "MONTH"
	INTERVAL_UNIT_YEAR  IntervalUnit = "YEAR"
)

type TenureType string

const (
	TENURE_TYPE_REGULAR TenureType = "REGULAR"
	TENURE_TYPE_TRIAL   TenureType = "TRIAL"
)

type SetupFeeFailureAction string

const (
	SETUP_FEE_FAILURE_ACTION_CONTINUE SetupFeeFailureAction = "CONTINUE"
	SETUP_FEE_FAILURE_ACTION_CANCEL   SetupFeeFailureAction = "CANCEL"
)

type ShippingPreference string

const (
	ShippingPreference_GET_FROM_FILE        ShippingPreference = "GET_FROM_FILE"
	ShippingPreference_NO_SHIPPING          ShippingPreference = "NO_SHIPPING"
	ShippingPreference_SET_PROVIDED_ADDRESS ShippingPreference = "SET_PROVIDED_ADDRESS"
)

type UserAction string

const (
	UserAction_CONTINUE      UserAction = "CONTINUE"
	UserAction_SUBSCRIBE_NOW UserAction = "SUBSCRIBE_NOW"
)

type SubscriptionStatus string

const (
	SubscriptionStatus_APPROVAL_PENDING SubscriptionStatus = "APPROVAL_PENDING"
	SubscriptionStatus_APPROVED         SubscriptionStatus = "APPROVED"
	SubscriptionStatus_ACTIVE           SubscriptionStatus = "ACTIVE"
	SubscriptionStatus_SUSPENDED        SubscriptionStatus = "SUSPENDED"
	SubscriptionStatus_CANCELLED        SubscriptionStatus = "CANCELLED"
	SubscriptionStatus_EXPIRED          SubscriptionStatus = "EXPIRED"
)

//Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#definition-transaction
type SubscriptionTransactionStatus string
const (
	SubscriptionCaptureStatus_COMPLETED          SubscriptionTransactionStatus = "COMPLETED"
	SubscriptionCaptureStatus_DECLINED           SubscriptionTransactionStatus = "DECLINED"
	SubscriptionCaptureStatus_PARTIALLY_REFUNDED SubscriptionTransactionStatus = "PARTIALLY_REFUNDED"
	SubscriptionCaptureStatus_PENDING            SubscriptionTransactionStatus = "PENDING"
	SubscriptionCaptureStatus_REFUNDED           SubscriptionTransactionStatus = "REFUNDED"
)

type CaptureType string
const (
	CaptureType_OUTSTANDING_BALANCE CaptureType = "OUTSTANDING_BALANCE"
)