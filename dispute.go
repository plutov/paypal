package paypal

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// ListDisputes - Lists disputes with a summary set of details.
// Endpoint: GET /v1/customer/disputes
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_list
func (c *Client) ListDisputes(ctx context.Context, req *ListDisputesRequest) (*ListDisputesResponse, error) {
	response := &ListDisputesResponse{}

	r, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s/v1/customer/disputes", c.APIBase), nil)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()

	if req.StartTime != nil {
		q.Add("start_time", req.StartTime.Format(time.RFC3339))
	}
	if req.DisputedTransactionId != nil {
		q.Add("disputed_transaction_id", *req.DisputedTransactionId)
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
	if err = c.SendWithAuth(req, response); err != nil {
		return nil, err
	}
	return response, err
}

// Partially updates a dispute, by ID
// Endpoint: PATCH /v1/customer/disputes/{id}
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_patch
func (c *Client) UpdateDispute(ctx context.Context, disputeId string, params *UpdateDisputeParams) error {
	req, err := c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("%s/v1/customer/disputes/%s", c.APIBase, disputeId), params)
	if err != nil {
		return err
	}
	return c.SendWithAuth(req, nil)
}

// Provides evidence for a dispute, by ID
// Endpoint: POST /v1/customer/disputes/{id}/provide-evidence
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_provide-evidence
func (c *Client) DisputeProvideEvidence(ctx context.Context, disputeId string, params *DisputeProvideEvidenceParams) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/customer/disputes/%s/provide-evidence", c.APIBase, disputeId), params)
	if err != nil {
		return err
	}
	return c.SendWithAuth(req, nil)
}

// Appeals a dispute, by ID
// Endpoint: POST /v1/customer/disputes/{id}/appeal
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_appeal
func (c *Client) DisputeAppeal(ctx context.Context, disputeId string, params *DisputeAppealParams) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/customer/disputes/%s/appeal", c.APIBase, disputeId), params)
	if err != nil {
		return err
	}
	return c.SendWithAuth(req, nil)
}

// Accepts liability for a claim, by ID
// Endpoint: POST /v1/customer/disputes/{id}/accept-claim
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_accept-claim
func (c *Client) DisputeAcceptClaim(ctx context.Context, disputeId string, params *DisputeAcceptClaimParams) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/customer/disputes/%s/accept-claim", c.APIBase, disputeId), params)
	if err != nil {
		return err
	}
	return c.SendWithAuth(req, nil)
}

// Settles a dispute in either the customer's or merchant's favor.
// Endpoint: POST /v1/customer/disputes/{id}/adjudicate
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_adjudicate
func (c *Client) SettleDispute(ctx context.Context, disputeId string, adjudicateOutcome AdjudicationOutcome) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/customer/disputes/%s/adjudicate", c.APIBase, disputeId), map[string]string{"adjudication_outcome": string(adjudicateOutcome)})
	if err != nil {
		return err
	}
	return c.SendWithAuth(req, nil)
}

// Updates the status of a dispute, by ID
// Endpoint: POST /v1/customer/disputes/{id}/require-evidence
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_require-evidence
// Important: This method is for sandbox use only.
func (c *Client) DisputeUpdateStatus(ctx context.Context, disputeId string, adjudicateOutcome AdjudicationOutcome) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/customer/disputes/%s/require-evidence", c.APIBase, disputeId), map[string]string{"action": string(adjudicateOutcome)})
	if err != nil {
		return err
	}
	return c.SendWithAuth(req, nil)
}

// Escalates the dispute, by ID
// Endpoint: POST /v1/customer/disputes/{id}/escalate
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_escalate
func (c *Client) DisputeEscalateToClaim(ctx context.Context, disputeId string, note string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/customer/disputes/%s/escalate", c.APIBase, disputeId), map[string]string{"note": note})
	if err != nil {
		return err
	}
	return c.SendWithAuth(req, nil)
}

// Sends a message about a dispute, by ID, to the other party in the dispute.
// Endpoint: POST /v1/customer/disputes/{id}/send-message
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_send-message
func (c *Client) DisputeSendMessageToOtherParty(ctx context.Context, disputeId string, message string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/customer/disputes/%s/send-message", c.APIBase, disputeId), map[string]string{"message": message})
	if err != nil {
		return err
	}
	return c.SendWithAuth(req, nil)
}

// Makes an offer to the other party to resolve a dispute, by ID
// Endpoint: POST /v1/customer/disputes/{id}/make-offer
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_make-offer
func (c *Client) DisputeMakeOffer(ctx context.Context, disputeId string, params *DisputeMakeOfferParams) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/customer/disputes/%s/make-offer", c.APIBase, disputeId), params)
	if err != nil {
		return err
	}
	return c.SendWithAuth(req, nil)
}

// The customer accepts the offer from merchant to resolve a dispute, by ID
// Endpoint: POST /v1/customer/disputes/{id}/accept-offer
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_accept-offer
func (c *Client) DisputeAcceptOffer(ctx context.Context, disputeId string, note string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/customer/disputes/%s/accept-offer", c.APIBase, disputeId), map[string]string{"note": note})
	if err != nil {
		return err
	}
	return c.SendWithAuth(req, nil)
}

// Denies an offer that the merchant proposes for a dispute, by ID.
// Endpoint: POST /v1/customer/disputes/{id}/deny-offer
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_deny-offer
func (c *Client) DisputeDenyOffer(ctx context.Context, disputeId string, note string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/customer/disputes/%s/deny-offer", c.APIBase, disputeId), map[string]string{"note": note})
	if err != nil {
		return err
	}
	return c.SendWithAuth(req, nil)
}

// Acknowledges that the customer returned an item for a dispute, by ID.
// Endpoint: POST /v1/customer/disputes/{id}/acknowledge-return-item
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_acknowledge-return-item
func (c *Client) DisputeAcknowledgeReturnItem(ctx context.Context, disputeId string, params *DisputeAcknowledgeReturnItemParams) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/customer/disputes/%s/acknowledge-return-item", c.APIBase, disputeId), params)
	if err != nil {
		return err
	}
	return c.SendWithAuth(req, nil)
}

// Provides supporting information for a dispute, by ID.
// Endpoint: POST /v1/customer/disputes/{id}/provide-supporting-info
// https://developer.paypal.com/docs/api/customer-disputes/v1/#disputes_provide-supporting-info
func (c *Client) DisputeProvideSupportingInfo(ctx context.Context, disputeId string, notes string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/customer/disputes/%s/provide-supporting-info", c.APIBase, disputeId), map[string]string{"notes": notes})
	if err != nil {
		return err
	}
	return c.SendWithAuth(req, nil)
}
