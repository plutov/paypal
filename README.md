[![Build Status](https://travis-ci.org/logpacker/paypalsdk.svg?branch=master)](https://travis-ci.org/logpacker/paypalsdk)

PayPal REST API

#### Usage

```
// Create a client instance
c, err := paypalsdk.NewClient("clietnid", "secret", paypalsdk.APIBaseSandBox)
```

```
// Redirect client to this URL with provided redirect URI and necessary scopes. It's necessary to retreive authorization_code
authCodeURL, err := c.GetAuthorizationCodeURL("https://example.com/redirect-uri", []string{"address"})
```

```
// When you will have authorization_code you can get an access_token
accessToken, err := c.GetAccessToken(authCode, "https://example.com/redirect-uri")
```
