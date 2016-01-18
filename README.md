[![Build Status](https://travis-ci.org/logpacker/PayPal-Go-SDK.svg?branch=master)](https://travis-ci.org/logpacker/PayPal-Go-SDK)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/logpacker/PayPal-Go-SDK)

#### GO client for PayPal REST API

#### Coverage
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

#### Missing endpoints
It is possible that some endpoints are missing in this SDK Client, but you can use built-in **paypalsdk** functions to perform a request: **NewClient -> NewRequest -> SendWithAuth**

#### Create Client

```go
import "github.com/logpacker/PayPal-Go-SDK"
```

```go
// Create a client instance
c, err := paypalsdk.NewClient("clientID", "secretID", paypalsdk.APIBaseSandBox)
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
amount := paypalsdk.Amount{
    Total:    "7.00",
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
authID := "2DC87612EK520411B"
auth, err := c.GetAuthorization(authID)
```

#### Capture authorization

```go
capture, err := c.CaptureAuthorization(authID, &paypalsdk.Amount{Total: "7.00", Currency: "USD"}, true)
```

#### Void authorization

```go
auth, err := c.VoidAuthorization(authID)
```

#### Reauthorize authorization

```go
auth, err := c.ReauthorizeAuthorization(authID, &paypalsdk.Amount{Total: "7.00", Currency: "USD"})
```

#### Get Sale by ID

```go
saleID := "36C38912MN9658832"
sale, err := c.GetSale(saleID)
```

#### Refund Sale by ID

```go
// Full
refund, err := c.RefundSale(saleID, nil)
// Partial
refund, err := c.RefundSale(saleID, &paypalsdk.Amount{Total: "7.00", Currency: "USD"})
```

#### Get Refund by ID

```go
orderID := "O-4J082351X3132253H"
refund, err := c.GetRefund(orderID)
```

#### Get Order by ID

```go
order, err := c.GetOrder(orderID)
```

#### Authorize Order

```go
auth, err := c.AuthorizeOrder(orderID, &paypalsdk.Amount{Total: "7.00", Currency: "USD"})
```

#### Capture Order

```go
capture, err := c.CaptureOrder(orderID, &paypalsdk.Amount{Total: "7.00", Currency: "USD"}, true, nil)
```

#### Void Order

```go
order, err := c.VoidOrder(orderID)
```

#### Identity

```go
// Retreive tolen by authorization code
token, err := c.GrantNewAccessTokenFromAuthCode("<Authorization-Code>", "http://example.com/myapp/return.php")
// ... or by refresh token
token, err := c.GrantNewAccessTokenFromRefreshToken("<Refresh-Token>")
```

#### How to Contribute

* Fork a repository
* Add/Fix something
* Run ./before-commit.sh to check that tests passed and code is formatted well
* Push to your repository
* Create pull request
