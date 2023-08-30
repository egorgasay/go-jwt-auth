package domains

import (
	"context"
	"time"
)

type GeneratorService interface {
	AccessToken(ctx context.Context, guid string, key []byte, accessTTL time.Duration) (token string, iat int64, err error)
	RefreshToken(ctx context.Context, refreshTTL time.Duration) (token string, iat int64, err error)
}
