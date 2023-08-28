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
) (access string, err error) {

	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	if guid == "" {
		return "", constants.ErrInvalidGUID
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			_guid: guid,
		})

	access, err = t.SignedString(key)
	if err != nil {
		g.logger.Error("can't sign token", zap.Error(err))
		return "", constants.ErrSignToken
	}

	return access, nil
}

func (g *GeneratorService) RefreshToken(
	ctx context.Context,
) (token string, err error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	uuidObj, err := uuid.NewUUID()
	if err != nil {
		g.logger.Error("can't generate uuid", zap.Error(err))
		return "", constants.ErrGenerateToken
	}

	return uuidObj.String(), nil
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
