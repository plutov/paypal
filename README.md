[![Build Status](https://travis-ci.org/logpacker/paypalsdk.svg?branch=master)](https://travis-ci.org/logpacker/paypalsdk)

PayPal REST API

#### Usage

```go
// Create a client instance
c, err := paypalsdk.NewClient("clietnid", "secret", paypalsdk.APIBaseSandBox)
c.SetLogFile("/tpm/paypal-debug.log") // Set log file if necessary
```

```go
// When you will have authorization_code you can get an access_token
accessToken, err := c.GetAccessToken()
```

```go
// Now we can create a paypal payment
amount := Amount{
    Total:    15.1111,
    Currency: "USD",
}
paymentResult, err := c.CreateDirectPaypalPayment(amount, "http://example.com/redirect-uri", "http://example.com/cancel-uri", "Description for this payment")

// If paymentResult.ID is not empty and paymentResult.Links is also
// we can redirect user to approval page (paymentResult.Links[0]).
// After approval user will be redirected to return_url from Request with PaymentID
```

```go
// And the last step is to execute approved payment
// paymentID is returned via return_url
paymentID := "PAY-17S8410768582940NKEE66EQ"
// payerID is returned via return_url
payerID := "7E7MGXCWTTKK2"
executeResult, err := c.ExecuteApprovedPayment(paymentID, payerID)
```

```go
// Get created payment info
payment, err := c.GetPayment(paymentID)
```

```go
// Get all payments slice
payments, err := c.GetPayments()
```
