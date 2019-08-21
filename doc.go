/*
Package paypal provides a wrapper to PayPal API (https://developer.paypal.com/webapps/developer/docs/api/).
The first thing you do is to create a Client (you can select API base URL using paypal contants).
  c, err := paypal.NewClient("clientID", "secretID", paypal.APIBaseSandBox)
Then you can get an access token from PayPal:
  accessToken, err := c.GetAccessToken()
After you have an access token you can call built-in functions to get data from PayPal.
paypal will assign all responses to go structures.
*/
package paypal
