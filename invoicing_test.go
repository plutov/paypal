package paypal_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/plutov/paypal/v4"
	"github.com/stretchr/testify/assert"
)

// All test values are defined here
var devTestClientID = "AXy9orp-CDaHhBZ9C78QHW2BKZpACgroqo85_NIOa9mIfJ9QnSVKzY-X_rivR_fTUUr6aLjcJsj6sDur"
var devTestSecret = "EBoIiUSkCKeSk49hHSgTem1qnjzzJgRQHDEHvGpzlLEf_nIoJd91xu8rPOBDCdR_UYNKVxJE-UgS2iCw"
var devAPIBaseSandBox = "https://api.sandbox.paypal.com"

func TestGenerateInvoiceNumber(t *testing.T) {
	ctx := context.Background()

	c, _ := paypal.NewClient(devTestClientID, devTestSecret, devAPIBaseSandBox)
	_, err := c.GetAccessToken(ctx)
	assert.Equal(t, nil, err)

	_, err = c.GenerateInvoiceNumber(ctx)
	assert.Equal(t, nil, err)
}
func assertTwoInvoices(t *testing.T, invoice paypal.Invoice, testInvoice paypal.Invoice) {

	// additional_recipients
	assert.Equal(t, len(invoice.AdditionalRecipients), len(testInvoice.AdditionalRecipients))
	// additional_recipients --> email_address !! EQUALITY OF SPLICE OF STRUCT REMAINING !!

	// amount
	assert.Equal(t, invoice.AmountSummary.Currency, testInvoice.AmountSummary.Currency)
	assert.Equal(t, invoice.AmountSummary.Value, testInvoice.AmountSummary.Value)
	// amount-->breakdown-->custom-->amount
	assert.Equal(t, invoice.AmountSummary.Breakdown.Custom.Amount.Currency, testInvoice.AmountSummary.Breakdown.Custom.Amount.Currency)
	assert.Equal(t, invoice.AmountSummary.Breakdown.Custom.Amount.Value, testInvoice.AmountSummary.Breakdown.Custom.Amount.Value)
	// amount-->breakdown-->custom-->label
	assert.Equal(t, invoice.AmountSummary.Breakdown.Custom.Label, testInvoice.AmountSummary.Breakdown.Custom.Label)
	// amount-->breakdown-->discount-->amount
	assert.Equal(t, invoice.AmountSummary.Breakdown.Discount.InvoiceDiscount.DiscountAmount.Currency, testInvoice.AmountSummary.Breakdown.Discount.InvoiceDiscount.DiscountAmount.Currency)
	assert.Equal(t, invoice.AmountSummary.Breakdown.Discount.InvoiceDiscount.DiscountAmount.Value, testInvoice.AmountSummary.Breakdown.Discount.InvoiceDiscount.DiscountAmount.Value)
	// amount-->breakdown-->discount-->percent
	assert.Equal(t, invoice.AmountSummary.Breakdown.Discount.InvoiceDiscount.Percent, testInvoice.AmountSummary.Breakdown.Discount.InvoiceDiscount.Percent)
	// amount-->breakdown-->discount-->item_discount
	assert.Equal(t, invoice.AmountSummary.Breakdown.Discount.ItemDiscount.Currency, testInvoice.AmountSummary.Breakdown.Discount.ItemDiscount.Currency)
	assert.Equal(t, invoice.AmountSummary.Breakdown.Discount.ItemDiscount.Value, testInvoice.AmountSummary.Breakdown.Discount.ItemDiscount.Value)
	// amount-->breakdown-->item_total
	assert.Equal(t, invoice.AmountSummary.Breakdown.ItemTotal.Currency, testInvoice.AmountSummary.Breakdown.ItemTotal.Currency)
	assert.Equal(t, invoice.AmountSummary.Breakdown.ItemTotal.Value, testInvoice.AmountSummary.Breakdown.ItemTotal.Value)
	// amount-->breakdown-->shipping-->amount
	assert.Equal(t, invoice.AmountSummary.Breakdown.Shipping.Amount.Currency, testInvoice.AmountSummary.Breakdown.Shipping.Amount.Currency)
	assert.Equal(t, invoice.AmountSummary.Breakdown.Shipping.Amount.Value, testInvoice.AmountSummary.Breakdown.Shipping.Amount.Value)
	// amount-->breakdown-->shipping-->tax
	assert.Equal(t, invoice.AmountSummary.Breakdown.Shipping.Tax.Amount.Currency, testInvoice.AmountSummary.Breakdown.Shipping.Tax.Amount.Currency)
	assert.Equal(t, invoice.AmountSummary.Breakdown.Shipping.Tax.Amount.Value, testInvoice.AmountSummary.Breakdown.Shipping.Tax.Amount.Value)
	assert.Equal(t, invoice.AmountSummary.Breakdown.Shipping.Tax.ID, testInvoice.AmountSummary.Breakdown.Shipping.Tax.ID)
	assert.Equal(t, invoice.AmountSummary.Breakdown.Shipping.Tax.Name, testInvoice.AmountSummary.Breakdown.Shipping.Tax.Name)
	assert.Equal(t, invoice.AmountSummary.Breakdown.Shipping.Tax.Percent, testInvoice.AmountSummary.Breakdown.Shipping.Tax.Percent)
	// amount-->breakdown-->tax_total
	assert.Equal(t, invoice.AmountSummary.Breakdown.TaxTotal.Currency, testInvoice.AmountSummary.Breakdown.TaxTotal.Currency)
	assert.Equal(t, invoice.AmountSummary.Breakdown.TaxTotal.Value, testInvoice.AmountSummary.Breakdown.TaxTotal.Value)

	// configuration
	assert.Equal(t, invoice.Configuration.AllowTip, testInvoice.Configuration.AllowTip)
	assert.Equal(t, invoice.Configuration.TaxCalculatedAfterDiscount, testInvoice.Configuration.TaxCalculatedAfterDiscount)
	assert.Equal(t, invoice.Configuration.TaxInclusive, testInvoice.Configuration.TaxInclusive)
	assert.Equal(t, invoice.Configuration.TemplateId, testInvoice.Configuration.TemplateId)
	// configuration --> partial_payment
	assert.Equal(t, invoice.Configuration.PartialPayment.AllowPartialPayment, testInvoice.Configuration.PartialPayment.AllowPartialPayment)
	assert.Equal(t, invoice.Configuration.PartialPayment.MinimumAmountDue.Currency, testInvoice.Configuration.PartialPayment.MinimumAmountDue.Currency)
	assert.Equal(t, invoice.Configuration.PartialPayment.MinimumAmountDue.Value, testInvoice.Configuration.PartialPayment.MinimumAmountDue.Value)

	// detail
	assert.Equal(t, invoice.Detail.CurrencyCode, testInvoice.Detail.CurrencyCode)
	assert.Equal(t, invoice.Detail.InvoiceDate, testInvoice.Detail.InvoiceDate)
	assert.Equal(t, invoice.Detail.InvoiceNumber, testInvoice.Detail.InvoiceNumber)
	assert.Equal(t, invoice.Detail.Memo, testInvoice.Detail.Memo)
	assert.Equal(t, invoice.Detail.Note, testInvoice.Detail.Note)
	assert.Equal(t, invoice.Detail.Reference, testInvoice.Detail.Reference)
	assert.Equal(t, invoice.Detail.TermsAndConditions, testInvoice.Detail.TermsAndConditions)
	// detail --> attachments !! EQUALITY OF SPLICE OF STRUCT REMAINING !!
	assert.Equal(t, len(invoice.Detail.Attachments), len(testInvoice.Detail.Attachments))
	// detail --> metadata
	assert.Equal(t, invoice.Detail.Metadata.CancelTime, testInvoice.Detail.Metadata.CancelTime)
	assert.Equal(t, invoice.Detail.Metadata.CancellledTimeBy, testInvoice.Detail.Metadata.CancellledTimeBy)
	assert.Equal(t, invoice.Detail.Metadata.CreateTime, testInvoice.Detail.Metadata.CreateTime)
	assert.Equal(t, invoice.Detail.Metadata.CreatedBy, testInvoice.Detail.Metadata.CreatedBy)
	assert.Equal(t, invoice.Detail.Metadata.CreatedByFlow, testInvoice.Detail.Metadata.CreatedByFlow)
	assert.Equal(t, invoice.Detail.Metadata.FirstSentTime, testInvoice.Detail.Metadata.FirstSentTime)
	assert.Equal(t, invoice.Detail.Metadata.InvoicerViewUrl, testInvoice.Detail.Metadata.InvoicerViewUrl)
	assert.Equal(t, invoice.Detail.Metadata.LastSentBy, testInvoice.Detail.Metadata.LastSentBy)
	assert.Equal(t, invoice.Detail.Metadata.LastSentTime, testInvoice.Detail.Metadata.LastSentTime)
	assert.Equal(t, invoice.Detail.Metadata.LastUpdateTime, testInvoice.Detail.Metadata.LastUpdateTime)
	assert.Equal(t, invoice.Detail.Metadata.LastUpdatedBy, testInvoice.Detail.Metadata.LastUpdatedBy)
	assert.Equal(t, invoice.Detail.Metadata.RecipientViewUrl, testInvoice.Detail.Metadata.RecipientViewUrl)
	// detail --> payment_term
	assert.Equal(t, invoice.Detail.PaymentTerm.DueDate, testInvoice.Detail.PaymentTerm.DueDate)
	assert.Equal(t, invoice.Detail.PaymentTerm.TermType, testInvoice.Detail.PaymentTerm.TermType)

	// due_amount
	assert.Equal(t, invoice.DueAmount.Currency, testInvoice.DueAmount.Currency)
	assert.Equal(t, invoice.DueAmount.Value, testInvoice.DueAmount.Value)

	// gratuity
	assert.Equal(t, invoice.Gratuity.Currency, testInvoice.Gratuity.Currency)
	assert.Equal(t, invoice.Gratuity.Value, testInvoice.Gratuity.Value)

	// id
	assert.Equal(t, invoice.ID, testInvoice.ID)

	// invoicer
	assert.Equal(t, invoice.Invoicer.AdditionalNotes, testInvoice.Invoicer.AdditionalNotes)
	assert.Equal(t, invoice.Invoicer.EmailAddress, testInvoice.Invoicer.EmailAddress)
	assert.Equal(t, invoice.Invoicer.LogoUrl, testInvoice.Invoicer.LogoUrl)
	assert.Equal(t, invoice.Invoicer.TaxId, testInvoice.Invoicer.TaxId)
	assert.Equal(t, invoice.Invoicer.Website, testInvoice.Invoicer.Website)
	// !!!  SPLICE EQUALITY STILL REMAINING !!!!!
	// invoicer --> phones
	assert.Equal(t, len(invoice.Invoicer.Phones), len(testInvoice.Invoicer.Phones))

	// items
	// !!!  SPLICE EQUALITY STILL REMAINING !!!!!
	assert.Equal(t, len(invoice.Items), len(testInvoice.Items))

	// links
	// !!!  SPLICE EQUALITY STILL REMAINING !!!!!
	assert.Equal(t, len(invoice.Links), len(testInvoice.Links))

	// parent_id
	assert.Equal(t, invoice.ParentID, testInvoice.ParentID)

	// payments
	assert.Equal(t, invoice.Payments.PaidAmount.Currency, testInvoice.Payments.PaidAmount.Currency)
	assert.Equal(t, invoice.Payments.PaidAmount.Value, testInvoice.Payments.PaidAmount.Value)
	// payments --> transactions
	assert.Equal(t, len(invoice.Payments.Transactions), len(testInvoice.Payments.Transactions))

	// primary_recipients
	// !!!  SPLICE EQUALITY STILL REMAINING !!!!!
	assert.Equal(t, len(invoice.PrimaryRecipients), len(testInvoice.PrimaryRecipients))

	// refunds
	assert.Equal(t, invoice.Refunds.RefundAmount.Currency, testInvoice.Refunds.RefundAmount.Currency)
	assert.Equal(t, invoice.Refunds.RefundAmount.Value, testInvoice.Refunds.RefundAmount.Value)
	assert.Equal(t, len(invoice.Refunds.RefundDetails), len(testInvoice.Refunds.RefundDetails))

	// status
	assert.Equal(t, invoice.Status, testInvoice.Status)

}

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
				"id": "TAX-9AU06895VD287170A",
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

	ctx := context.Background()

	c, _ := paypal.NewClient(devTestClientID, devTestSecret, devAPIBaseSandBox)
	_, _ = c.GetAccessToken(ctx)

	invoice, err := c.GetInvoiceDetails(ctx, "INV2-XFXV-YW42-ZANU-4F33")
	assert.Equal(t, nil, err) // if passed, means that request was successful
	assertTwoInvoices(t, *invoice, testInvoice)
}
