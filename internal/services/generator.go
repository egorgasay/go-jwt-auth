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
			_iat:  exp,
		})

	access, err = t.SignedString(key)
	if err != nil {
		g.logger.Error("can't sign token", zap.Error(err))
		return "", 0, constants.ErrSignToken
	}

	return access, exp, nil
}

func (g *GeneratorService) RefreshToken(
	ctx context.Context,
	refreshTTL time.Duration,
) (token string, exp int64, err error) {
	if ctx.Err() != nil {
		return "", 0, ctx.Err()
	}

	uuidObj, err := uuid.NewUUID()
	if err != nil {
		g.logger.Error("can't generate uuid", zap.Error(err))
		return "", 0, constants.ErrGenerateToken
	}

	return uuidObj.String(), time.Now().Add(refreshTTL).Unix(), nil
}

// bcryptHashFrom generates a bcrypt hash from a given token.
func bcryptHashFrom(token []byte) ([]byte, error) {
	if len(token) == 0 {
		return nil, constants.ErrInvalidToken
	}

	bcryptHash, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return bcryptHash, nil
}
