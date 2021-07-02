package main

import (
	"testing"

	paypal "github.com/plutov/paypal/v4"
	"github.com/shawnhankim/paypal/mycodes/util"
)

func TestGetAccessToken(t *testing.T) {
	res, err := util.GetPaymentGWAccessToken(
		TestClientID, TestSecret, paypal.APIBaseSandBox,
	)
	if err != nil {
		panic(err)
	}
	util.DispPaymentGWAccessToken(res)
}
