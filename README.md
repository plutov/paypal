[![Build Status](https://travis-ci.org/logpacker/paypalsdk.svg?branch=master)](https://travis-ci.org/logpacker/paypalsdk)

#### GO client for PayPal REST API

#### Coverage
 * POST /v1/oauth2/token
 * POST /v1/payments/payment
 * GET /v1/payments/payment/%ID%
 * GET /v1/payments/payment
 * GET /v1/payments/authorization/%ID%
 * POST /v1/payments/authorization/%ID%/capture
 * POST /v1/payments/authorization/%ID%/void
 * POST /v1/payments/authorization/%ID%/reauthorize

#### Create client

```go
import "github.com/logpacker/paypalsdk"
```

```go
// Create a client instance
c, err := paypalsdk.NewClient("clietnid", "secret", paypalsdk.APIBaseSandBox)
c.SetLogFile("/tpm/paypal-debug.log") // Set log file if necessary
```

#### Get access token
```go
// When you will have authorization_code you can get an access_token
accessToken, err := c.GetAccessToken()
```

#### Create direct paypal payment

```go
// Now we can create a paypal payment
amount := Amount{
    Total:    "15.11",
    Currency: "USD",
}
redirectURI := "http://example.com/redirect-uri"
cancelURI := "http://example.com/cancel-uri"
description := "Description for this payment"
paymentResult, err := c.CreateDirectPaypalPayment(amount, redirectURI, cancelURI, description)

// If paymentResult.ID is not empty and paymentResult.Links is also
// we can redirect user to approval page (paymentResult.Links[0]).
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
            Total:    "200",
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

#### Execute approved payment

```go
// And the last step is to execute approved payment
// paymentID is returned via return_url
paymentID := "PAY-17S8410768582940NKEE66EQ"
// payerID is returned via return_url
payerID := "7E7MGXCWTTKK2"
executeResult, err := c.ExecuteApprovedPayment(paymentID, payerID)
```

#### Get payment by ID

```go
// Get created payment info
payment, err := c.GetPayment(paymentID)
```

#### Get list of payments

```go
// Get all payments slice
payments, err := c.GetPayments()
```

#### Get authorization by ID

```go
auth, err := c.GetAuthorization("AUTH-1")
```

#### Capture authorization

```go
capture, err := c.CaptureAuthorization("AUTH-1", &paypalsdk.Amount{Total: "200", Currency: "USD"}, true)
```

#### Void authorization

```go
auth, err := c.VoidAuthorization("AUTH-1")
```

#### Reauthorize authorization

```go
auth, err := c.ReauthorizeAuthorization("AUTH-1", &paypalsdk.Amount{Total: "200", Currency: "USD"})
```
