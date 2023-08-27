package domains

import (
	"context"
	"go-jwt-auth/internal/models"
)

type Database interface {
	SaveToken(ctx context.Context, t models.TokenData) error
	GetTokensDataByGUID(ctx context.Context, guid string) (t []models.TokenData, err error)
	DeleteTokenData(ctx context.Context, guid, hash string) error
}
