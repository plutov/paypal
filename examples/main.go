package main

import (
	"paypalsdk"
	"strconv"

	"fmt"
	"os"
)

func main() {
	client, err := paypalsdk.NewClient("AZgwu4yt5Ba0gyTu1dGBH3txHCJbMuFNvrmQxBaQbfDncDiCs6W_rwJD8Ir-0pZrN-_eq7n9zVd8Y-5f", "EBzA1wRl5t73OMugOieDj_tI3vihfJmGl47ukQT-cpctooIzDu0K7IPESNC0cKodlLSOXzwI8qXSM0rd", paypalsdk.APIBaseSandBox)
	if err == nil {
		fmt.Println("DEBUG: ClientID=" + client.ClientID + " APIBase=" + client.APIBase)
	} else {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	token, err := client.GetAccessToken()
	if err == nil {
		fmt.Println("DEBUG: AccessToken=" + token.Token)
	} else {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	payment, err := client.GetPayment("PAY-TEST-123")
	fmt.Println("DEBUG: PaymentID=" + payment.ID)

	payments, err := client.GetPayments()
	if err == nil {
		fmt.Println("DEBUG: PaymentsCount=" + strconv.Itoa(len(payments)))
	} else {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
	fmt.Println("OK")
}
