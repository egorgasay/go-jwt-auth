package services

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go-jwt-auth/internal/constants"
	"go-jwt-auth/internal/domains"
	"go-jwt-auth/internal/lib"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type GeneratorService struct {
	logger lib.Logger
}

func NewGeneratorService(logger lib.Logger) domains.GeneratorService {
	return &GeneratorService{
		logger: logger,
	}
}

func (g *GeneratorService) RefreshToken(ctx context.Context, refreshTTL time.Duration) (string, int64, error) {
	if ctx.Err() != nil {
		return "", 0, ctx.Err()
	}

	uuidObj, err := uuid.NewUUID()
	if err != nil {
		g.logger.Error("can't generate uuid for refresh token", zap.Error(err))
		return "", 0, constants.ErrGenerateUUID
	}
	refreshExp := time.Now().Add(refreshTTL).Unix()

	return uuidObj.String(), refreshExp, nil
}

func (g *GeneratorService) AccessToken(
	ctx context.Context,
	guid string, key []byte,
	accessTTL time.Duration,
) (access string, exp int64, err error) {

	if ctx.Err() != nil {
		return "", 0, ctx.Err()
	}

	if guid == "" {
		return "", 0, constants.ErrInvalidGUID
	}

	exp = time.Now().Add(accessTTL).Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			_guid: guid,
			_exp:  exp,
		})

	access, err = t.SignedString(key)
	if err != nil {
		g.logger.Error("can't sign token", zap.Error(err))
		return "", exp, constants.ErrSignToken
	}

	return access, exp, nil
}

// bcryptHashFrom generates a bcrypt hash from a given token.
func bcryptHashFrom(token []byte) (string, error) {
	bcryptHash, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bcryptHash), nil
}
