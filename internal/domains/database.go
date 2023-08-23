package domains

import (
	"context"
	"go-jwt-auth/internal/models"
)

type Database interface {
	SaveRefresh(ctx context.Context, guid string, refresh models.RefreshToken) error
	GetRefTokenAndGUID(ctx context.Context, refresh string) (guid string, rt models.RefreshToken, err error)
}
