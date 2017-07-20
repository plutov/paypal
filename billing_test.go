package paypalsdk

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonStructure(t *testing.T) {
	expected_str := `{
  "name": "Plan with Regular and Trial Payment Definitions",
  "description": "Plan with regular and trial payment definitions.",
  "type": "fixed",
  "payment_definitions": [
  {
    "name": "Regular payment definition",
    "type": "REGULAR",
    "frequency": "MONTH",
    "frequency_interval": "2",
    "amount":
    {
      "value": "100",
      "currency": "USD"
    },
    "cycles": "12",
    "charge_models": [
    {
      "type": "SHIPPING",
      "amount":
      {
        "value": "10",
        "currency": "USD"
      }
    },
    {
      "type": "TAX",
      "amount":
      {
        "value": "12",
        "currency": "USD"
      }
    }]
  },
  {
    "name": "Trial payment definition",
    "type": "trial",
    "frequency": "week",
    "frequency_interval": "5",
    "amount":
    {
      "value": "9.19",
      "currency": "USD"
    },
    "cycles": "2",
    "charge_models": [
    {
      "type": "SHIPPING",
      "amount":
      {
        "value": "1",
        "currency": "USD"
      }
    },
    {
      "type": "TAX",
      "amount":
      {
        "value": "2",
        "currency": "USD"
      }
    }]
  }],
  "merchant_preferences":
  {
    "setup_fee":
    {
      "value": "1",
      "currency": "USD"
    },
    "return_url": "http://www.paypal.com",
    "cancel_url": "http://www.paypal.com/cancel",
    "auto_bill_amount": "YES",
    "initial_fail_amount_action": "CONTINUE",
    "max_fail_attempts": "0"
  }
}`
	plan := BillingPlan{
		Name:        "Plan with Regular and Trial Payment Definitions",
		Description: "Plan with regular and trial payment definitions.",
		Type:        "fixed",
		PaymentDefinitions: []PaymentDefinition{
			PaymentDefinition{
				Name:              "Regular payment definition",
				Type:              "REGULAR",
				Frequency:         "MONTH",
				FrequencyInterval: "2",
				Amount: AmountPayout{
					Value:    "100",
					Currency: "USD",
				},
				Cycles: "12",
				ChargeModels: []ChargeModel{
					ChargeModel{
						Type: "SHIPPING",
						Amount: AmountPayout{
							Value:    "10",
							Currency: "USD",
						},
					},
					ChargeModel{
						Type: "TAX",
						Amount: AmountPayout{
							Value:    "12",
							Currency: "USD",
						},
					},
				},
			},
			PaymentDefinition{
				Name:              "Trial payment definition",
				Type:              "trial",
				Frequency:         "week",
				FrequencyInterval: "5",
				Amount: AmountPayout{
					Value:    "9.19",
					Currency: "USD",
				},
				Cycles: "2",
				ChargeModels: []ChargeModel{
					ChargeModel{
						Type: "SHIPPING",
						Amount: AmountPayout{
							Value:    "1",
							Currency: "USD",
						},
					},
					ChargeModel{
						Type: "TAX",
						Amount: AmountPayout{
							Value:    "2",
							Currency: "USD",
						},
					},
				},
			},
		},
		MerchantPreferences: MerchantPreferences{
			SetupFee: AmountPayout{
				Value:    "1",
				Currency: "USD",
			},
			ReturnUrl:               "http://www.paypal.com",
			CancelUrl:               "http://www.paypal.com/cancel",
			AutoBillAmount:          "YES",
			InitialFailAmountAction: "CONTINUE",
			MaxFailAttempts:         "0",
		},
	}
	_, err := json.Marshal(&plan)
	expected := new(BillingPlan)
	json.Unmarshal([]byte(expected_str), expected)
	assert.NoError(t, err)
	assert.Equal(t, expected, &plan, "Wrong billing builder")
}
