package paypal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreateOrderPaypalPayerInformation(t *testing.T) {
	for _, tt := range []struct {
		name   string
		source PaymentSourcePaypal
		want   string
	}{
		{
			name: "payer information",
			source: PaymentSourcePaypal{
				EmailAddress: "buyer@example.com",
				Name:         &CreateOrderPayerName{GivenName: "Jane", Surname: "Doe"},
				Phone: &PhoneWithType{
					PhoneType:   "MOBILE",
					PhoneNumber: &PhoneWithTypeNumber{NationalNumber: "2125551234"},
				},
				BirthDate: "1990-01-02",
				TaxInfo:   &TaxInfo{TaxID: "12345678901", TaxIDType: "BR_CPF"},
				Address: &ShippingDetailAddressPortable{
					AddressLine1: "123 Main St",
					AddressLine2: "Apt 4",
					AdminArea1:   "NY",
					AdminArea2:   "New York",
					PostalCode:   "10001",
					CountryCode:  "US",
				},
			},
			want: `{
				"email_address":"buyer@example.com",
				"name":{"given_name":"Jane","surname":"Doe"},
				"phone":{"phone_type":"MOBILE","phone_number":{"national_number":"2125551234"}},
				"birth_date":"1990-01-02",
				"tax_info":{"tax_id":"12345678901","tax_id_type":"BR_CPF"},
				"address":{"address_line_1":"123 Main St","address_line_2":"Apt 4","admin_area_1":"NY","admin_area_2":"New York","postal_code":"10001","country_code":"US"}
			}`,
		},
		{
			name: "optional payer information omitted",
			want: `{}`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt.source.ExperienceContext = PaymentSourcePaypalExperienceContext{
				PaymentMethodPreference: "IMMEDIATE_PAYMENT_REQUIRED",
				BrandName:               "Example Store",
				Locale:                  "en-US",
				LandingPage:             "LOGIN",
				ShippingPreference:      "NO_SHIPPING",
				UserAction:              "PAY_NOW",
				ReturnURL:               "https://example.com/return",
				CancelURL:               "https://example.com/cancel",
			}

			var want map[string]any
			if err := json.Unmarshal([]byte(tt.want), &want); err != nil {
				t.Fatal(err)
			}
			want["experience_context"] = map[string]any{
				"payment_method_preference": "IMMEDIATE_PAYMENT_REQUIRED",
				"brand_name":                "Example Store",
				"locale":                    "en-US",
				"landing_page":              "LOGIN",
				"shipping_preference":       "NO_SHIPPING",
				"user_action":               "PAY_NOW",
				"return_url":                "https://example.com/return",
				"cancel_url":                "https://example.com/cancel",
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/v2/checkout/orders" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
				}
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				paymentSource, ok := body["payment_source"].(map[string]any)
				if !ok || !reflect.DeepEqual(paymentSource["paypal"], want) {
					t.Errorf("payment_source = %#v, want paypal = %#v", body["payment_source"], want)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id":"ORDER123","status":"CREATED"}`))
			}))
			defer server.Close()

			client, err := NewClient("client-id", "secret", server.URL)
			if err != nil {
				t.Fatal(err)
			}
			client.SetAccessToken("token")
			order, err := client.CreateOrder(context.Background(), "CAPTURE", []PurchaseUnitRequest{
				{Amount: &PurchaseUnitAmount{Currency: "USD", Value: "10.00"}},
			}, &PaymentSource{Paypal: &tt.source}, nil)
			if err != nil {
				t.Fatal(err)
			}
			if order.ID != "ORDER123" {
				t.Errorf("order ID = %q, want ORDER123", order.ID)
			}
		})
	}
}
