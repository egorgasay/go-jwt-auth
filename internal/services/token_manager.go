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
func (tm *TokenManager) GetTokens(ctx context.Context, guid string) (access string, refresh string, err error) {
	access, err = tm.generator.AccessToken(ctx, guid, tm.key)
	if err != nil {
		tm.logger.Error("can't generate access token", zap.Error(err))
		return "", "", errors.Join(constants.ErrGenerate, err)
	}

	refresh, err = tm.generator.RefreshToken(ctx)
	if err != nil {
		tm.logger.Error("can't generate refresh token", zap.Error(err))
		return "", "", errors.Join(constants.ErrGenerate, err)
	}

	bcryptHash, err := bcryptHashFrom([]byte(refresh))
	if err != nil {
		tm.logger.Error("can't hash refresh token", zap.Error(err))
		return "", "", constants.ErrCantHashToken
	}

	if err := tm.repository.SaveToken(ctx, models.TokenData{
		GUID:        guid,
		RefreshHash: string(bcryptHash),
		RefreshExp:  time.Now().Add(tm.refreshTTL).Unix(),
		AccessExp:   time.Now().Add(tm.accessTTL).Unix(),
	}); err != nil {
		tm.logger.Error("can't save token", zap.Error(err))
		if errors.Is(err, constants.ErrAlreadyExists) {
			return "", "", constants.ErrAlreadyExists
		}

		return "", "", constants.ErrRepository
	}

	refreshB64 := base64.StdEncoding.EncodeToString([]byte(refresh))
	accessB64 := base64.StdEncoding.EncodeToString([]byte(access))

	return accessB64, refreshB64, nil
}

// RefreshTokens retrieves the access and refresh tokens for a given GUID.
func (tm *TokenManager) RefreshTokens(ctx context.Context, oldAccessB64, oldRefreshB64 string) (access string, refresh string, err error) {
	if oldRefreshB64 == "" {
		return "", "", constants.ErrMissingRefreshToken
	} else if oldAccessB64 == "" {
		return "", "", constants.ErrMissingAccessToken
	}

	oldRefreshBytes, err := base64.StdEncoding.DecodeString(oldRefreshB64)
	if err != nil {
		tm.logger.Error("can't decode refresh token", zap.Error(err))
		return "", "", constants.ErrInvalidToken
	}

	oldAccessBytes, err := base64.StdEncoding.DecodeString(oldAccessB64)
	if err != nil {
		tm.logger.Error("can't decode access token", zap.Error(err))
		return "", "", constants.ErrInvalidToken
	}

	guid, err := tm.guidFromJWT(string(oldAccessBytes))
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
		if err = validateTokenHash([]byte(tokenData.RefreshHash), oldRefreshBytes); err != nil {
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

// validateTokenHash validates the hash of a given token.
func validateTokenHash(hash []byte, incoming []byte) error {
	err := bcrypt.CompareHashAndPassword(hash, incoming)
	if err != nil {
		return err
	}

	return nil
}

// guidFromJWT extracts the GUID from a given JWT token.
func (tm *TokenManager) guidFromJWT(token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return tm.key, nil
	})
	if err != nil {
		tm.logger.Error("can't parse token", zap.Error(err))
		return "", constants.ErrInvalidToken
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		tm.logger.Error("can't extract claims from token", zap.Error(err))
		return "", constants.ErrInvalidToken
	}

	guid, ok := claims[_guid]
	if !ok {
		tm.logger.Error("can't extract guid from token", zap.Error(err))
		return "", constants.ErrInvalidToken
	}

	guidStr, ok := guid.(string)
	if !ok {
		tm.logger.Error("can't convert guid to string", zap.Error(err))
		return "", constants.ErrInvalidToken
	}

	return guidStr, nil
}
