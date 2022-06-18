package paypal

import (
	"context"
	"testing"
)

var testClientID = "AXy9orp-CDaHhBZ9C78QHW2BKZpACgroqo85_NIOa9mIfJ9QnSVKzY-X_rivR_fTUUr6aLjcJsj6sDur"
var testSecret = "EBoIiUSkCKeSk49hHSgTem1qnjzzJgRQHDEHvGpzlLEf_nIoJd91xu8rPOBDCdR_UYNKVxJE-UgS2iCw"

func TestUpdateOrder(t *testing.T) {
	ctx := context.Background()

	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	_, _ = c.GetAccessToken(ctx)

	orderResponse, err := c.CreateOrder(
		ctx,
		OrderIntentCapture,
		[]PurchaseUnitRequest{
			{
				Amount: &PurchaseUnitAmount{
					Value:    "7.00",
					Currency: "USD",
				},
			},
		},
		&CreateOrderPayer{},
		&ApplicationContext{},
	)
	if err != nil {
		t.Errorf("Not expected error for CreateOrder(), got %s", err.Error())
	}

	order, err := c.GetOrder(ctx, orderResponse.ID)
	if err != nil {
		t.Errorf("Not expected error for GetOrder(), got %s", err.Error())
	}

	if order.PurchaseUnits[0].Amount.Value != "7.00" {
		t.Errorf("CreateOrder amount incorrect")
	}

	err = c.UpdateOrder(
		ctx,
		orderResponse.ID,
		"replace",
		"/purchase_units/@reference_id=='default'/amount",
		map[string]string{
			"currency_code": "USD",
			"value":         "2.00",
		},
	)
	if err != nil {
		t.Errorf("Not expected error for UpdateOrder(), got %s", err.Error())
	}

	order, err = c.GetOrder(ctx, orderResponse.ID)
	if err != nil {
		t.Errorf("Not expected error for GetOrder(), got %s", err.Error())
	}

	if order.PurchaseUnits[0].Amount.Value != "2.00" {
		t.Errorf("CreateOrder after update amount incorrect")
	}
}
