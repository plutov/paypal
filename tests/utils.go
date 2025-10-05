package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mockserver "github.com/plutov/paypal/v4/mockserver"
	srvPayouts "github.com/plutov/paypal/v4/mockserver/payments_payouts_batch_v1"
	srvShipping "github.com/plutov/paypal/v4/mockserver/shipping_shipment_tracking_v1"
)

const (
	clientId     = "test-client-id"
	clientSecret = "test-client-secret"
)

func createTestServer() *httptest.Server {
	mock := &mockserver.MockServer{}
	mux := http.NewServeMux()

	payoutsSI := srvPayouts.NewStrictHandler(mock, nil)
	srvPayouts.HandlerWithOptions(payoutsSI, srvPayouts.StdHTTPServerOptions{
		BaseRouter: mux,
	})

	shippingSI := srvShipping.NewStrictHandler(mock, nil)
	srvShipping.HandlerWithOptions(shippingSI, srvShipping.StdHTTPServerOptions{
		BaseRouter: mux,
	})

	return httptest.NewServer(mux)
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func assertEqual(t *testing.T, expected, actual any) {
	t.Helper()
	if expected != actual {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}
