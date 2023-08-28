package domains

import "context"

type TokenManager interface {
	GetTokens(ctx context.Context, guid string) (string, string, error)
	RefreshTokens(ctx context.Context, access, refresh string) (string, string, error)
}
