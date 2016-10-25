package paypalsdk

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type webprofileTestServer struct {
	t *testing.T
}

// ServeHTTP implements http.Handler
func (ts *webprofileTestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ts.t.Log(r.RequestURI)
	if r.RequestURI == "/v1/payment-experience/web-profiles" {
		if r.Method == http.MethodPost {
			ts.create(w, r)
		}
		if r.Method == http.MethodGet {
			ts.list(w, r)
		}
	}
	if r.RequestURI == "/v1/payment-experience/web-profiles/XP-CP6S-W9DY-96H8-MVN2" {
		if r.Method == http.MethodGet {
			ts.getvalid(w, r)
		}
		if r.Method == http.MethodPut {
			ts.updatevalid(w, r)
		}
		if r.Method == http.MethodDelete {
			ts.deletevalid(w, r)
		}
	}
	if r.RequestURI == "/v1/payment-experience/web-profiles/foobar" {
		if r.Method == http.MethodGet {
			ts.getinvalid(w, r)
		}
		if r.Method == http.MethodPut {
			ts.updateinvalid(w, r)
		}
		if r.Method == http.MethodDelete {
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

	res, err := c.CreateWebProfile(wp)

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

	_, err := c.CreateWebProfile(wp)

	if err == nil {
		t.Fatalf("expecting an error got nil")
	}
}

func TestGetWebProfile_valid(t *testing.T) {
	ts := httptest.NewServer(&webprofileTestServer{t: t})
	defer ts.Close()

	c, _ := NewClient("foo", "bar", ts.URL)

	res, err := c.GetWebProfile("XP-CP6S-W9DY-96H8-MVN2")

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

	_, err := c.GetWebProfile("foobar")

	if err == nil {
		t.Fatalf("expecting an error got nil")
	}
}

func TestGetWebProfiles(t *testing.T) {
	ts := httptest.NewServer(&webprofileTestServer{t: t})
	defer ts.Close()

	c, _ := NewClient("foo", "bar", ts.URL)

	res, err := c.GetWebProfiles()

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

	err := c.SetWebProfile(wp)

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

	err := c.SetWebProfile(wp)

	if err == nil {
		t.Fatal(err)
	}

	wp = WebProfile{}

	err = c.SetWebProfile(wp)

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

	err := c.SetWebProfile(wp)

	if err != nil {
		t.Fatal(err)
	}

}

func TestDeleteWebProfile_invalid(t *testing.T) {
	ts := httptest.NewServer(&webprofileTestServer{t: t})
	defer ts.Close()

	c, _ := NewClient("foo", "bar", ts.URL)

	err := c.DeleteWebProfile("foobar")

	if err == nil {
		t.Fatal(err)
	}

}
