package paypal

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type TransactionSearchRequest struct {
	TransactionID               *string
	TransactionType             *string
	TransactionStatus           *string
	TransactionAmount           *string
	TransactionCurrency         *string
	StartDate                   time.Time
	EndDate                     time.Time
	PaymentInstrumentType       *string
	StoreID                     *string
	TerminalID                  *string
	Fields                      *string
	BalanceAffectingRecordsOnly *string
	PageSize                    *int
	Page                        *int
}

type TransactionSearchResponse struct {
	TransactionDetails  []SearchTransactionDetails `json:"transaction_details"`
	AccountNumber       string                     `json:"account_number"`
	StartDate           JSONTime                   `json:"start_date"`
	EndDate             JSONTime                   `json:"end_date"`
	LastRefreshDatetime JSONTime                   `json:"last_refreshed_datetime"`
	Page                int                        `json:"page"`
	SharedListResponse
}

// ListTransactions - Use this to search PayPal transactions from the last 31 days.
// Endpoint: GET /v1/reporting/transactions
func (c *Client) ListTransactions(ctx context.Context, req *TransactionSearchRequest) (*TransactionSearchResponse, error) {
	response := &TransactionSearchResponse{}

	r, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/reporting/transactions"), nil)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()

	q.Add("start_date", req.StartDate.Format(time.RFC3339))
	q.Add("end_date", req.EndDate.Format(time.RFC3339))

	if req.TransactionID != nil {
		q.Add("transaction_id", *req.TransactionID)
	}
	if req.TransactionType != nil {
		q.Add("transaction_type", *req.TransactionType)
	}
	if req.TransactionStatus != nil {
		q.Add("transaction_status", *req.TransactionStatus)
	}
	if req.TransactionAmount != nil {
		q.Add("transaction_amount", *req.TransactionAmount)
	}
	if req.TransactionCurrency != nil {
		q.Add("transaction_currency", *req.TransactionCurrency)
	}
	if req.PaymentInstrumentType != nil {
		q.Add("payment_instrument_type", *req.PaymentInstrumentType)
	}
	if req.StoreID != nil {
		q.Add("store_id", *req.StoreID)
	}
	if req.TerminalID != nil {
		q.Add("terminal_id", *req.TerminalID)
	}
	if req.Fields != nil {
		q.Add("fields", *req.Fields)
	}
	if req.BalanceAffectingRecordsOnly != nil {
		q.Add("balance_affecting_records_only", *req.BalanceAffectingRecordsOnly)
	}
	if req.PageSize != nil {
		q.Add("page_size", strconv.Itoa(*req.PageSize))
	}
	if req.Page != nil {
		q.Add("page", strconv.Itoa(*req.Page))
	}

	r.URL.RawQuery = q.Encode()

	if err = c.SendWithAuth(r, response); err != nil {
		return nil, err
	}

	return response, nil
}
