package paypalsdk

import "testing"

func TestCreateSinglePayout(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	c.GetAccessToken()

	payout := Payout{
		SenderBatchHeader: &SenderBatchHeader{
			EmailSubject: "Subject will be displayed on PayPal",
		},
		Items: []PayoutItem{
			{
				RecipientType: "EMAIL",
				Receiver:      "single-email-payout@mail.com",
				Amount: &AmountPayout{
					Value:    "15.11",
					Currency: "USD",
				},
				Note:         "Optional note",
				SenderItemID: "Optional Item ID",
			},
		},
	}

	p, err := c.CreateSinglePayout(payout)

	if err != nil || len(p.Items) != 1 {
		t.Errorf("Test single payout is not created")
	}
}
