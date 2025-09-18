package paypal_test

import (
	"encoding/json"
	"testing"

	"github.com/plutov/paypal/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetInvoice(t *testing.T) {
	testInvoiceJSONData := []byte(`
	{
		"amount": {
		  "breakdown": {
			"custom": {
			  "amount": {
				"currency_code": "USD",
				"value": "10.00"
			  },
			  "label": "Packing Charges"
			},
			"discount": {
			  "invoice_discount": {
				"amount": {
				  "currency_code": "USD",
				  "value": "-2.63"
				},
				"percent": "5"
			  },
			  "item_discount": {
				"currency_code": "USD",
				"value": "-7.50"
			  }
			},
			"item_total": {
			  "currency_code": "USD",
			  "value": "60.00"
			},
			"shipping": {
			  "amount": {
				"currency_code": "USD",
				"value": "10.00"
			  },
			  "tax": {
				"amount": {
				  "currency_code": "USD",
				  "value": "0.73"
				},
				"id": "TAX-5XV24702TP4910056",
				"name": "Sales Tax",
				"percent": "7.25"
			  }
			},
			"tax_total": {
			  "currency_code": "USD",
			  "value": "4.34"
			}
		  },
		  "currency_code": "USD",
		  "value": "74.21"
		},
		"configuration": {
		  "allow_tip": true,
		  "partial_payment": {
			"allow_partial_payment": true,
			"minimum_amount_due": {
			  "currency_code": "USD",
			  "value": "20.00"
			}
		  },
		  "tax_calculated_after_discount": true,
		  "tax_inclusive": false,
		  "template_id": "TEMP-4NW98229SC0703920"
		},
		"detail": {
		  "additional_data": "2-4",
		  "archived": false,
		  "category_code": "SHIPPABLE",
		  "currency_code": "USD",
		  "group_draft": false,
		  "invoice_date": "2018-11-12",
		  "invoice_number": "0001",
		  "memo": "This is a long contract",
		  "metadata": {
			"caller_type": "API_V2_INVOICE",
			"create_time": "2022-10-25T16:54:50Z",
			"created_by_flow": "REGULAR_SINGLE",
			"invoicer_view_url": "https://www.sandbox.paypal.com/invoice/details/INV2-XFXV-YW42-ZANU-4F33",
			"last_update_time": "2022-10-25T16:54:50Z",
			"recipient_view_url": "https://www.sandbox.paypal.com/invoice/p/#XFXVYW42ZANU4F33"
		  },
		  "note": "Thank you for your business.",
		  "payment_term": {
			"due_date": "2018-11-22",
			"term_type": "NET_10"
		  },
		  "reference": "deal-ref",
		  "viewed_by_recipient": false
		},
		"due_amount": {
		  "currency_code": "USD",
		  "value": "74.21"
		},
		"id": "INV2-XFXV-YW42-ZANU-4F33",
		"invoicer": {
		  "additional_notes": "2-4",
		  "address": {
			"address_line_1": "1234 First Street",
			"address_line_2": "337673 Hillside Court",
			"admin_area_1": "CA",
			"admin_area_2": "Anytown",
			"country_code": "US",
			"postal_code": "98765"
		  },
		  "email_address": "merchant@example.com",
		  "logo_url": "https://example.com/logo.PNG",
		  "name": {
			"full_name": "David Larusso",
			"given_name": "David",
			"surname": "Larusso"
		  },
		  "phones": [
			{
			  "country_code": "001",
			  "national_number": "4085551234",
			  "phone_type": "MOBILE"
			}
		  ],
		  "tax_id": "ABcNkWSfb5ICTt73nD3QON1fnnpgNKBy- Jb5SeuGj185MNNw6g",
		  "website": "www.test.com"
		},
		"items": [
		  {
			"description": "Elastic mat to practice yoga.",
			"discount": {
			  "amount": {
				"currency_code": "USD",
				"value": "-2.50"
			  },
			  "percent": "5"
			},
			"id": "ITEM-5335764681676603X",
			"name": "Yoga Mat",
			"quantity": "1",
			"tax": {
			  "amount": {
				"currency_code": "USD",
				"value": "3.27"
			  },
			  "id": "TAX-5XV24702TP4910056",
			  "name": "Sales Tax",
			  "percent": "7.25"
			},
			"unit_amount": {
			  "currency_code": "USD",
			  "value": "50.00"
			},
			"unit_of_measure": "QUANTITY"
		  },
		  {
			"discount": {
			  "amount": {
				"currency_code": "USD",
				"value": "-5.00"
			  }
			},
			"id": "ITEM-1B467958Y9218273X",
			"name": "Yoga t-shirt",
			"quantity": "1",
			"tax": {
			  "amount": {
				"currency_code": "USD",
				"value": "0.34"
			  },
			  "id": "TAX-5XV24702TP4910056",
			  "name": "Sales Tax",
			  "percent": "7.25"
			},
			"unit_amount": {
			  "currency_code": "USD",
			  "value": "10.00"
			},
			"unit_of_measure": "QUANTITY"
		  }
		],
		"links": [
		  {
			"href": "https://api.sandbox.paypal.com/v2/invoicing/invoices/INV2-XFXV-YW42-ZANU-4F33",
			"method": "GET",
			"rel": "self"
		  },
		  {
			"href": "https://api.sandbox.paypal.com/v2/invoicing/invoices/INV2-XFXV-YW42-ZANU-4F33/send",
			"method": "POST",
			"rel": "send"
		  },
		  {
			"href": "https://api.sandbox.paypal.com/v2/invoicing/invoices/INV2-XFXV-YW42-ZANU-4F33",
			"method": "PUT",
			"rel": "replace"
		  },
		  {
			"href": "https://api.sandbox.paypal.com/v2/invoicing/invoices/INV2-XFXV-YW42-ZANU-4F33",
			"method": "DELETE",
			"rel": "delete"
		  },
		  {
			"href": "https://api.sandbox.paypal.com/v2/invoicing/invoices/INV2-XFXV-YW42-ZANU-4F33/payments",
			"method": "POST",
			"rel": "record-payment"
		  }
		],
		"primary_recipients": [
		  {
			"billing_info": {
			  "address": {
				"address_line_1": "1234 Main Street",
				"admin_area_1": "CA",
				"admin_area_2": "Anytown",
				"country_code": "US",
				"postal_code": "98765"
			  },
			  "email_address": "bill-me@example.com",
			  "name": {
				"full_name": "Stephanie Meyers",
				"given_name": "Stephanie",
				"surname": "Meyers"
			  }
			},
			"shipping_info": {
			  "address": {
				"address_line_1": "1234 Main Street",
				"admin_area_1": "CA",
				"admin_area_2": "Anytown",
				"country_code": "US",
				"postal_code": "98765"
			  },
			  "name": {
				"full_name": "Stephanie Meyers",
				"given_name": "Stephanie",
				"surname": "Meyers"
			  }
			}
		  }
		],
		"status": "DRAFT",
		"unilateral": false
	  }
	`)
	var testInvoice paypal.Invoice
	err := json.Unmarshal(testInvoiceJSONData, &testInvoice)
	assert.Equal(t, nil, err) // if passed, means unmarshalling was successful
}
