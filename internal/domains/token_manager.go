package domains

import "context"

type TokenManager interface {
	GetTokens(ctx context.Context, guid string) (string, string, error)
	RefreshTokens(ctx context.Context, refresh string) (string, string, error)
}
