<p>
    <a href="https://github.com/plutov/paypal/releases"><img src="https://img.shields.io/github/release/plutov/paypal.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/plutov/paypal?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
</p>

# Go client for PayPal REST API

## Paypal REST API Docs

- [Get started with PayPal REST APIs](https://developer.paypal.com/api/rest/)
- [Specification](https://github.com/paypal/paypal-rest-api-specifications)

## Missing endpoints

It is possible that some endpoints are missing in this client, but you can use built-in `paypal` functions to perform a request: `NewClient -> NewRequest -> SendWithAuth`

## Usage

```go
import "github.com/plutov/paypal/v4"

c, err := paypal.NewClient("clientID", "secretID", paypal.APIBaseSandBox) // or paypal.APIBaseLive
c.SetLog(os.Stdout)
```

### Get authorization

```go
auth, err := c.GetAuthorization("2DC87612EK520411B")
```

### Capture authorization

```go
capture, err := c.CaptureAuthorization(authID, &paypal.Amount{Total: "7.00", Currency: "USD"}, true)
```

### Void authorization

```go
auth, err := c.VoidAuthorization(authID)
```

### Reauthorize authorization

```go
auth, err := c.ReauthorizeAuthorization(authID, &paypal.Amount{Total: "7.00", Currency: "USD"})
```

### Get Refund

```go
refund, err := c.GetRefund("O-4J082351X3132253H")
```

### Get Order

```go
order, err := c.GetOrder("O-4J082351X3132253H")
```

### Create Order

```go
units := []paypal.PurchaseUnitRequest{}
source := &paypal.PaymentSource{}
appCtx := &paypal.ApplicationContext{}
order, err := c.CreateOrder(context.TODO(), paypal.OrderIntentCapture, units, ource, appCtx)
```

### Update Order

```go
order, err := c.UpdateOrder("O-4J082351X3132253H", []paypal.PurchaseUnitRequest{})
```

### Authorize Order

```go
auth, err := c.AuthorizeOrder(orderID, paypal.AuthorizeOrderRequest{})
```

### Capture Order

```go
capture, err := c.CaptureOrder(orderID, paypal.CaptureOrderRequest{})
```

### Identity

```go
token, err := c.GrantNewAccessTokenFromAuthCode("<Authorization-Code>", "http://example.com/myapp/return.php")

// ... or by refresh token
token, err := c.GrantNewAccessTokenFromRefreshToken("<Refresh-Token>")
```

### Get user info

```go
userInfo, err := c.GetUserInfo("openid")
```

### Create single payout to email

```go
payout := paypal.Payout{
    SenderBatchHeader: &paypal.SenderBatchHeader{
        EmailSubject: "Subject will be displayed on PayPal",
    },
    Items: []paypal.PayoutItem{
        paypal.PayoutItem{
            RecipientType: "EMAIL",
            Receiver:      "single-email-payout@mail.com",
            Amount: &paypal.AmountPayout{
                Value:    "15.11",
                Currency: "USD",
            },
            Note:         "Optional note",
            SenderItemID: "Optional Item ID",
        },
    },
}

payoutResp, err := c.CreatePayout(payout)
```

### Get payout

```go
payout, err := c.GetPayout("PayoutBatchID")
```

### Get payout item

```go
payoutItem, err := c.GetPayoutItem("PayoutItemID")
```

### Cancel unclaimed payout item

```go
payoutItem, err := c.CancelPayoutItem("PayoutItemID")
```

### Create web experience profile

```go
webprofile := WebProfile{
    Name: "YeowZa! T-Shirt Shop",
    Presentation: Presentation{
        BrandName:  "YeowZa! Paypal",
        LogoImage:  "http://www.yeowza.com",
        LocaleCode: "US",
    },

    InputFields: InputFields{
        AllowNote:       true,
        NoShipping:      NoShippingDisplay,
        AddressOverride: AddrOverrideFromCall,
    },

    FlowConfig: FlowConfig{
        LandingPageType:   LandingPageTypeBilling,
        BankTXNPendingURL: "http://www.yeowza.com",
    },
}

result, err := c.CreateWebProfile(webprofile)
```

### Get web experience profile

```go
webprofile, err := c.GetWebProfile("XP-CP6S-W9DY-96H8-MVN2")
```

### List web experience profile

```go
webprofiles, err := c.GetWebProfiles()
```

### Update web experience profile

```go
webprofile := WebProfile{
    ID: "XP-CP6S-W9DY-96H8-MVN2",
    Name: "Shop YeowZa! YeowZa! ",
}
err := c.SetWebProfile(webprofile)
```

### Delete web experience profile

```go
err := c.DeleteWebProfile("XP-CP6S-W9DY-96H8-MVN2")
```

### Vault

```go
c.StoreCreditCard(paypal.CreditCard{
    Number:      "4417119669820331",
    Type:        "visa",
    ExpireMonth: "11",
    ExpireYear:  "2020",
    CVV2:        "874",
    FirstName:   "Foo",
    LastName:    "Bar",
})

c.DeleteCreditCard("CARD-ID-123")

c.PatchCreditCard("CARD-ID-123", []paypal.CreditCardField{
    paypal.CreditCardField{
        Operation: "replace",
        Path:      "/billing_address/line1",
        Value:     "New value",
    },
})

c.GetCreditCard("CARD-ID-123")

c.GetCreditCards(nil)
```

### Webhooks

```go
c.CreateWebhook(paypal.CreateWebhookRequest{
    URL: "webhook URL",
    EventTypes: []paypal.WebhookEventType{
        paypal.WebhookEventType{
            Name: "PAYMENT.AUTHORIZATION.CREATED",
        },
    },
})

c.UpdateWebhook("WebhookID", []paypal.WebhookField{
    paypal.WebhookField{
        Operation: "replace",
        Path:      "/event_types",
        Value: []interface{}{
            map[string]interface{}{
                "name": "PAYMENT.SALE.REFUNDED",
            },
        },
    },
})

c.GetWebhook("WebhookID")

c.DeleteWebhook("WebhookID")

c.ListWebhooks(paypal.AncorTypeApplication)
```

### Generate Next Invoice Number

```go
c.GenerateInvoiceNumber(ctx) // might return something like "0001" or "0010".
```

### Get Invoice Details

```go
invoice, err := c.GetInvoiceDetails(ctx, "INV2-XFXV-YW42-ZANU-4F33")
```

## Contributing

Check out [./CONTRIBUTING.md](CONTRIBUTING.md).

## Testing

This projects uses [paypal-rest-api-specifications](https://github.com/paypal/paypal-rest-api-specifications) to generate a mock server for testing purposes.


```
make test
```

To update the mock server, first we need to get the latest OpenAPI specification files and then generate the mock server code. Paypal has multiple files that have to be combined into one OpenAPI file. The following commands will do that:

```bash
make oapi
```

Then make sure to add missing mock implementations in `mockserver/impl.go`.

## Linting

Make sure to run the linter before submitting a PR:

```bash
make lint
```
