package domains

import (
	"context"
	"go-jwt-auth/internal/models"
)

type Database interface {
	SaveToken(ctx context.Context, guid string, t models.Token) error
	GetRefTokenByGUID(ctx context.Context, refresh string) (t models.Token, err error)
}
