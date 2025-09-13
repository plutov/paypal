package paypal

import "context"

type CacheToken struct {
	RefreshToken string `json:"refresh_token"`
	Token        string `json:"access_token"`
	Type         string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	ExpiresAt    int64  `json:"expires_at"`
}

type Cache interface {
	GetToken(ctx context.Context, clientId string) (*CacheToken, error)
	SetToken(ctx context.Context, clientId string, token *CacheToken) error
}
