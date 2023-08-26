package domains

import (
	"context"
	"time"
)

type GeneratorService interface {
	JWTToken(ctx context.Context, guid string, key []byte, ttl time.Duration) (token string, exp int64, err error)
}
