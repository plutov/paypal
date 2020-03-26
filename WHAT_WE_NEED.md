# Subscribe a.k.a. Billing Agreement

## Table of Contents:

- [Payments](#payments-api-docshttpsdeveloperpaypalcomdocsapipaymentsv1)
    - [Payment](#payment)
    - [Sale](#sale)
    - [Authorization](#authorization)
    - [Order](#orders)
    - [Capture](#capture)
    - [Refund](#refund)
- [Billing Plans](#billing-plans-api-docshttpsdeveloperpaypalcomdocsapipaymentsbilling-plansv1)
- [Billing Agreements](#billing-agreements-api-docshttpsdeveloperpaypalcomdocsapipaymentsbilling-agreementsv1)
- [Hooks](#hooks)

## Payments [API Docs](https://developer.paypal.com/docs/api/payments/v1/)

### Payment 
- [ ] Show [Payment](https://developer.paypal.com/docs/api/payments/v1/#payment_get) **GET** `/v1/payments/payment/{payment_id}`

<!--
- [ ] List [Payments](https://developer.paypal.com/docs/api/payments/v1/#payment_list) **GET** `/v1/payments/payment`
- [ ] Partially Update [Payment](https://developer.paypal.com/docs/api/payments/v1/#payment_update) **PATCH** `/v1/payments/payment/{payment_id}`
- [ ] Execute [Payment](https://developer.paypal.com/docs/api/payments/v1/#payment_execute) **POST** `/v1/payments/payment/{payment_id}/execute`
- [ ] Create [Payment](https://developer.paypal.com/docs/api/payments/v1/#payment_create) **POST** `/v1/payments/payment`
-->

### Sale

- [ ] Show [Sale](https://developer.paypal.com/docs/api/payments/v1/#sale_get) **GET** `/v1/payments/sale/{sale_id}`
- [ ] Refund [Sale](https://developer.paypal.com/docs/api/payments/v1/#sale_refund) **POST** `/v1/payments/sale/{sale_id}/refund`

<!--
### Authorization

- [ ] Show [Authorization](https://developer.paypal.com/docs/api/payments/v1/#authorization_get) **GET** `/v1/payments/authorization/{authorization_id}`
- [ ] Capture [Authorization](https://developer.paypal.com/docs/api/payments/v1/#authorization_capture) **POST** `/v1/payments/authorization/{authorization_id}/capture`
- [ ] Re-authorize [Payment](https://developer.paypal.com/docs/api/payments/v1/#authorization_reauthorize) **POST** `/v1/payments/authorization/{authorization_id}/reauthorize`
- [ ] Void [Authorization](https://developer.paypal.com/docs/api/payments/v1/#authorization_void) **POST** `/v1/payments/authorization/{authorization_id}/void`

### Orders

- [ ] Show [Order Details](https://developer.paypal.com/docs/api/payments/v1/#orders_get) **GET** `/v1/payments/orders/{order_id}`
- [ ] Authorize [Order](https://developer.paypal.com/docs/api/payments/v1/#orders_authorize) **POST** `/v1/payments/orders/{order_id}/authorize`
- [ ] Capture [Order](https://developer.paypal.com/docs/api/payments/v1/#orders_capture) **POST** `/v1/payments/orders/{order_id}/capture`
- [ ] Void [Order](https://developer.paypal.com/docs/api/payments/v1/#orders_void) **POST** `/v1/payments/orders/{order_id}/do-void`

### Capture

- [ ] Show [Capture Payment Details](https://developer.paypal.com/docs/api/payments/v1/#capture_get) **GET** `/v1/payments/capture/{capture_id}`
- [ ] Refund [Capture](https://developer.paypal.com/docs/api/payments/v1/#capture_refund) **POST** `/v1/payments/capture/{capture_id}/refund`
-->

### Refund

- [ ] Show [Refund](https://developer.paypal.com/docs/api/payments/v1/#refund_get) **GET** `/v1/payments/refund/{refund_id}`

## Billing Plans [API Docs](https://developer.paypal.com/docs/api/payments.billing-plans/v1/)

- [ ] Create [Billing Plan](https://developer.paypal.com/docs/api/payments.billing-plans/v1/) **POST** `/v1/payments/billing-plans`
- [ ] List [Billing Plans](https://developer.paypal.com/docs/api/payments.billing-plans/v1/#billing-plans_list) **GET** `/v1/payments/billing-plans`
- [ ] Update [Billing Plan](https://developer.paypal.com/docs/api/payments.billing-plans/v1/#billing-plans_patch) **PATCH** `/v1/payments/billing-plans/{id}`
- [ ] Show [Billing Plan](https://developer.paypal.com/docs/api/payments.billing-plans/v1/#billing-plans_get) **GET** `/v1/payments/billing-plans/{id}`

## Billing Agreements [API Docs](https://developer.paypal.com/docs/api/payments.billing-agreements/v1/)
<!--
- [ ] Create [Billing Agreement](https://developer.paypal.com/docs/api/payments.billing-agreements/v1/#billing-agreements_post) **POST** `/v1/payments/billing-agreements`
-->

- [ ] Update [Billing Agreement](https://developer.paypal.com/docs/api/payments.billing-agreements/v1/#billing-agreements_patch) **PATCH** `/v1/payments/billing-agreements/{agreement_id}`
- [ ] Show [Billing Agreement](https://developer.paypal.com/docs/api/payments.billing-agreements/v1/#billing-agreements_get) **GET** `/v1/payments/billing-agreements/{agreement_id}`
- [ ] Bill [Billing Agreement](https://developer.paypal.com/docs/api/payments.billing-agreements/v1/#billing-agreements_bill-balance) Balance **POST** `/v1/payments/billing-agreements/{agreement_id}/bill-balance`
- [ ] Cancel [Billing Agreement](https://developer.paypal.com/docs/api/payments.billing-agreements/v1/#billing-agreements_cancel) **POST** `/v1/payments/billing-agreements/{agreement_id}/cancel`
- [ ] Set [Billing Agreement](https://developer.paypal.com/docs/api/payments.billing-agreements/v1/#billing-agreements_set-balance) Balance **POST** `/v1/payments/billing-agreements/{agreement_id}/set-balance`
- [ ] Suspend [Billing Agreement](https://developer.paypal.com/docs/api/payments.billing-agreements/v1/#billing-agreements_suspend) **POST** `/v1/payments/billing-agreements/{agreement_id}/suspend`
- [ ] Reactivate [Billing Agreement](https://developer.paypal.com/docs/api/payments.billing-agreements/v1/#billing-agreements_re-activate) **POST** `/v1/payments/billing-agreements/{agreement_id}/re-activate`
- [ ] List [Billing Agreement](https://developer.paypal.com/docs/api/payments.billing-agreements/v1/#billing-agreements_transactions) Transactions **GET** `/v1/payments/billing-agreements/{agreement_id}/transactions`
- [ ] Execute [Billing Agreement](https://developer.paypal.com/docs/api/payments.billing-agreements/v1/#billing-agreements_agreement-execute) **POST** `/v1/payments/billing-agreements/{agreement_id}/agreement-execute`

## Hooks

- [ ] BILLING.PLAN.CREATED
- [ ] BILLING.PLAN.UPDATED
- [ ] BILLING.SUBSCRIPTION.CREATED
- [ ] BILLING.SUBSCRIPTION.CANCELLED
- [ ] BILLING.SUBSCRIPTION.SUSPENDED
- [ ] BILLING.SUBSCRIPTION.RE-ACTIVATED
- [ ] PAYMENT.SALE.COMPLETED
- [ ] PAYMENT.SALE.DENIED
- [ ] PAYMENT.SALE.PENDING
- [ ] PAYMENT.SALE.REFUNDED
- [ ] PAYMENT.SALE.REVERSED