package paypal_test

import (
	"fmt"
	"time"

	paypal "github.com/plutov/paypal"
)

func BillingExample() {
	plan := paypal.BillingPlan{
		Name:        "Plan with Regular and Trial Payment Definitions",
		Description: "Plan with regular and trial payment definitions.",
		Type:        "fixed",
		PaymentDefinitions: []paypal.PaymentDefinition{
			{
				Name:              "Regular payment definition",
				Type:              "REGULAR",
				Frequency:         "MONTH",
				FrequencyInterval: "2",
				Amount: paypal.AmountPayout{
					Value:    "100",
					Currency: "USD",
				},
				Cycles: "12",
				ChargeModels: []paypal.ChargeModel{
					{
						Type: "SHIPPING",
						Amount: paypal.AmountPayout{
							Value:    "10",
							Currency: "USD",
						},
					},
					{
						Type: "TAX",
						Amount: paypal.AmountPayout{
							Value:    "12",
							Currency: "USD",
						},
					},
				},
			},
			{
				Name:              "Trial payment definition",
				Type:              "trial",
				Frequency:         "week",
				FrequencyInterval: "5",
				Amount: paypal.AmountPayout{
					Value:    "9.19",
					Currency: "USD",
				},
				Cycles: "2",
				ChargeModels: []paypal.ChargeModel{
					{
						Type: "SHIPPING",
						Amount: paypal.AmountPayout{
							Value:    "1",
							Currency: "USD",
						},
					},
					{
						Type: "TAX",
						Amount: paypal.AmountPayout{
							Value:    "2",
							Currency: "USD",
						},
					},
				},
			},
		},
		MerchantPreferences: &paypal.MerchantPreferences{
			SetupFee: &paypal.AmountPayout{
				Value:    "1",
				Currency: "USD",
			},
			ReturnURL:               "http://www.paypal.com",
			CancelURL:               "http://www.paypal.com/cancel",
			AutoBillAmount:          "YES",
			InitialFailAmountAction: "CONTINUE",
			MaxFailAttempts:         "0",
		},
	}
	c, err := paypal.NewClient("clientID", "secretID", paypal.APIBaseSandBox)
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
	agreement := paypal.BillingAgreement{
		Name:        "Fast Speed Agreement",
		Description: "Agreement for Fast Speed Plan",
		StartDate:   paypal.JSONTime(time.Now().Add(time.Hour * 24)),
		Plan:        paypal.BillingPlan{ID: planResp.ID},
		Payer: paypal.Payer{
			PaymentMethod: "paypal",
		},
	}
	resp, err := c.CreateBillingAgreement(agreement)
	fmt.Println(err, resp)

	bps, err := c.ListBillingPlans(paypal.BillingPlanListParams{Status: "ACTIVE"})
	fmt.Println(err, bps)
}
