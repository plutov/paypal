[![Go Report Card](https://goreportcard.com/badge/logpacker/PayPal-Go-SDK)](https://goreportcard.com/report/logpacker/PayPal-Go-SDK)
[![Build Status](https://travis-ci.org/logpacker/PayPal-Go-SDK.svg?branch=master)](https://travis-ci.org/logpacker/PayPal-Go-SDK)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/logpacker/PayPal-Go-SDK)
[![Chat at https://gitter.im/logpacker/PayPal-Go-SDK](https://img.shields.io/badge/gitter-dev_chat-46bc99.svg)](https://gitter.im/logpacker/PayPal-Go-SDK?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Sourcegraph](https://sourcegraph.com/github.com/logpacker/PayPal-Go-SDK/-/badge.svg)](https://sourcegraph.com/github.com/logpacker/PayPal-Go-SDK?badge)

### GO client for PayPal REST API

### Coverage
 * POST /v1/oauth2/token
 * POST /v1/payments/payment
 * GET /v1/payments/payment/**ID**
 * GET /v1/payments/payment
 * GET /v1/payments/authorization/**ID**
 * POST /v1/payments/authorization/**ID**/capture
 * POST /v1/payments/authorization/**ID**/void
 * POST /v1/payments/authorization/**ID**/reauthorize
 * GET /v1/payments/sale/**ID**
 * POST /v1/payments/sale/**ID**/refund
 * GET /v1/payments/refund/**ID**
 * GET /v1/payments/orders/**ID**
 * POST /v1/payments/orders/**ID**/authorize
 * POST /v1/payments/orders/**ID**/capture
 * POST /v1/payments/orders/**ID**/do-void
 * POST /v1/identity/openidconnect/tokenservice
 * GET /v1/identity/openidconnect/userinfo/?schema=**SCHEMA**
 * POST /v1/payments/payouts?sync_mode=true
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

### Missing endpoints
It is possible that some endpoints are missing in this SDK Client, but you can use built-in **paypalsdk** functions to perform a request: **NewClient -> NewRequest -> SendWithAuth**

### Create Client

```go
import "github.com/logpacker/PayPal-Go-SDK"
```

```go
// Create a client instance
c, err := paypalsdk.NewClient("clientID", "secretID", paypalsdk.APIBaseSandBox)
c.SetLog(os.Stdout) // Set log to terminal stdout
```

### Get access token
```go
// When you will have authorization_code you can get an access_token
accessToken, err := c.GetAccessToken()
```

### Create direct paypal payment

```go
// Now we can create a paypal payment
amount := paypalsdk.Amount{
    Total:    "7.00",
    Currency: "USD",
}
redirectURI := "http://example.com/redirect-uri"
cancelURI := "http://example.com/cancel-uri"
description := "Description for this payment"
paymentResult, err := c.CreateDirectPaypalPayment(amount, redirectURI, cancelURI, description)

// If paymentResult.ID is not empty and paymentResult.Links is also
// we can redirect user to approval page (paymentResult.Links[1]).
// After approval user will be redirected to return_url from Request with PaymentID
```

### Create any payment
```go
p := paypalsdk.Payment{
    Intent: "sale",
    Payer: &paypalsdk.Payer{
        PaymentMethod: "credit_card",
        FundingInstruments: []paypalsdk.FundingInstrument{paypalsdk.FundingInstrument{
            CreditCard: &paypalsdk.CreditCard{
                Number:      "4111111111111111",
                Type:        "visa",
                ExpireMonth: "11",
                ExpireYear:  "2020",
                CVV2:        "777",
                FirstName:   "John",
                LastName:    "Doe",
            },
        }},
    },
    Transactions: []paypalsdk.Transaction{paypalsdk.Transaction{
        Amount: &paypalsdk.Amount{
            Currency: "USD",
            Total:    "7.00",
        },
        Description: "My Payment",
    }},
    RedirectURLs: &paypalsdk.RedirectURLs{
        ReturnURL: "http://...",
        CancelURL: "http://...",
    },
}
paymentResponse, err := client.CreatePayment(p)
```

### Execute approved payment

```go
// And the last step is to execute approved payment
// paymentID is returned via return_url
paymentID := "PAY-17S8410768582940NKEE66EQ"
// payerID is returned via return_url
payerID := "7E7MGXCWTTKK2"
executeResult, err := c.ExecuteApprovedPayment(paymentID, payerID)
```

### Get payment by ID

```go
// Get created payment info
payment, err := c.GetPayment(paymentID)
```

### Get list of payments

```go
// Get all payments slice
payments, err := c.GetPayments()
```

### Get authorization by ID

```go
authID := "2DC87612EK520411B"
auth, err := c.GetAuthorization(authID)
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
saleID := "36C38912MN9658832"
sale, err := c.GetSale(saleID)
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
orderID := "O-4J082351X3132253H"
refund, err := c.GetRefund(orderID)
```

### Get Order by ID

```go
order, err := c.GetOrder(orderID)
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
// Retreive tolen by authorization code
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
// https://developer.paypal.com/docs/api/vault/

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

// get all stored credit cards
c.GetCreditCards(nil)
```

### How to Contribute

* Fork a repository
* Add/Fix something
* Run ./before-commit.sh to check that tests passed and code is formatted well
* Push to your repository
* Create pull request
