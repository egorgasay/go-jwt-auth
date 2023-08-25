package domains

import (
	"context"
	"time"
)

type GeneratorService interface {
	Access(ctx context.Context, guid string, key []byte, accessTTL time.Duration) (string, int64, error)
	Refresh(ctx context.Context, refreshTTL time.Duration) (string, int64, error)
}
