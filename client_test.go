package paypalsdk

import (
	"fmt"
	"testing"
)

// All test values are defined here
var testClientID = "AZgwu4yt5Ba0gyTu1dGBH3txHCJbMuFNvrmQxBaQbfDncDiCs6W_rwJD8Ir-0pZrN-_eq7n9zVd8Y-5f"
var testSecret = "EBzA1wRl5t73OMugOieDj_tI3vihfJmGl47ukQT-cpctooIzDu0K7IPESNC0cKodlLSOXzwI8qXSM0rd"
var testAuthID = "2DC87612EK520411B"
var testFakeAuthID = "FAKE-2DC87612EK520411B"
var testOrderID = "O-4J082351X3132253H"
var testFakeOrderID = "FAKE-O-4J082351X3132253H"
var testSaleID = "4CF18861HF410323U"
var testPaymentID = "PAY-5YK922393D847794YKER7MUI"
var testPayerID = "CR87QHB7JTRSC"
var testUserID = "https://www.paypal.com/webapps/auth/identity/user/WEssgRpQij92sE99_F9MImvQ8FPYgUEjrvCja2qH2H8"

func TestNewClient(t *testing.T) {
	_, err := NewClient("", "", "")
	if err == nil {
		t.Errorf("All arguments are required in NewClient()")
	} else {
		fmt.Println(err.Error())
	}

	_, err = NewClient(testClientID, testSecret, APIBaseSandBox)
	if err != nil {
		t.Errorf("NewClient() must not return error for valid creds: " + err.Error())
	}
}
