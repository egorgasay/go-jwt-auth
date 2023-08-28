package domains

import (
	"context"
)

type GeneratorService interface {
	AccessToken(ctx context.Context, guid string, key []byte) (token string, err error)
	RefreshToken(ctx context.Context) (token string, err error)
}
