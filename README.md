[![Go Report Card](https://goreportcard.com/badge/logpacker/PayPal-Go-SDK)](https://goreportcard.com/report/logpacker/PayPal-Go-SDK)
[![Build Status](https://travis-ci.org/logpacker/PayPal-Go-SDK.svg?branch=master)](https://travis-ci.org/logpacker/PayPal-Go-SDK)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/logpacker/PayPal-Go-SDK)
[![Chat at https://gitter.im/logpacker/PayPal-Go-SDK](https://img.shields.io/badge/gitter-dev_chat-46bc99.svg)](https://gitter.im/logpacker/PayPal-Go-SDK?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Sourcegraph](https://sourcegraph.com/github.com/logpacker/PayPal-Go-SDK/-/badge.svg)](https://sourcegraph.com/github.com/logpacker/PayPal-Go-SDK?badge)

### Go client for PayPal REST API

Currently supports **v2** only, if you want to use **v1**, use **v1.1.4** git tag.

### Coverage

 * POST /v1/oauth2/token
 * POST /v1/identity/openidconnect/tokenservice
 * GET /v1/identity/openidconnect/userinfo/?schema=**SCHEMA**
 * POST /v1/payments/payouts
 * GET /v1/payments/payouts/**ID**
 * GET /v1/payments/payouts-item/**ID**
 * POST /v1/payments/payouts-item/**ID**/cancel
 * GET /v1/payment-experience/web-profiles
 * POST /v1/payment-experience/web-profiles
 * GET /v1/payment-experience/web-profiles/**ID**
 * PUT /v1/payment-experience/web-profiles/**ID**
 * DELETE /v1/payment-experience/web-profiles/**ID**
 * POST /v1/vault/credit-cards
 * DELETE /v1/vault/credit-cards/**ID**
 * PATCH /v1/vault/credit-cards/**ID**
 * GET /v1/vault/credit-cards/**ID**
 * GET /v1/vault/credit-cards
 * GET /v2/payments/authorization/**ID**
 * POST /v2/payments/authorization/**ID**/capture
 * POST /v2/payments/authorization/**ID**/void
 * POST /v2/payments/authorization/**ID**/reauthorize
 * GET /v2/payments/sale/**ID**
 * POST /v2/payments/sale/**ID**/refund
 * GET /v2/payments/refund/**ID**
 * POST /v2/checkout/orders
 * GET /v2/checkout/orders/**ID**
 * POST /v2/checkout/orders/**ID**/authorize
 * POST /v2/checkout/orders/**ID**/capture
 * POST /v2/payments/billing-plans
 * PATCH /v2/payments/billing-plans/***ID***
 * POST /v2/payments/billing-agreements
 * POST /v2/payments/billing-agreements/***TOKEN***/agreement-execute

### Missing endpoints
It is possible that some endpoints are missing in this SDK Client, but you can use built-in **paypalsdk** functions to perform a request: **NewClient -> NewRequest -> SendWithAuth**

### New Client

```go
import "github.com/logpacker/PayPal-Go-SDK"

// Create a client instance
c, err := paypalsdk.NewClient("clientID", "secretID", paypalsdk.APIBaseSandBox)
c.SetLog(os.Stdout) // Set log to terminal stdout

accessToken, err := c.GetAccessToken()
```

### Get authorization by ID

```go
auth, err := c.GetAuthorization("2DC87612EK520411B")
```

### Capture authorization

```go
capture, err := c.CaptureAuthorization(authID, &paypalsdk.Amount{Total: "7.00", Currency: "USD"}, true)
```

### Void authorization

```go
auth, err := c.VoidAuthorization(authID)
```

### Reauthorize authorization

```go
auth, err := c.ReauthorizeAuthorization(authID, &paypalsdk.Amount{Total: "7.00", Currency: "USD"})
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
refund, err := c.RefundSale(saleID, &paypalsdk.Amount{Total: "7.00", Currency: "USD"})
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
order, err := c.CreateOrder(paypalsdk.OrderIntentCapture, []paypalsdk.PurchaseUnitRequest{paypalsdk.PurchaseUnitRequest{ReferenceID: "ref-id", Amount: paypalsdk.Amount{Total: "7.00", Currency: "USD"}}})
```

### Authorize Order

```go
auth, err := c.AuthorizeOrder(orderID, &paypalsdk.Amount{Total: "7.00", Currency: "USD"})
```

### Capture Order

```go
capture, err := c.CaptureOrder(orderID, &paypalsdk.Amount{Total: "7.00", Currency: "USD"}, true, nil)
```

### Void Order

```go
order, err := c.VoidOrder(orderID)
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
payout := paypalsdk.Payout{
    SenderBatchHeader: &paypalsdk.SenderBatchHeader{
        EmailSubject: "Subject will be displayed on PayPal",
    },
    Items: []paypalsdk.PayoutItem{
        paypalsdk.PayoutItem{
            RecipientType: "EMAIL",
            Receiver:      "single-email-payout@mail.com",
            Amount: &paypalsdk.AmountPayout{
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
c.StoreCreditCard(paypalsdk.CreditCard{
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
c.PatchCreditCard("CARD-ID-123", []paypalsdk.CreditCardField{
    paypalsdk.CreditCardField{
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

### How to Contribute

* Fork a repository
* Add/Fix something
* Check that tests are passing
* Create PR

### Tests

* Unit tests: `go test`
* Integration tests: `go test -tags=integration`
