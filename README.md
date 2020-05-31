[![Go Report Card](https://goreportcard.com/badge/plutov/paypal)](https://goreportcard.com/report/plutov/paypal)
[![Build Status](https://travis-ci.org/plutov/paypal.svg?branch=master)](https://travis-ci.org/plutov/paypal)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/plutov/paypal)

### Go client for PayPal REST API

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
 * GET /v2/payments/authorizations/**ID**
 * POST /v2/payments/authorizations/**ID**/capture
 * POST /v2/payments/authorizations/**ID**/void
 * POST /v2/payments/authorizations/**ID**/reauthorize
 * GET /v1/payments/sale/**ID**
 * POST /v1/payments/sale/**ID**/refund
 * GET /v2/payments/refund/**ID**
 * POST /v1/reporting/transactions
 #Vault
  * POST /v1/vault/credit-cards
  * DELETE /v1/vault/credit-cards/**ID**
  * PATCH /v1/vault/credit-cards/**ID**
  * GET /v1/vault/credit-cards/**ID**
  * GET /v1/vault/credit-cards
 #Checkout
 * POST /v2/checkout/orders
 * GET /v2/checkout/orders/**ID**
 * PATCH /v2/checkout/orders/**ID**
 * POST /v2/checkout/orders/**ID**/authorize
 * POST /v2/checkout/orders/**ID**/capture
 #Billing plans (payments)
 * GET /v1/payments/billing-plans
 * POST /v1/payments/billing-plans
 * PATCH /v1/payments/billing-plans/***ID***
 * POST /v1/payments/billing-agreements
 * POST /v1/payments/billing-agreements/***TOKEN***/agreement-execute
 #Notifications 
 * POST /v1/notifications/webhooks
 * GET /v1/notifications/webhooks
 * GET /v1/notifications/webhooks/**ID**
 * PATCH /v1/notifications/webhooks/**ID**
 * DELETE /v1/notifications/webhooks/**ID**
 * POST /v1/notifications/verify-webhook-signature
 #Products (Catalog)
 * POST /v1/catalogs/products
 * PATCH /v1/catalogs/products/**ID**
 * GET /v1/catalogs/products/**ID**
 * GET /v1/catalogs/products
 #Billing Plans (Subscriptions)
 * POST  /v1/billing/plans
 * PATCH /v1/billing/plans/**ID**
 * GET   /v1/billing/plans/**ID**
 * GET   /v1/billing/plans
 * POST  /v1/billing/plans/**ID**/activate
 * POST  /v1/billing/plans/**ID**/deactivate
 * POST  /v1/billing/plans/**ID**/update-pricing-schemes
 #Subscriptions
 * POST /v1/billing/subscriptions
 * PATCH /v1/billing/subscriptions/**ID**
 * GET /v1/billing/subscriptions/**ID**
 * POST /v1/billing/subscriptions/**ID**/activate
 * POST /v1/billing/subscriptions/**ID**/cancel
 * POST /v1/billing/subscriptions/**ID**/capture
 * POST /v1/billing/subscriptions/**ID**/suspend
 * GET /v1/billing/subscriptions/**ID**/transactions
 
### Missing endpoints
It is possible that some endpoints are missing in this SDK Client, but you can use built-in **paypal** functions to perform a request: **NewClient -> NewRequest -> SendWithAuth**

### Usage

```go
// If using Go Modules
// import "github.com/plutov/paypal/v3" 
import "github.com/plutov/paypal"

// Create a client instance
c, err := paypal.NewClient("clientID", "secretID", paypal.APIBaseSandBox)
c.SetLog(os.Stdout) // Set log to terminal stdout

accessToken, err := c.GetAccessToken()
```

### How to Contribute

* Fork a repository
* Add/Fix something
* Check that tests are passing
* Create PR

Current contributors:

- [Roopak Venkatakrishnan](https://github.com/roopakv)
- [Alex Pliutau](https://github.com/plutov)

### Tests

* Unit tests: `go test -v ./...`
* Integration tests: `go test -tags=integration`
