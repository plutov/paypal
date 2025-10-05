package mockserver

import (
	"context"
	"net/http"
	"time"

	srv "github.com/plutov/paypal/v4/mockserver/payments_payouts_batch_v1"
	srvShipping "github.com/plutov/paypal/v4/mockserver/shipping_shipment_tracking_v1"
)

type MockServer struct{}

var (
	dummyTime              = time.Now()
	dummyTransactionStatus = srv.TransactionEnum("SUCCESS")
	dummyTracker           = srvShipping.TrackersGet200JSONResponse{
		TransactionId: "TXN123",
		Status:        "SHIPPED",
	}
)

type dummyPutResponse struct{}

func (d dummyPutResponse) VisitTrackersPutResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type dummyBatchPostResponse struct{}

func (d dummyBatchPostResponse) VisitTrackersBatchPostResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

var dummyPayout = srv.Payout{
	BatchHeader: &srv.PayoutHeader{
		BatchStatus:   "SUCCESS",
		PayoutBatchId: "BATCH123",
		SenderBatchHeader: srv.PayoutSenderBatchHeader{
			SenderBatchId: stringPtr("SENDER123"),
		},
		TimeCreated: &dummyTime,
	},
}

var dummyPayoutItem = srv.PayoutItem{
	PayoutBatchId: "BATCH123",
	PayoutItem: srv.PayoutItemDetail{
		Amount: srv.Currency{
			Currency: "USD",
			Value:    "10.00",
		},
		Receiver: "receiver@example.com",
	},
	PayoutItemId:      "ITEM123",
	TimeProcessed:     &dummyTime,
	TransactionStatus: &dummyTransactionStatus,
}

var dummyPayoutBatch = srv.PayoutBatch{
	BatchHeader: &srv.PayoutBatchHeader{
		BatchStatus:   "SUCCESS",
		PayoutBatchId: "BATCH123",
		SenderBatchHeader: srv.PayoutSenderBatchHeader{
			SenderBatchId: stringPtr("SENDER123"),
		},
		TimeCreated: &dummyTime,
	},
	TotalItems: intPtr(1),
	TotalPages: intPtr(1),
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func (m *MockServer) PayoutsPost(ctx context.Context, request srv.PayoutsPostRequestObject) (srv.PayoutsPostResponseObject, error) {
	return srv.PayoutsPost201JSONResponse(dummyPayout), nil
}

func (m *MockServer) PayoutsItemGet(ctx context.Context, request srv.PayoutsItemGetRequestObject) (srv.PayoutsItemGetResponseObject, error) {
	return srv.PayoutsItemGet200JSONResponse(dummyPayoutItem), nil
}

func (m *MockServer) PayoutsItemCancel(ctx context.Context, request srv.PayoutsItemCancelRequestObject) (srv.PayoutsItemCancelResponseObject, error) {
	return srv.PayoutsItemCancel200JSONResponse(dummyPayoutItem), nil
}

func (m *MockServer) PayoutsGet(ctx context.Context, request srv.PayoutsGetRequestObject) (srv.PayoutsGetResponseObject, error) {
	return srv.PayoutsGet200JSONResponse(dummyPayoutBatch), nil
}

func (m *MockServer) TrackersBatchGet(ctx context.Context, request srvShipping.TrackersBatchGetRequestObject) (srvShipping.TrackersBatchGetResponseObject, error) {
	return srvShipping.TrackersBatchGet200JSONResponse(srvShipping.Tracker{
		TransactionId: "TXN123",
		Status:        "SHIPPED",
	}), nil
}

func (m *MockServer) TrackersPost(ctx context.Context, request srvShipping.TrackersPostRequestObject) (srvShipping.TrackersPostResponseObject, error) {
	return srvShipping.TrackersPost200JSONResponse(srvShipping.TrackerIdentifierCollection{
		TrackerIdentifiers: &[]srvShipping.TrackerIdentifier{{TransactionId: "TXN123"}},
	}), nil
}

func (m *MockServer) TrackersBatchPost(ctx context.Context, request srvShipping.TrackersBatchPostRequestObject) (srvShipping.TrackersBatchPostResponseObject, error) {
	return dummyBatchPostResponse{}, nil
}

func (m *MockServer) TrackersGet(ctx context.Context, request srvShipping.TrackersGetRequestObject) (srvShipping.TrackersGetResponseObject, error) {
	return dummyTracker, nil
}

func (m *MockServer) TrackersPut(ctx context.Context, request srvShipping.TrackersPutRequestObject) (srvShipping.TrackersPutResponseObject, error) {
	return dummyPutResponse{}, nil
}
