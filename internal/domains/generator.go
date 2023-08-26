package domains

import (
	"context"
	"time"
)

type GeneratorService interface {
	AccessToken(ctx context.Context, guid string, key []byte, accessTTL time.Duration) (string, int64, error)
	RefreshToken(ctx context.Context, refreshTTL time.Duration) (string, int64, error)
}
