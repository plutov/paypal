package paypal

import "time"

type ListDisputesRequest struct {
	StartTime             *time.Time
	DisputedTransactionId *string
	PageSize              *int
	NextPageToken         *string
	DisputeState          *string
	UpdateTimeBefore      *time.Time
	UpdateTimeAfter       *time.Time
}

type ListDisputesResponse struct {
	Items []*DisputeItem `json:"items,omitempty"`
	Links []Link         `json:"links,omitempty"`
}

type DisputeItem struct {
	DisputeID             string                     `json:"dispute_id,omitempty"`
	CreateTime            time.Time                  `json:"create_time"`
	UpdateTime            time.Time                  `json:"update_time"`
	DisputedTransactions  []*DisputedTransactionItem `json:"disputed_transactions,omitempty"`
	Reason                string                     `json:"reason,omitempty"`
	Status                string                     `json:"status,omitempty"`
	DisputeState          DisputeState               `json:"dispute_state,omitempty"`
	DisputeAmount         *Money                     `json:"dispute_amount,omitempty"`
	Outcome               string                     `json:"outcome,omitempty"`
	DisputeLifeCycleStage string                     `json:"dispute_life_cycle_stage,omitempty"`
	DisputeChannel        string                     `json:"dispute_channel,omitempty"`
	DisputeFlow           string                     `json:"dispute_flow,omitempty"`
	Links                 []Link                     `json:"links,omitempty"`
}

type DisputedTransactionItem struct {
	BuyerTransactionID string `json:"buyer_transaction_id,omitempty"`
	Buyer              Buyer  `json:"buyer"`
	Seller             Seller `json:"seller"`
}

type Buyer struct {
	PayerID string `json:"payer_id,omitempty"`
	Name    string `json:"name,omitempty"`
}

type Seller struct {
	Email      string `json:"email,omitempty"`
	MerchantID string `json:"merchant_id,omitempty"`
	Name       string `json:"name,omitempty"`
}

type GetDisputeDetailResponse struct {
	DisputeID             string                       `json:"dispute_id,omitempty"`
	CreateTime            time.Time                    `json:"create_time"`
	UpdateTime            time.Time                    `json:"update_time"`
	DisputedTransactions  []*DisputedTransactionDetail `json:"disputed_transactions,omitempty"`
	Reason                string                       `json:"reason,omitempty"`
	Status                string                       `json:"status,omitempty"`
	DisputeAmount         *Money                       `json:"dispute_amount,omitempty"`
	DisputeOutcome        *DisputeOutcome              `json:"dispute_outcome,omitempty"`
	DisputeLifeCycleStage string                       `json:"dispute_life_cycle_stage,omitempty"`
	DisputeChannel        string                       `json:"dispute_channel,omitempty"`
	Messages              []*Message                   `json:"messages,omitempty"`
	Extensions            *Extensions                  `json:"extensions,omitempty"`
	Offer                 *Offer                       `json:"offer,omitempty"`
	Links                 []Link                       `json:"links,omitempty"`
}

type DisputeOutcome struct {
	OutcomeCode    string `json:"outcome_code,omitempty"`
	AmountRefunded *Money `json:"amount_refunded,omitempty"`
}

type DisputedTransactionDetail struct {
	SellerTransactionID string    `json:"seller_transaction_id,omitempty"`
	CreateTime          time.Time `json:"create_time"`
	TransactionStatus   string    `json:"transaction_status,omitempty"`
	GrossAmount         *Money    `json:"gross_amount,omitempty"`
	Buyer               Buyer     `json:"buyer"`
	Seller              Seller    `json:"seller"`
}

type Extensions struct {
	MerchandizeDisputeProperties MerchandizeDisputeProperties `json:"merchandize_dispute_properties"`
}

type MerchandizeDisputeProperties struct {
	IssueType      string          `json:"issue_type,omitempty"`
	ServiceDetails *ServiceDetails `json:"service_details,omitempty"`
}

type ServiceDetails struct {
	SubReasons  []string `json:"sub_reasons,omitempty"`
	PurchaseURL string   `json:"purchase_url,omitempty"`
}

type Message struct {
	PostedBy   string      `json:"posted_by,omitempty"`
	TimePosted time.Time   `json:"time_posted"`
	Content    string      `json:"content,omitempty"`
	Documents  []*Document `json:"documents,omitempty"`
}

type Document struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type Offer struct {
	BuyerRequestedAmount *Money     `json:"buyer_requested_amount,omitempty"`
	OfferType            string     `json:"offer_type,omitempty"`
	History              []*History `json:"history,omitempty"`
}

type History struct {
	OfferTime             time.Time `json:"offer_time"`
	Actor                 string    `json:"actor,omitempty"`
	EventType             string    `json:"event_type,omitempty"`
	OfferType             string    `json:"offer_type,omitempty"`
	OfferAmount           *Money    `json:"offer_amount,omitempty"`
	Notes                 string    `json:"notes,omitempty"`
	DisputeLifeCycleStage string    `json:"dispute_life_cycle_stage,omitempty"`
}

type UpdateDisputeValue struct {
	ID         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	CreateTime time.Time `json:"create_time"`
	Reason     string    `json:"reason,omitempty"`
	Status     string    `json:"status,omitempty"`
}

type UpdateDisputeParams struct {
	Op    string              `json:"op,omitempty"`
	Path  string              `json:"path,omitempty"`
	Value *UpdateDisputeValue `json:"value,omitempty"`
}

type DisputeAcceptClaimParams struct {
	Note                  string                 `json:"note,omitempty"`
	AcceptClaimReason     AcceptClaimReason      `json:"accept_claim_reason,omitempty"`
	InvoiceId             string                 `json:"invoice_id,omitempty"`
	RefundAmount          *Money                 `json:"refund_amount,omitempty"`
	ReturnShippingAddress *ReturnShippingAddress `json:"return_shipping_address,omitempty"`
	ReturnShipmentInfo    []*ReturnShipmentInfo  `json:"return_shipment_info,omitempty"`
	AcceptClaimType       AcceptClaimType        `json:"accept_claim_type,omitempty"`
}

type ReturnShipmentInfo struct {
	ShipmentLabel *ShipmentLabel `json:"shipment_label"`
	TrackingInfo  *TrackingInfo  `json:"tracking_info"`
}

type ShipmentLabel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TrackingInfo struct {
	CarrierName      string `json:"carrier_name,omitempty"`
	TrackingNumber   string `json:"tracking_number,omitempty"`
	CarrierNameOther string `json:"carrier_name_other,omitempty"`
	TrackingUrl      string `json:"tracking_url,omitempty"`
}

type ReturnShippingAddress struct {
	AddressLine1 string `json:"address_line_1,omitempty"`
	AddressLine2 string `json:"address_line_2,omitempty"`
	CountryCode  string `json:"country_code,omitempty"`
	AdminArea1   string `json:"admin_area_1,omitempty"`
	AdminArea2   string `json:"admin_area_2,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
}

type DisputeMakeOfferParams struct {
	Note                  string                 `json:"note,omitempty"`
	InvoiceId             string                 `json:"invoice_id,omitempty"`
	OfferAmount           *Money                 `json:"offer_amount,omitempty"`
	ReturnShippingAddress *ReturnShippingAddress `json:"return_shipping_address,omitempty"`
	OfferType             MakeOfferType          `json:"offer_type,omitempty"`
}

type DisputeAcknowledgeReturnItemParams struct {
	Note      string             `json:"note,omitempty"`
	Evidences []*DisputeEvidence `json:"evidences,omitempty"`
}

type DisputeEvidence struct {
	EvidenceType EvidenceType         `json:"evidence_type,omitempty"`
	Documents    []*Document          `json:"documents,omitempty"`
	Notes        string               `json:"notes,omitempty"`
	ItemId       string               `json:"item_id,omitempty"`
	EvidenceInfo *DisputeEvidenceInfo `json:"evidence_info,omitempty"`
}

type DisputeEvidenceInfo struct {
	TrackingInfo []*TrackingInfo `json:"tracking_info,omitempty"`
	RefundIds    []string        `json:"refund_ids,omitempty"`
}

type DisputeProvideEvidenceParams struct {
	Evidences             *DisputeEvidence       `json:"evidences,omitempty"`
	ReturnShippingAddress *ReturnShippingAddress `json:"return_shipping_address,omitempty"`
}

type DisputeAppealParams struct {
	Evidences             *DisputeEvidence       `json:"evidences,omitempty"`
	ReturnShippingAddress *ReturnShippingAddress `json:"return_shipping_address,omitempty"`
}

// https://developer.paypal.com/docs/api/customer-disputes/v1/#definition-dispute_info
type DisputeState string

const (
	DisputeStateOpenInquiries            DisputeState = "OPEN_INQUIRIES"
	DisputeStateRequiredAction           DisputeState = "REQUIRED_ACTION"
	DisputeStateRequiredOtherPartyAction DisputeState = "REQUIRED_OTHER_PARTY_ACTION"
	DisputeStateUnderPaypalReview        DisputeState = "UNDER_PAYPAL_REVIEW"
	DisputeStateAppealable               DisputeState = "APPEALABLE"
	DisputeStateResolved                 DisputeState = "RESOLVED"
)

type AcceptClaimReason string

const (
	AcceptClaimReasonDidNotShipItem   AcceptClaimReason = "DID_NOT_SHIP_ITEM"
	AcceptClaimReasonTooTimeConsuming AcceptClaimReason = "TOO_TIME_CONSUMING"
	AcceptClaimReasonLostInMail       AcceptClaimReason = "LOST_IN_MAIL"
	AcceptClaimReasonNotAbleToWin     AcceptClaimReason = "NOT_ABLE_TO_WIN"
	AcceptClaimReasonCompanyPolicy    AcceptClaimReason = "COMPANY_POLICY"
	AcceptClaimReasonReasonNotSet     AcceptClaimReason = "REASON_NOT_SET"
)

type AcceptClaimType string

const (
	AcceptClaimTypeRefund                        AcceptClaimType = "REFUND"
	AcceptClaimTypeRefundWithReturn              AcceptClaimType = "REFUND_WITH_RETURN"
	AcceptClaimTypePartialRefund                 AcceptClaimType = "PARTIAL_REFUND"
	AcceptClaimTypeRefundWithReturnShipmentLabel AcceptClaimType = "REFUND_WITH_RETURN_SHIPMENT_LABEL"
)

type AdjudicationOutcome string

const (
	AdjudicationOutcomeBuyerFavor  AdjudicationOutcome = "BUYER_FAVOR"
	AdjudicationOutcomeSellerFavor AdjudicationOutcome = "SELLER_FAVOR"
)

type DisputeStatusAction string

const (
	DisputeStatusBuyerEvidence  DisputeStatusAction = "BUYER_EVIDENCE"
	DisputeStatusSellerEvidence DisputeStatusAction = "SELLER_EVIDENCE"
)

type DisputeStatus string

const (
	DisputeStatusOpen                     DisputeStatus = "OPEN"
	DisputeStatusWaitingForBuyerResponse  DisputeStatus = "WAITING_FOR_BUYER_RESPONSE"
	DisputeStatusWaitingForSellerResponse DisputeStatus = "WAITING_FOR_SELLER_RESPONSE"
	DisputeStatusUnderReview              DisputeStatus = "UNDER_REVIEW"
	DisputeStatusResolved                 DisputeStatus = "RESOLVED"
	DisputeStatusOther                    DisputeStatus = "OTHER"
)

type MakeOfferType string

const (
	MakeOfferTypeRefund                   MakeOfferType = "REFUND"
	MakeOfferTypeRefundWithReturn         MakeOfferType = "REFUND_WITH_RETURN"
	MakeOfferTypeRefundWithReplacement    MakeOfferType = "REFUND_WITH_REPLACEMENT"
	MakeOfferTypeReplacementWithoutRefund MakeOfferType = "REPLACEMENT_WITHOUT_REFUND"
)

type AcknowledgementType string

const (
	AcknowledgementTypeItemReceived            AcknowledgementType = "ITEM_RECEIVED"
	AcknowledgementTypeItemNotReceived         AcknowledgementType = "ITEM_NOT_RECEIVED"
	AcknowledgementTypeDamaged                 AcknowledgementType = "DAMAGED"
	AcknowledgementTypeEmptyPackageOrDifferent AcknowledgementType = "EMPTY_PACKAGE_OR_DIFFERENT"
	AcknowledgementTypeMissingItems            AcknowledgementType = "MISSING_ITEMS"
)

type EvidenceType string

const (
	EvidenceTypeProofOfDamage                                   EvidenceType = "PROOF_OF_DAMAGE"
	EvidenceTypeThirdpartyProofForDamageOrSignificantDifference EvidenceType = "THIRDPARTY_PROOF_FOR_DAMAGE_OR_SIGNIFICANT_DIFFERENCE"
	EvidenceTypeDelcaration                                     EvidenceType = "DECLARATION"
	EvidenceTypeProofOfMissingItems                             EvidenceType = "PROOF_OF_MISSING_ITEMS"
	EvidenceTypeProofOfEmptyPackageOrDifferentItem              EvidenceType = "PROOF_OF_EMPTY_PACKAGE_OR_DIFFERENT_ITEM"
	EvidenceTypeProofOfItemNotReceived                          EvidenceType = "PROOF_OF_ITEM_NOT_RECEIVED"
	EvidenceTypeProofOfFulfillment                              EvidenceType = "PROOF_OF_FULFILLMENT"
	EvidenceTypeProofOfRefund                                   EvidenceType = "PROOF_OF_REFUND"
	EvidenceTypeProofOfDeliverySignature                        EvidenceType = "PROOF_OF_DELIVERY_SIGNATURE"
	EvidenceTypeProofOfReceiptCopy                              EvidenceType = "PROOF_OF_RECEIPT_COPY"
	EvidenceTypeReturnPolicy                                    EvidenceType = "RETURN_POLICY"
	EvidenceTypeBillingAgreement                                EvidenceType = "BILLING_AGREEMENT"
	EvidenceTypeProofOfReshipment                               EvidenceType = "PROOF_OF_RESHIPMENT"
	EvidenceTypeItemDescription                                 EvidenceType = "ITEM_DESCRIPTION"
	EvidenceTypePoliceReport                                    EvidenceType = "POLICE_REPORT"
	EvidenceTypeAffidavit                                       EvidenceType = "AFFIDAVIT"
	EvidenceTypePaidWithOtherMethod                             EvidenceType = "PAID_WITH_OTHER_METHOD"
	EvidenceTypeCopyOfContract                                  EvidenceType = "COPY_OF_CONTRACT"
	EvidenceTypeTerminalAtmReceipt                              EvidenceType = "TERMINAL_ATM_RECEIPT"
	EvidenceTypePriceDifferenceReason                           EvidenceType = "PRICE_DIFFERENCE_REASON"
	EvidenceTypeSourceConversionRate                            EvidenceType = "SOURCE_CONVERSION_RATE"
	EvidenceTypeBankStatement                                   EvidenceType = "BANK_STATEMENT"
	EvidenceTypeCreditDueReason                                 EvidenceType = "CREDIT_DUE_REASON"
	EvidenceTypeRequestCreditReceipt                            EvidenceType = "REQUEST_CREDIT_RECEIPT"
	EvidenceTypeProofOfReturn                                   EvidenceType = "PROOF_OF_RETURN"
	EvidenceTypeCreate                                          EvidenceType = "CREATE"
	EvidenceTypeChangeReason                                    EvidenceType = "CHANGE_REASON"
	EvidenceTypeProofOfRefundOutsidePaypa                       EvidenceType = "PROOF_OF_REFUND_OUTSIDE_PAYPA"
	EvidenceTypeReceiptOfMerchandise                            EvidenceType = "RECEIPT_OF_MERCHANDISE"
	EvidenceTypeCustomsDocument                                 EvidenceType = "CUSTOMS_DOCUMENT"
	EvidenceTypeCustomsFeeReceipt                               EvidenceType = "CUSTOMS_FEE_RECEIPT"
	EvidenceTypeInformationOnResolution                         EvidenceType = "INFORMATION_ON_RESOLUTION"
	EvidenceTypeAdditionalInformationOfItem                     EvidenceType = "ADDITIONAL_INFORMATION_OF_ITEM"
	EvidenceTypeDetailsOfPurchase                               EvidenceType = "DETAILS_OF_PURCHASE"
	EvidenceTypeProofOfSignificantDifference                    EvidenceType = "PROOF_OF_SIGNIFICANT_DIFFERENCE"
	EvidenceTypeProofOfSoftwareOrServiceNotAsDescribed          EvidenceType = "PROOF_OF_SOFTWARE_OR_SERVICE_NOT_AS_DESCRIBED"
	EvidenceTypeProofOfConfiscation                             EvidenceType = "PROOF_OF_CONFISCATION	Documentation"
	EvidenceTypeCopyOfLawEnforcementAgencyReport                EvidenceType = "COPY_OF_LAW_ENFORCEMENT_AGENCY_REPORT"
	EvidenceTypeAdditionalProofOfShipment                       EvidenceType = "ADDITIONAL_PROOF_OF_SHIPMENT"
	EvidenceTypeProofOfDenialByCarrier                          EvidenceType = "PROOF_OF_DENIAL_BY_CARRIER"
	EvidenceTypeValidSupportingDocument                         EvidenceType = "VALID_SUPPORTING_DOCUMENT"
	EvidenceTypeLegibleSupportingDocument                       EvidenceType = "LEGIBLE_SUPPORTING_DOCUMENT"
	EvidenceTypeReturnTrackingInformation                       EvidenceType = "RETURN_TRACKING_INFORMATION"
	EvidenceTypeDeliveryReceipt                                 EvidenceType = "DELIVERY_RECEIPT"
	EvidenceTypeProofOfInstoreReceipt                           EvidenceType = "PROOF_OF_INSTORE_RECEIPT"
	EvidenceTypeAdditionalTrackingInformation                   EvidenceType = "ADDITIONAL_TRACKING_INFORMATION"
	EvidenceTypeProofOfShipmentPostage                          EvidenceType = "PROOF_OF_SHIPMENT_POSTAGE"
	EvidenceTypeOnlineTrackingInformation                       EvidenceType = "ONLINE_TRACKING_INFORMATION"
	EvidenceTypeProofOfInstoreRefund                            EvidenceType = "PROOF_OF_INSTORE_REFUND"
	EvidenceTypeProofForSoftwareOrServiceDelivered              EvidenceType = "PROOF_FOR_SOFTWARE_OR_SERVICE_DELIVERED"
	EvidenceTypeReturnAddressForShipping                        EvidenceType = "RETURN_ADDRESS_FOR_SHIPPING"
	EvidenceTypeCopyOfTheEparcelManifest                        EvidenceType = "COPY_OF_THE_EPARCEL_MANIFEST"
	EvidenceTypeCopyOfShippingManifest                          EvidenceType = "COPY_OF_SHIPPING_MANIFEST"
	EvidenceTypeAppealAffidavit                                 EvidenceType = "APPEAL_AFFIDAVIT"
	EvidenceTypeReceiptOfReplacement                            EvidenceType = "RECEIPT_OF_REPLACEMENT"
	EvidenceTypeCopyOfDriversLicense                            EvidenceType = "COPY_OF_DRIVERS_LICENSE"
	EvidenceTypeAccountChangeInformation                        EvidenceType = "ACCOUNT_CHANGE_INFORMATION"
	EvidenceTypeDeliveryAddress                                 EvidenceType = "DELIVERY_ADDRESS"
	EvidenceTypeConfirmationOfResolution                        EvidenceType = "CONFIRMATION_OF_RESOLUTION"
	EvidenceTypeMerchantResponse                                EvidenceType = "MERCHANT_RESPONSE"
	EvidenceTypePermissionDescription                           EvidenceType = "PERMISSION_DESCRIPTION"
	EvidenceTypeStatusOfMerchandise                             EvidenceType = "STATUS_OF_MERCHANDISE"
	EvidenceTypeLostCardDetails                                 EvidenceType = "LOST_CARD_DETAILS"
	EvidenceTypeLastValidTransactionDetails                     EvidenceType = "LAST_VALID_TRANSACTION_DETAILS"
	EvidenceTypeAdditionalProofOfReturn                         EvidenceType = "ADDITIONAL_PROOF_OF_RETURN"
	EvidenceTypeOrderDetails                                    EvidenceType = "ORDER_DETAILS"
	EvidenceTypeListingUrl                                      EvidenceType = "LISTING_URL"
	EvidenceTypeShippingInsurance                               EvidenceType = "SHIPPING_INSURANCE"
	EvidenceTypeBuyerResponse                                   EvidenceType = "BUYER_RESPONSE"
	EvidenceTypePhotosOfShippedItem                             EvidenceType = "PHOTOS_OF_SHIPPED_ITEM"
	EvidenceTypeOther                                           EvidenceType = "OTHER"
	EvidenceTypeCancellationDetails                             EvidenceType = "CANCELLATION_DETAILS"
	EvidenceTypeMerchantContactDetails                          EvidenceType = "MERCHANT_CONTACT_DETAILS"
	EvidenceTypeItemDetails                                     EvidenceType = "ITEM_DETAILS"
	EvidenceTypeCorrectRecipientInformation                     EvidenceType = "CORRECT_RECIPIENT_INFORMATION"
	EvidenceTypeExplanationOfFundsNotDelivered                  EvidenceType = "EXPLANATION_OF_FUNDS_NOT_DELIVERED"
)
