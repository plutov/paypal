[Docs](https://pkg.go.dev/github.com/plutov/paypal)

# Go client for PayPal REST API

## Coverage

### Auth

* POST /v1/oauth2/token

### /v1/payments

* POST /v1/payments/payouts
* GET /v1/payments/payouts/:id
* GET /v1/payments/payouts-item/:id
* POST /v1/payments/payouts-item/:id/cancel
* GET /v1/payments/sale/:id
* POST /v1/payments/sale/:id/refund
* GET /v1/payments/billing-plans
* POST /v1/payments/billing-plans
* PATCH /v1/payments/billing-plans/:id
* POST /v1/payments/billing-agreements
* POST /v1/payments/billing-agreements/:token/agreement-execute

### /v2/payments

* GET /v2/payments/authorizations/:id
* GET /v2/payments/captures/:id
* POST /v2/payments/authorizations/:id/capture
* POST /v2/payments/authorizations/:id/void
* POST /v2/payments/authorizations/:id/reauthorize
* GET /v2/payments/refund/:id

### Identity
* POST /v1/identity/openidconnect/tokenservice
* GET /v1/identity/openidconnect/userinfo/?schema=:schema

### /v1/payment-experience

* GET /v1/payment-experience/web-profiles
* POST /v1/payment-experience/web-profiles
* GET /v1/payment-experience/web-profiles/:id
* PUT /v1/payment-experience/web-profiles/:id
* DELETE /v1/payment-experience/web-profiles/:id

### /v1/reporting

* POST /v1/reporting/transactions

### Vault

* POST /v1/vault/credit-cards
* DELETE /v1/vault/credit-cards/:id
* PATCH /v1/vault/credit-cards/:id
* GET /v1/vault/credit-cards/:id
* GET /v1/vault/credit-cards

### Checkout

* POST /v2/checkout/orders
* GET /v2/checkout/orders/:id
* PATCH /v2/checkout/orders/:id
* POST /v2/checkout/orders/:id/authorize
* POST /v2/checkout/orders/:id/capture

### Notifications
* POST /v1/notifications/webhooks
* GET /v1/notifications/webhooks
* GET /v1/notifications/webhooks/:id
* PATCH /v1/notifications/webhooks/:id
* DELETE /v1/notifications/webhooks/:id
* POST /v1/notifications/verify-webhook-signature

### Products (Catalog)

* POST /v1/catalogs/products
* PATCH /v1/catalogs/products/:id
* GET /v1/catalogs/products/:id
* GET /v1/catalogs/products

### Billing Plans (Subscriptions)

* POST  /v1/billing/plans
* PATCH /v1/billing/plans/:id
* GET   /v1/billing/plans/:id
* GET   /v1/billing/plans
* POST  /v1/billing/plans/:id/activate
* POST  /v1/billing/plans/:id/deactivate
* POST  /v1/billing/plans/:id/update-pricing-schemes

### Subscriptions

* POST /v1/billing/subscriptions
* PATCH /v1/billing/subscriptions/:id
* GET /v1/billing/subscriptions/:id
* POST /v1/billing/subscriptions/:id/activate
* POST /v1/billing/subscriptions/:id/cancel
* POST /v1/billing/subscriptions/:id/revise
* POST /v1/billing/subscriptions/:id/capture
* POST /v1/billing/subscriptions/:id/suspend
* GET /v1/billing/subscriptions/:id/transactions

### Invoicing

* POST /v2/invoicing/generate-next-invoice-number
* GET /v2/invoicing/invoices/:id
 
## Missing endpoints

It is possible that some endpoints are missing in this Client, but you can use built-in `paypal` functions to perform a request: `NewClient -> NewRequest -> SendWithAuth`

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

### Get Sale by ID

```go
sale, err := c.GetSale("36C38912MN9658832")
```

### Refund Sale by ID

```go
// Full
refund, err := c.RefundSale(saleID, nil)
// Partial
refund, err := c.RefundSale(saleID, &paypal.Amount{Total: "7.00", Currency: "USD"})
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
order, err := c.CreateOrder(paypal.OrderIntentCapture, []paypal.PurchaseUnitRequest{paypal.PurchaseUnitRequest{ReferenceID: "ref-id", Amount: paypal.Amount{Total: "7.00", Currency: "USD"}}})
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

payoutResp, err := c.CreateSinglePayout(payout)
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
* for now, we are yet to implement the ShowAllInvoices endpoint, so use the following cURL request for the same(this gives you the list of invoice-IDs for this customer)
    ```bash
    curl -v -X GET https://api-m.sandbox.paypal.com/v2/invoicing/invoices?total_required=true \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <Token>" 
    ```

* refer to the beginning of this Usage section for obtaining a Token.
    

## How to Contribute

* Fork a repository
* Add/Fix something
* Check that tests are passing
* Create PR

Current contributors:

- [Roopak Venkatakrishnan](https://github.com/roopakv)
- [Alex Pliutau](https://github.com/plutov)

## Tests

* Unit tests: `go test -v ./...`
* Integration tests: `go test -tags=integration`
