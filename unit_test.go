package paypal

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type webprofileTestServer struct {
	t *testing.T
}

func TestNewClient(t *testing.T) {
	c, err := NewClient("", "", "")
	if err == nil {
		t.Errorf("Expected error for NewClient('','','')")
	}
	if c != nil {
		t.Errorf("Expected nil Client for NewClient('','',''), got %v", c)
	}

	c, err = NewClient("1", "2", "3")
	if err != nil {
		t.Errorf("Not expected error for NewClient(1, 2, 3), got %v", err)
	}
	if c == nil {
		t.Errorf("Expected non-nil Client for NewClient(1, 2, 3)")
	}
}

func TestTypeUserInfo(t *testing.T) {
	response := `{
    "user_id": "https://www.paypal.com/webapps/auth/server/64ghr894040044",
    "name": "Peter Pepper",
    "given_name": "Peter",
    "family_name": "Pepper",
    "email": "ppuser@example.com"
    }`

	u := &UserInfo{}
	err := json.Unmarshal([]byte(response), u)
	if err != nil {
		t.Errorf("UserInfo Unmarshal failed")
	}

	if u.ID != "https://www.paypal.com/webapps/auth/server/64ghr894040044" ||
		u.Name != "Peter Pepper" ||
		u.GivenName != "Peter" ||
		u.FamilyName != "Pepper" ||
		u.Email != "ppuser@example.com" {
		t.Errorf("UserInfo decoded result is incorrect, Given: %v", u)
	}
}

func TestTypeItem(t *testing.T) {
	response := `{
    "name":"Item",
    "quantity":"1"
}`

	i := &Item{}
	err := json.Unmarshal([]byte(response), i)
	if err != nil {
		t.Errorf("Item Unmarshal failed")
	}

	if i.Name != "Item" ||
		i.Quantity != "1" {
		t.Errorf("Item decoded result is incorrect, Given: %v", i)
	}
}

func TestTypeErrorResponseOne(t *testing.T) {
	response := `{
		"name":"USER_BUSINESS_ERROR",
		"message":"User business error.",
		"debug_id":"f05063556a338",
		"information_link":"https://developer.paypal.com/docs/api/payments.payouts-batch/#errors",
		"details":[
			{
				"field":"SENDER_BATCH_ID",
				"issue":"Batch with given sender_batch_id already exists",
				"link":[
					{
						"href":"https://api.sandbox.paypal.com/v1/payments/payouts/CR9VS2K4X4846",
						"rel":"self",
						"method":"GET"
					}
				]
			}
		]
	}`

	i := &ErrorResponse{}
	err := json.Unmarshal([]byte(response), i)
	if err != nil {
		t.Errorf("ErrorResponse Unmarshal failed")
	}

	if i.Name != "USER_BUSINESS_ERROR" ||
		i.Message != "User business error." ||
		i.DebugID != "f05063556a338" ||
		i.InformationLink != "https://developer.paypal.com/docs/api/payments.payouts-batch/#errors" ||
		len(i.Details) != 1 ||
		i.Details[0].Field != "SENDER_BATCH_ID" ||
		i.Details[0].Issue != "Batch with given sender_batch_id already exists" ||
		len(i.Details[0].Links) != 1 ||
		i.Details[0].Links[0].Href != "https://api.sandbox.paypal.com/v1/payments/payouts/CR9VS2K4X4846" {
		t.Errorf("ErrorResponse decoded result is incorrect, Given: %v", i)
	}
}

func TestTypeErrorResponseTwo(t *testing.T) {
	response := `{
		"name":"VALIDATION_ERROR",
		"message":"Invalid request - see details.",
		"debug_id":"662121ee369c0",
		"information_link":"https://developer.paypal.com/docs/api/payments.payouts-batch/#errors",
		"details":[
			{
				"field":"items[0].recipient_type",
				"issue":"Value is invalid (must be EMAILor PAYPAL_ID or PHONE)"
			}
		]
	}`

	i := &ErrorResponse{}
	err := json.Unmarshal([]byte(response), i)
	if err != nil {
		t.Errorf("ErrorResponse Unmarshal failed")
	}

	if i.Name != "VALIDATION_ERROR" ||
		i.Message != "Invalid request - see details." ||
		len(i.Details) != 1 ||
		i.Details[0].Field != "items[0].recipient_type" {
		t.Errorf("ErrorResponse decoded result is incorrect, Given: %v", i)
	}
}

func TestTypePayoutResponse(t *testing.T) {
	response := `{
		"batch_header":{
			"payout_batch_id":"G4E6WJE6Y4853",
			"batch_status":"SUCCESS",
			"time_created":"2017-11-01T23:08:25Z",
			"time_completed":"2017-11-01T23:08:46Z",
			"sender_batch_header":{
				"sender_batch_id":"2017110109",
				"email_subject":"Payment"
			},
			"amount":{
				"currency":"USD",
				"value":"6.37"
			},
			"fees":{
				"currency":"USD",
				"value":"0.25"
			}
		},
		"items":[
			{
				"payout_item_id":"9T35G83YA546X",
				"transaction_id":"4T328230B1D337285",
				"transaction_status":"UNCLAIMED",
				"payout_item_fee":{
					"currency":"USD",
					"value":"0.25"
				},
				"payout_batch_id":"G4E6WJE6Y4853",
				"payout_item":{
					"recipient_type":"EMAIL",
					"amount":{
						"currency":"USD",
						"value":"6.37"
					},
					"note":"Optional note",
					"receiver":"ppuser@example.com",
					"sender_item_id":"1"
				},
				"time_processed":"2017-11-01T23:08:43Z",
				"errors":{
					"name":"RECEIVER_UNREGISTERED",
					"message":"Receiver is unregistered",
					"information_link":"https://developer.paypal.com/docs/api/payments.payouts-batch/#errors",
					"details":[]
				},
				"links":[
					{
						"href":"https://api.sandbox.paypal.com/v1/payments/payouts-item/9T35G83YA546X",
						"rel":"item",
						"method":"GET",
						"encType":"application/json"
					}
				]
			}
		],
		"links":[
			{
				"href":"https://api.sandbox.paypal.com/v1/payments/payouts/G4E6WJE6Y4853?page_size=1000&page=1",
				"rel":"self",
				"method":"GET",
				"encType":"application/json"
			}
		]
	}`

	pr := &PayoutResponse{}
	err := json.Unmarshal([]byte(response), pr)
	if err != nil {
		t.Errorf("PayoutResponse Unmarshal failed")
	}

	if pr.BatchHeader.BatchStatus != "SUCCESS" ||
		pr.BatchHeader.PayoutBatchID != "G4E6WJE6Y4853" ||
		len(pr.Items) != 1 ||
		pr.Items[0].PayoutItemID != "9T35G83YA546X" ||
		pr.Items[0].TransactionID != "4T328230B1D337285" ||
		pr.Items[0].TransactionStatus != "UNCLAIMED" ||
		pr.Items[0].Error.Name != "RECEIVER_UNREGISTERED" {
		t.Errorf("PayoutResponse decoded result is incorrect, Given: %v", pr)
	}
}

func TestOrderUnmarshal(t *testing.T) {
	response := `{
		"id": "5O190127TN364715T",
		"status": "CREATED",
		"links": [
		  {
			"href": "https://api.paypal.com/v2/checkout/orders/5O190127TN364715T",
			"rel": "self",
			"method": "GET"
		  },
		  {
			"href": "https://api.sandbox.paypal.com/checkoutnow?token=5O190127TN364715T",
			"rel": "approve",
			"method": "GET"
		  },
		  {
			"href": "https://api.paypal.com/v2/checkout/orders/5O190127TN364715T/capture",
			"rel": "capture",
			"method": "POST"
		  }
		]
	}`

	order := &Order{}
	err := json.Unmarshal([]byte(response), order)
	if err != nil {
		t.Errorf("Order Unmarshal failed")
	}

	if order.ID != "5O190127TN364715T" ||
		order.Status != "CREATED" ||
		order.Links[0].Href != "https://api.paypal.com/v2/checkout/orders/5O190127TN364715T" {
		t.Errorf("Order decoded result is incorrect, Given: %+v", order)
	}
}

func TestTypePayoutItemResponse(t *testing.T) {
	response := `{
		"payout_item_id":"9T35G83YA546X",
		"transaction_id":"4T328230B1D337285",
		"transaction_status":"UNCLAIMED",
		"payout_item_fee":{
			"currency":"USD",
			"value":"0.25"
		},
		"payout_batch_id":"G4E6WJE6Y4853",
		"payout_item":{
			"recipient_type":"EMAIL",
			"amount":{
				"currency":"USD",
				"value":"6.37"
			},
			"note":"Optional note",
			"receiver":"ppuser@example.com",
			"sender_item_id":"1"
		},
		"time_processed":"2017-11-01T23:08:43Z",
		"errors":{
			"name":"RECEIVER_UNREGISTERED",
			"message":"Receiver is unregistered",
			"information_link":"https://developer.paypal.com/docs/api/payments.payouts-batch/#errors",
			"details":[]
		},
		"links":[
			{
				"href":"https://api.sandbox.paypal.com/v1/payments/payouts-item/3YA546X9T35G8",
				"rel":"self",
				"method":"GET",
				"encType":"application/json"
			},
			{
				"href":"https://api.sandbox.paypal.com/v1/payments/payouts/6Y4853G4E6WJE",
				"rel":"batch",
				"method":"GET",
				"encType":"application/json"
			}
		]
	}`

	pir := &PayoutItemResponse{}
	err := json.Unmarshal([]byte(response), pir)
	if err != nil {
		t.Errorf("PayoutItemResponse Unmarshal failed")
	}

	if pir.PayoutItemID != "9T35G83YA546X" ||
		pir.PayoutBatchID != "G4E6WJE6Y4853" ||
		pir.TransactionID != "4T328230B1D337285" ||
		pir.TransactionStatus != "UNCLAIMED" ||
		pir.Error.Name != "RECEIVER_UNREGISTERED" {
		t.Errorf("PayoutItemResponse decoded result is incorrect, Given: %+v", pir)
	}
}

func TestTypePaymentPatch(t *testing.T) {
	// test unmarshaling
	response := `{
		"op": "replace",
		"path": "/transactions/0/amount",
		"value": "5"
	}`
	pp := &PaymentPatch{}
	err := json.Unmarshal([]byte(response), pp)
	if err != nil {
		t.Errorf("TestTypePaymentPatch Unmarshal failed")
	}
	if pp.Operation != "replace" ||
		pp.Path != "/transactions/0/amount" ||
		pp.Value != "5" {
		t.Errorf("PaymentPatch decoded result is incorrect, Given: %+v", pp)
	}
}

func TestTypePaymentPatchMarshal(t *testing.T) {
	// test marshalling
	p2 := &PaymentPatch{
		Operation: "add",
		Path:      "/transactions/0/amount",
		Value: map[string]interface{}{
			"total":    "18.37",
			"currency": "EUR",
			"details": map[string]interface{}{
				"subtotal": "13.37",
				"shipping": "5.00",
			},
		},
	}
	p2expectedresponse := `{"op":"add","path":"/transactions/0/amount","value":{"currency":"EUR","details":{"shipping":"5.00","subtotal":"13.37"},"total":"18.37"}}`
	response2, _ := json.Marshal(p2)
	if string(response2) != string(p2expectedresponse) {
		t.Errorf("PaymentPatch response2 is incorrect,\n Given:    %+v\n Expected: %+v", string(response2), string(p2expectedresponse))
	}
}

// ServeHTTP implements http.Handler
func (ts *webprofileTestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ts.t.Log(r.RequestURI)
	if r.RequestURI == "/v1/payment-experience/web-profiles" {
		if r.Method == "POST" {
			ts.create(w, r)
		}
		if r.Method == "GET" {
			ts.list(w, r)
		}
	}
	if r.RequestURI == "/v1/payment-experience/web-profiles/XP-CP6S-W9DY-96H8-MVN2" {
		if r.Method == "GET" {
			ts.getvalid(w, r)
		}
		if r.Method == "PUT" {
			ts.updatevalid(w, r)
		}
		if r.Method == "DELETE" {
			ts.deletevalid(w, r)
		}
	}
	if r.RequestURI == "/v1/payment-experience/web-profiles/foobar" {
		if r.Method == "GET" {
			ts.getinvalid(w, r)
		}
		if r.Method == "PUT" {
			ts.updateinvalid(w, r)
		}
		if r.Method == "DELETE" {
			ts.deleteinvalid(w, r)
		}
	}
}

func (ts *webprofileTestServer) create(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.Unmarshal(body, &data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var raw map[string]string

	w.Header().Set("Content-Type", "application/json")

	if name, ok := data["name"]; !ok || name == "" {
		raw = map[string]string{
			"name":    "VALIDATION_ERROR",
			"message": "should have name",
		}
		w.WriteHeader(http.StatusBadRequest)
	} else {
		raw = map[string]string{
			"id": "XP-CP6S-W9DY-96H8-MVN2",
		}
		w.WriteHeader(http.StatusCreated)
	}

	res, _ := json.Marshal(raw)
	w.Write(res)
}

func (ts *webprofileTestServer) createWithoutName(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.Unmarshal(body, &data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var raw map[string]string

	w.Header().Set("Content-Type", "application/json")

	raw = map[string]string{
		"id": "B-12345678901234567",
	}
	w.WriteHeader(http.StatusCreated)

	res, _ := json.Marshal(raw)
	w.Write(res)
}

func (ts *webprofileTestServer) updatevalid(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.Unmarshal(body, &data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if ID, ok := data["id"]; !ok || ID != "XP-CP6S-W9DY-96H8-MVN2" {
		raw := map[string]string{
			"name":    "INVALID_RESOURCE_ID",
			"message": "id invalid",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(raw)
		w.Write(res)
		return
	}

	if name, ok := data["name"]; !ok || name == "" {
		raw := map[string]string{
			"name":    "VALIDATION_ERROR",
			"message": "should have name",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(raw)
		w.Write(res)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (ts *webprofileTestServer) updateinvalid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	raw := map[string]interface{}{
		"name":    "INVALID_RESOURCE_ID",
		"message": "foobar not found",
	}

	res, _ := json.Marshal(raw)
	w.Write(res)
}

func (ts *webprofileTestServer) getvalid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	raw := map[string]interface{}{
		"id":   "XP-CP6S-W9DY-96H8-MVN2",
		"name": "YeowZa! T-Shirt Shop",
		"presentation": map[string]interface{}{
			"brand_name":  "YeowZa! Paypal",
			"logo_image":  "http://www.yeowza.com",
			"locale_code": "US",
		},

		"input_fields": map[string]interface{}{
			"allow_note":       true,
			"no_shipping":      0,
			"address_override": 1,
		},

		"flow_config": map[string]interface{}{
			"landing_page_type":    "Billing",
			"bank_txn_pending_url": "http://www.yeowza.com",
		},
	}

	res, _ := json.Marshal(raw)
	w.Write(res)
}

func (ts *webprofileTestServer) getinvalid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	raw := map[string]interface{}{
		"name":    "INVALID_RESOURCE_ID",
		"message": "foobar not found",
	}

	res, _ := json.Marshal(raw)
	w.Write(res)
}

func (ts *webprofileTestServer) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	raw := []interface{}{
		map[string]interface{}{
			"id":   "XP-CP6S-W9DY-96H8-MVN2",
			"name": "YeowZa! T-Shirt Shop",
		},
		map[string]interface{}{
			"id":   "XP-96H8-MVN2-CP6S-W9DY",
			"name": "Shop T-Shirt YeowZa! ",
		},
	}

	res, _ := json.Marshal(raw)
	w.Write(res)
}

func (ts *webprofileTestServer) deleteinvalid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	raw := map[string]interface{}{
		"name":    "INVALID_RESOURCE_ID",
		"message": "foobar not found",
	}

	res, _ := json.Marshal(raw)
	w.Write(res)
}

func (ts *webprofileTestServer) deletevalid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func TestCreateWebProfile_valid(t *testing.T) {
	ts := httptest.NewServer(&webprofileTestServer{t: t})
	defer ts.Close()

	c, _ := NewClient("foo", "bar", ts.URL)

	wp := WebProfile{
		Name: "YeowZa! T-Shirt Shop",
		Presentation: Presentation{
			BrandName:  "YeowZa! Paypal",
			LogoImage:  "http://www.yeowza.com",
			LocaleCode: "US",
		},

		InputFields: InputFields{
			AllowNote:       true,
			NoShipping:      NoShippingDisplay,
			AddressOverride: AddrOverrideFromCall,
		},

		FlowConfig: FlowConfig{
			LandingPageType:   LandingPageTypeBilling,
			BankTXNPendingURL: "http://www.yeowza.com",
		},
	}

	res, err := c.CreateWebProfile(context.Background(), wp)

	if err != nil {
		t.Fatal(err)
	}

	if res.ID != "XP-CP6S-W9DY-96H8-MVN2" {
		t.Fatalf("expecting response to have ID = `XP-CP6S-W9DY-96H8-MVN2` got `%s`", res.ID)
	}
}

func TestCreateWebProfile_invalid(t *testing.T) {
	ts := httptest.NewServer(&webprofileTestServer{t: t})
	defer ts.Close()

	c, _ := NewClient("foo", "bar", ts.URL)

	wp := WebProfile{}

	_, err := c.CreateWebProfile(context.Background(), wp)

	if err == nil {
		t.Fatalf("expecting an error got nil")
	}
}

func TestGetWebProfile_valid(t *testing.T) {
	ts := httptest.NewServer(&webprofileTestServer{t: t})
	defer ts.Close()

	c, _ := NewClient("foo", "bar", ts.URL)

	res, err := c.GetWebProfile(context.Background(), "XP-CP6S-W9DY-96H8-MVN2")

	if err != nil {
		t.Fatal(err)
	}

	if res.ID != "XP-CP6S-W9DY-96H8-MVN2" {
		t.Fatalf("expecting res.ID to have value = `XP-CP6S-W9DY-96H8-MVN2` but got `%s`", res.ID)
	}

	if res.Name != "YeowZa! T-Shirt Shop" {
		t.Fatalf("expecting res.Name to have value = `YeowZa! T-Shirt Shop` but got `%s`", res.Name)
	}
}

func TestGetWebProfile_invalid(t *testing.T) {
	ts := httptest.NewServer(&webprofileTestServer{t: t})
	defer ts.Close()

	c, _ := NewClient("foo", "bar", ts.URL)

	_, err := c.GetWebProfile(context.Background(), "foobar")

	if err == nil {
		t.Fatalf("expecting an error got nil")
	}
}

func TestGetWebProfiles(t *testing.T) {
	ts := httptest.NewServer(&webprofileTestServer{t: t})
	defer ts.Close()

	c, _ := NewClient("foo", "bar", ts.URL)

	res, err := c.GetWebProfiles(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 2 {
		t.Fatalf("expecting two results got %d", len(res))
	}
}

func TestSetWebProfile_valid(t *testing.T) {
	ts := httptest.NewServer(&webprofileTestServer{t: t})
	defer ts.Close()

	c, _ := NewClient("foo", "bar", ts.URL)

	wp := WebProfile{
		ID:   "XP-CP6S-W9DY-96H8-MVN2",
		Name: "Shop T-Shirt YeowZa!",
	}

	err := c.SetWebProfile(context.Background(), wp)

	if err != nil {
		t.Fatal(err)
	}

}

func TestSetWebProfile_invalid(t *testing.T) {
	ts := httptest.NewServer(&webprofileTestServer{t: t})
	defer ts.Close()

	c, _ := NewClient("foo", "bar", ts.URL)

	wp := WebProfile{
		ID: "foobar",
	}

	err := c.SetWebProfile(context.Background(), wp)

	if err == nil {
		t.Fatal(err)
	}

	wp = WebProfile{}

	err = c.SetWebProfile(context.Background(), wp)

	if err == nil {
		t.Fatal(err)
	}
}

func TestDeleteWebProfile_valid(t *testing.T) {
	ts := httptest.NewServer(&webprofileTestServer{t: t})
	defer ts.Close()

	c, _ := NewClient("foo", "bar", ts.URL)

	wp := WebProfile{
		ID:   "XP-CP6S-W9DY-96H8-MVN2",
		Name: "Shop T-Shirt YeowZa!",
	}

	err := c.SetWebProfile(context.Background(), wp)

	if err != nil {
		t.Fatal(err)
	}

}

func TestDeleteWebProfile_invalid(t *testing.T) {
	ts := httptest.NewServer(&webprofileTestServer{t: t})
	defer ts.Close()

	c, _ := NewClient("foo", "bar", ts.URL)

	err := c.DeleteWebProfile(context.Background(), "foobar")

	if err == nil {
		t.Fatal(err)
	}

}