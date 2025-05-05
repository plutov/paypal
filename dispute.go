package paypal

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

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
	CreateTime            time.Time                  `json:"create_time,omitempty"`
	UpdateTime            time.Time                  `json:"update_time,omitempty"`
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
	Buyer              Buyer  `json:"buyer,omitempty"`
	Seller             Seller `json:"seller,omitempty"`
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
	CreateTime            time.Time                    `json:"create_time,omitempty"`
	UpdateTime            time.Time                    `json:"update_time,omitempty"`
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
	CreateTime          time.Time `json:"create_time,omitempty"`
	TransactionStatus   string    `json:"transaction_status,omitempty"`
	GrossAmount         *Money    `json:"gross_amount,omitempty"`
	Buyer               Buyer     `json:"buyer,omitempty"`
	Seller              Seller    `json:"seller,omitempty"`
}

type Extensions struct {
	MerchandizeDisputeProperties MerchandizeDisputeProperties `json:"merchandize_dispute_properties,omitempty"`
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
	TimePosted time.Time   `json:"time_posted,omitempty"`
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
	OfferTime             time.Time `json:"offer_time,omitempty"`
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
	CreateTime time.Time `json:"create_time,omitempty"`
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

// ListDisputes - Lists disputes with a summary set of details.
// Endpoint: GET /v1/customer/disputes
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_list
func (c *Client) ListDisputes(ctx context.Context, req *ListDisputesRequest) (*ListDisputesResponse, error) {
	response := &ListDisputesResponse{}

	r, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/customer/disputes"), nil)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()

	if req.StartTime != nil {
		q.Add("start_time", req.StartTime.Format(time.RFC3339))
	}
	if req.DisputedTransactionId != nil {
		q.Add("disputed_transactioniId", *req.DisputedTransactionId)
	}
	if req.PageSize != nil {
		q.Add("page_size", strconv.Itoa(*req.PageSize))
	}
	if req.NextPageToken != nil {
		q.Add("next_page_token", *req.NextPageToken)
	}
	if req.DisputeState != nil {
		q.Add("dispute_state", *req.DisputeState)
	}
	if req.UpdateTimeBefore != nil {
		q.Add("update_time_before", req.UpdateTimeBefore.Format(time.RFC3339))
	}
	if req.UpdateTimeAfter != nil {
		q.Add("update_time_after", req.UpdateTimeAfter.Format(time.RFC3339))
	}

	r.URL.RawQuery = q.Encode()

	if err = c.SendWithAuth(r, response); err != nil {
		return nil, err
	}

	return response, nil
}

// Shows details for a dispute, by ID
// Endpoint: GET /v1/customer/disputes/{dispute_id}
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_get
func (c *Client) GetDisputeDetail(ctx context.Context, DisputeID string) (*GetDisputeDetailResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/v1/customer/disputes/%s", c.APIBase, DisputeID), nil)
	response := &GetDisputeDetailResponse{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// Partially updates a dispute, by ID
// Endpoint: PATCH /v1/customer/disputes/{id}
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_patch
func (c *Client) UpdateDispute(ctx context.Context, disputeId string, params *UpdateDisputeParams) error {
	req, err := c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/customer/disputes/", disputeId), params)
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Provides evidence for a dispute, by ID
// Endpoint: POST /v1/customer/disputes/{id}/provide-evidence
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_provide-evidence
func (c *Client) DisputeProvideEvidence(ctx context.Context, disputeId string, params *DisputeProvideEvidenceParams) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s%s/provide-evidence", c.APIBase, "/v1/customer/disputes/", disputeId), params)
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Appeals a dispute, by ID
// Endpoint: POST /v1/customer/disputes/{id}/appeal
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_appeal
func (c *Client) DisputeAppeal(ctx context.Context, disputeId string, params *DisputeAppealParams) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s%s/appeal", c.APIBase, "/v1/customer/disputes/", disputeId), params)
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Accepts liability for a claim, by ID
// Endpoint: POST /v1/customer/disputes/{id}/accept-claim
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_accept-claim
func (c *Client) DisputeAcceptClaim(ctx context.Context, disputeId string, params *DisputeAcceptClaimParams) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s%s/accept-claim", c.APIBase, "/v1/customer/disputes/", disputeId), params)
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Settles a dispute in either the customer's or merchant's favor.
// Endpoint: POST /v1/customer/disputes/{id}/adjudicate
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_adjudicate
func (c *Client) SettleDispute(ctx context.Context, disputeId string, adjudicateOutcome AdjudicationOutcome) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s%s/adjudicate", c.APIBase, "/v1/customer/disputes/", disputeId), map[string]string{"adjudication_outcome": string(adjudicateOutcome)})
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Updates the status of a dispute, by ID
// Endpoint: POST /v1/customer/disputes/{id}/require-evidence
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_require-evidence
// Important: This method is for sandbox use only.
func (c *Client) DisputeUpdateStatus(ctx context.Context, disputeId string, adjudicateOutcome AdjudicationOutcome) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s%s/require-evidence", c.APIBase, "/v1/customer/disputes/", disputeId), map[string]string{"action": string(adjudicateOutcome)})
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Escalates the dispute, by ID
// Endpoint: POST /v1/customer/disputes/{id}/escalate
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_escalate
func (c *Client) DisputeEscalateToClaim(ctx context.Context, disputeId string, note string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s%s/escalate", c.APIBase, "/v1/customer/disputes/", disputeId), map[string]string{"note": note})
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Sends a message about a dispute, by ID, to the other party in the dispute.
// Endpoint: POST /v1/customer/disputes/{id}/send-message
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_send-message
func (c *Client) DisputeSendMessageToOtherParty(ctx context.Context, disputeId string, message string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s%s/send-message", c.APIBase, "/v1/customer/disputes/", disputeId), map[string]string{"message": message})
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Makes an offer to the other party to resolve a dispute, by ID
// Endpoint: POST /v1/customer/disputes/{id}/make-offer
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_make-offer
func (c *Client) DisputeMakeOffer(ctx context.Context, disputeId string, params *DisputeMakeOfferParams) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s%s/make-offer", c.APIBase, "/v1/customer/disputes/", disputeId), params)
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// The customer accepts the offer from merchant to resolve a dispute, by ID
// Endpoint: POST /v1/customer/disputes/{id}/accept-offer
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_accept-offer
func (c *Client) DisputeAcceptOffer(ctx context.Context, disputeId string, note string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s%s/accept-offer", c.APIBase, "/v1/customer/disputes/", disputeId), map[string]string{"note": note})
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Denies an offer that the merchant proposes for a dispute, by ID.
// Endpoint: POST /v1/customer/disputes/{id}/deny-offer
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_deny-offer
func (c *Client) DisputeDenyOffer(ctx context.Context, disputeId string, note string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s%s/deny-offer", c.APIBase, "/v1/customer/disputes/", disputeId), map[string]string{"note": note})
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Acknowledges that the customer returned an item for a dispute, by ID.
// Endpoint: POST /v1/customer/disputes/{id}/acknowledge-return-item
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_acknowledge-return-item
func (c *Client) DisputeAcknowledgeReturnItem(ctx context.Context, disputeId string, params *DisputeAcknowledgeReturnItemParams) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s%s/acknowledge-return-item", c.APIBase, "/v1/customer/disputes/", disputeId), params)
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Provides supporting information for a dispute, by ID.
// Endpoint: POST /v1/customer/disputes/{id}/provide-supporting-info
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_provide-supporting-info
func (c *Client) DisputeProvideSupportingInfo(ctx context.Context, disputeId string, notes string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s%s/provide-supporting-info", c.APIBase, "/v1/customer/disputes/", disputeId), map[string]string{"notes": notes})
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}
