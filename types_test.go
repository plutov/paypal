package paypalsdk

// These tests test responses convertion from JSON to golang structs

import (
	"encoding/json"
	"testing"
)

func TestTypeUserInfo(t *testing.T) {
	response := `{
    "user_id": "https://www.paypal.com/webapps/auth/server/64ghr894040044",
    "name": "Peter Pepper",
    "given_name": "Peter",
    "family_name": "Pepper",
    "email": "ppuser@example.com"
    }`

	u := &UserInfo{}
	err := json.Unmarshal([]byte(response), u)
	if err != nil {
		t.Errorf("UserInfo Unmarshal failed")
	}

	if u.ID != "https://www.paypal.com/webapps/auth/server/64ghr894040044" ||
		u.Name != "Peter Pepper" ||
		u.GivenName != "Peter" ||
		u.FamilyName != "Pepper" ||
		u.Email != "ppuser@example.com" {
		t.Errorf("UserInfo decoded result is incorrect")
	}
}
