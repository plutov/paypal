package paypalsdk_test

import (
	"fmt"
	pp "github.com/logpacker/PayPal-Go-SDK"
	"time"
)

func BillingExample() {
	plan := pp.BillingPlan{
		Name:        "Plan with Regular and Trial Payment Definitions",
		Description: "Plan with regular and trial payment definitions.",
		Type:        "fixed",
		PaymentDefinitions: []pp.PaymentDefinition{
			pp.PaymentDefinition{
				Name:              "Regular payment definition",
				Type:              "REGULAR",
				Frequency:         "MONTH",
				FrequencyInterval: "2",
				Amount: pp.AmountPayout{
					Value:    "100",
					Currency: "USD",
				},
				Cycles: "12",
				ChargeModels: []pp.ChargeModel{
					pp.ChargeModel{
						Type: "SHIPPING",
						Amount: pp.AmountPayout{
							Value:    "10",
							Currency: "USD",
						},
					},
					pp.ChargeModel{
						Type: "TAX",
						Amount: pp.AmountPayout{
							Value:    "12",
							Currency: "USD",
						},
					},
				},
			},
			pp.PaymentDefinition{
				Name:              "Trial payment definition",
				Type:              "trial",
				Frequency:         "week",
				FrequencyInterval: "5",
				Amount: pp.AmountPayout{
					Value:    "9.19",
					Currency: "USD",
				},
				Cycles: "2",
				ChargeModels: []pp.ChargeModel{
					pp.ChargeModel{
						Type: "SHIPPING",
						Amount: pp.AmountPayout{
							Value:    "1",
							Currency: "USD",
						},
					},
					pp.ChargeModel{
						Type: "TAX",
						Amount: pp.AmountPayout{
							Value:    "2",
							Currency: "USD",
						},
					},
				},
			},
		},
		MerchantPreferences: &pp.MerchantPreferences{
			SetupFee: &pp.AmountPayout{
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
	c, err := pp.NewClient("clientID", "secretID", pp.APIBaseSandBox)
	if err != nil {
		panic(err)
	}
	_, err = c.GetAccessToken()
	if err != nil {
		panic(err)
	}
	planResp, err := c.CreateBillingPlan(plan)
	if err != nil {
		panic(err)
	}
	err = c.ActivatePlan(planResp.ID)
	fmt.Println(err)
	agreement := pp.BillingAgreement{
		Name:        "Fast Speed Agreement",
		Description: "Agreement for Fast Speed Plan",
		StartDate:   pp.JsonTime(time.Now().Add(time.Hour * 24)),
		Plan:        pp.BillingPlan{ID: planResp.ID},
		Payer: pp.Payer{
			PaymentMethod: "paypal",
		},
	}
	resp, err := c.CreateBillingAgreement(agreement)
	fmt.Println(err, resp)
}
