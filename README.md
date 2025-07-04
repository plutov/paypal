[Docs](https://pkg.go.dev/github.com/plutov/paypal)

<p>
    <a href="https://github.com/plutov/paypal/releases"><img src="https://img.shields.io/github/release/plutov/paypal.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/plutov/paypal?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
</p>

![CodeRabbit Pull Request Reviews](https://img.shields.io/coderabbit/prs/github/plutov/paypal?utm_source=oss&utm_medium=github&utm_campaign=plutov%2Fpaypal&labelColor=171717&color=FF570A&link=https%3A%2F%2Fcoderabbit.ai&label=CodeRabbit+Reviews)

# Go client for PayPal REST API

## Paypal REST API Docs

[Get started with PayPal REST APIs](https://developer.paypal.com/api/rest/)

## Missing endpoints

It is possible that some endpoints are missing in this client, but you can use built-in `paypal` functions to perform a request: `NewClient -> NewRequest -> SendWithAuth`

## Usage

```go
import "github.com/plutov/paypal/v4"

// Create a client instance
c, err := paypal.NewClient("clientID", "secretID", paypal.APIBaseSandBox)
c.SetLog(os.Stdout) // Set log to terminal stdout
```

## Get authorization by ID

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

### Get Refund by ID

```go
refund, err := c.GetRefund("O-4J082351X3132253H")
```

### Get Order by ID

```go
order, err := c.GetOrder("O-4J082351X3132253H")
```

### Create an Order

```go
ctx := context.Background()
units := []paypal.PurchaseUnitRequest{}
source := &paypal.PaymentSource{}
appCtx := &paypal.ApplicationContext{}
order, err := c.CreateOrder(ctx, paypal.OrderIntentCapture, units, ource, appCtx)
```

### Update Order by ID

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

### Retreive user information

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

### Get payout by ID

```go
payout, err := c.GetPayout("PayoutBatchID")
```

### Get payout item by ID

```go
payoutItem, err := c.GetPayoutItem("PayoutItemID")
```

### Cancel unclaimed payout item by ID

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
// Store CC
c.StoreCreditCard(paypal.CreditCard{
    Number:      "4417119669820331",
    Type:        "visa",
    ExpireMonth: "11",
    ExpireYear:  "2020",
    CVV2:        "874",
    FirstName:   "Foo",
    LastName:    "Bar",
})

// Delete it
c.DeleteCreditCard("CARD-ID-123")

// Edit it
c.PatchCreditCard("CARD-ID-123", []paypal.CreditCardField{
    paypal.CreditCardField{
        Operation: "replace",
        Path:      "/billing_address/line1",
        Value:     "New value",
    },
})

// Get it
c.GetCreditCard("CARD-ID-123")

// Get all stored credit cards
c.GetCreditCards(nil)
```

### Webhooks

```go
// Create a webhook
c.CreateWebhook(paypal.CreateWebhookRequest{
    URL: "webhook URL",
    EventTypes: []paypal.WebhookEventType{
        paypal.WebhookEventType{
            Name: "PAYMENT.AUTHORIZATION.CREATED",
        },
    },
})

// Update a registered webhook
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

// Get a registered webhook
c.GetWebhook("WebhookID")

// Delete a webhook
c.DeleteWebhook("WebhookID")

// List registered webhooks
c.ListWebhooks(paypal.AncorTypeApplication)
```

### Generate Next Invoice Number

```go
// GenerateInvoiceNumber: generates the next invoice number that is available to the merchant.
c.GenerateInvoiceNumber(ctx) // might return something like "0001" or "0010".
```

### Get Invoice Details by ID

```go
// the second argument is an ID, it should be valid
invoice, err := c.GetInvoiceDetails(ctx, "INV2-XFXV-YW42-ZANU-4F33")
```

- for now, we are yet to implement the ShowAllInvoices endpoint, so use the following cURL request for the same(this gives you the list of invoice-IDs for this customer)

  ```bash
  curl -v -X GET https://api-m.sandbox.paypal.com/v2/invoicing/invoices?total_required=true \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <Token>"
  ```

- refer to the beginning of this Usage section for obtaining a Token.

## How to Contribute

- Fork a repository
- Add/Fix something
- Check that tests are passing
- Create PR

Main contributors:

- [Alex Pliutau](https://github.com/plutov)
- [Roopak Venkatakrishnan](https://github.com/roopakv)

## Tests

```
go test -v ./...
```
