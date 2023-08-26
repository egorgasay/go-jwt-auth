package services

import (
	"bytes"
	"context"
	"github.com/golang-jwt/jwt/v5"
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

func (g *GeneratorService) JWTToken(
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
	if len(token) == 0 {
		return "", constants.ErrInvalidToken
	}
	buf := bytes.Buffer{}
	last := len(token) > constants.MaxBcryptLength

	for len(token) > constants.MaxBcryptLength {
		bcryptHash, err := bcrypt.GenerateFromPassword(token[:constants.MaxBcryptLength], bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}

		buf.Write(bcryptHash)

		if len(token) > constants.MaxBcryptLength {
			token = token[constants.MaxBcryptLength:]
		} else if last {
			break
		} else {
			last = true
		}
		buf.WriteString(" ")
	}

	return buf.String(), nil
}
