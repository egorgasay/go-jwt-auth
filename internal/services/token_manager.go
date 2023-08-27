package services

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go-jwt-auth/internal/constants"
	"go-jwt-auth/internal/domains"
	"go-jwt-auth/internal/lib"
	"go-jwt-auth/internal/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// TokenManager is a service for managing tokens.
type TokenManager struct {
	repository domains.Repository
	logger     lib.Logger
	key        []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	generator  domains.GeneratorService
}

// NewTokenManager creates a new instance of TokenManager.
func NewTokenManager(
	st domains.Repository,
	logger lib.Logger,
	conf lib.Config,
	generator domains.GeneratorService,
) (domains.TokenManager, error) {

	accessTTL, err := time.ParseDuration(conf.JWT.AccessTTL)
	if err != nil {
		logger.Error("can't parse access_ttl", zap.Error(err))
		return nil, err
	}

	refreshTTL, err := time.ParseDuration(conf.JWT.RefreshTTL)
	if err != nil {
		logger.Error("can't parse refresh_ttl", zap.Error(err))
		return nil, err
	}

	return &TokenManager{
		repository: st,
		logger:     logger,
		key:        []byte(conf.JWT.Key),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		generator:  generator,
	}, nil
}

const (
	_guid = "guid"
	_exp  = "exp"
)

// GetTokens retrieves the access and refresh tokens for a given GUID.
//
// ctx - the context.Context object for the request
// guid - the unique identifier
// access - the access token
// refresh - the refresh token
// err - any error that occurred during token generation
func (tm *TokenManager) GetTokens(ctx context.Context, guid string) (access string, refresh string, err error) {
	access, accessExp, err := tm.generator.JWTToken(ctx, guid, tm.key, tm.accessTTL)
	if err != nil {
		tm.logger.Error("can't generate access token", zap.Error(err))
		return "", "", errors.Join(err, constants.ErrGenerate)
	}

	refresh, refreshExp, err := tm.generator.JWTToken(ctx, guid, tm.key, tm.refreshTTL)
	if err != nil {
		tm.logger.Error("can't generate access token", zap.Error(err))
		return "", "", errors.Join(err, constants.ErrGenerate)
	}

	bcryptHash, err := bcryptHashFrom([]byte(refresh))
	if err != nil {
		tm.logger.Error("can't hash refresh token", zap.Error(err))
		return "", "", constants.ErrCantHashToken
	}

	if err := tm.repository.SaveToken(ctx, models.TokenData{
		GUID:        guid,
		RefreshHash: bcryptHash,
		RefreshExp:  refreshExp,
		AccessExp:   accessExp,
	}); err != nil {
		tm.logger.Error("can't save token", zap.Error(err))
		if errors.Is(err, constants.ErrAlreadyExists) {
			return "", "", constants.ErrAlreadyExists
		}

		return "", "", constants.ErrRepository
	}

	refreshB64 := base64.StdEncoding.EncodeToString([]byte(refresh))

	return access, refreshB64, nil
}

// RefreshTokens refreshes the access and refresh tokens for a given old refresh token.
//
// ctx: the context.Context object for managing the lifecycle of the request.
// oldRefreshB64: the base64 encoded string of the old refresh token.
// access: the new access token.
// refresh: the new refresh token.
// err: any error that occurred during the token refresh.
// Returns the new access token, the new refresh token, and an error if any.
func (tm *TokenManager) RefreshTokens(ctx context.Context, oldRefreshB64 string) (access string, refresh string, err error) {
	if oldRefreshB64 == "" {
		return "", "", constants.ErrNoToken
	}

	oldRefreshBytes, err := base64.StdEncoding.DecodeString(oldRefreshB64)
	if err != nil {
		tm.logger.Error("can't decode refresh token", zap.Error(err))
		return "", "", constants.ErrInvalidToken
	}

	guid, err := tm.guidFromRefresh(string(oldRefreshBytes))
	if err != nil {
		return "", "", err
	}

	userTokens, err := tm.repository.GetTokensDataByGUID(ctx, guid)
	if err != nil {
		tm.logger.Debug("can't get token by guid", zap.Error(err))
		if errors.Is(err, constants.ErrNotFound) {
			return "", "", constants.ErrNotFound
		}
		return "", "", err
	}

	for _, tokenData := range userTokens {
		if err = validateTokenHash(tokenData.RefreshHash, oldRefreshBytes); err != nil {
			continue
		}

		if tokenData.RefreshExp < time.Now().Unix() {
			err = constants.ErrTokenExpired
			continue
		}

		if err = tm.repository.DeleteTokenData(ctx, guid, tokenData.RefreshHash); err != nil {
			tm.logger.Error("can't delete token", zap.Error(err))
			return "", "", constants.ErrRepository
		}
		break
	}

	if err != nil {
		tm.logger.Debug("can't validate token", zap.Error(err))
		return "", "", constants.ErrInvalidToken
	}

	return tm.GetTokens(ctx, guid)
}

func validateTokenHash(hash string, incoming []byte) error {
	hashParts := strings.Split(hash, " ")
	incomingPartsLen := len(incoming) / constants.MaxBcryptLength
	if len(hashParts) != incomingPartsLen && incomingPartsLen+1 != len(hashParts) {
		return constants.ErrInvalidToken
	}

	last := len(incoming) <= constants.MaxBcryptLength
	for i := 0; len(incoming) > constants.MaxBcryptLength; i++ {
		err := bcrypt.CompareHashAndPassword([]byte(hashParts[i]), incoming[:constants.MaxBcryptLength])
		if err != nil {
			return err
		}

		if len(incoming) > constants.MaxBcryptLength {
			incoming = incoming[constants.MaxBcryptLength:]
		} else if last {
			break
		} else {
			last = true
		}
	}

	return nil
}

func (tm *TokenManager) guidFromRefresh(refresh string) (string, error) {
	t, err := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
		return tm.key, nil
	})
	if err != nil {
		tm.logger.Debug("can't parse refresh token", zap.Error(err))
		return "", constants.ErrInvalidToken
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		tm.logger.Error("can't extract claims from refresh token", zap.Error(err))
		return "", constants.ErrInvalidToken
	}

	guid, ok := claims[_guid]
	if !ok {
		tm.logger.Error("can't extract guid from refresh token", zap.Error(err))
		return "", constants.ErrInvalidToken
	}

	guidStr, ok := guid.(string)
	if !ok {
		tm.logger.Error("can't convert guid to string", zap.Error(err))
		return "", constants.ErrInvalidToken
	}

	return guidStr, nil
}
