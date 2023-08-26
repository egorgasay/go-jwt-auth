package services

import (
	"context"
	"encoding/base64"
	"errors"
	"go-jwt-auth/internal/constants"
	"go-jwt-auth/internal/domains"
	"go-jwt-auth/internal/lib"
	"go-jwt-auth/internal/models"
	"go.uber.org/zap"
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
	access, _, err = tm.generator.AccessToken(ctx, guid, tm.key, tm.accessTTL)
	if err != nil {
		tm.logger.Error("can't generate access token", zap.Error(err))
		return "", "", errors.Join(err, constants.ErrGenerate)
	}

	refresh, refreshExp, err := tm.generator.RefreshToken(ctx, tm.refreshTTL)
	if err != nil {
		tm.logger.Error("can't generate refresh token", zap.Error(err))
		return "", "", errors.Join(err, constants.ErrGenerate)
	}

	bcryptHash, err := bcryptHashFrom([]byte(refresh))
	if err != nil {
		tm.logger.Error("can't hash refresh token", zap.Error(err))
		return "", "", constants.ErrCantHashToken
	}

	if err := tm.repository.SaveRefresh(ctx, guid, models.RefreshToken{
		RefreshTokenBCrypt: bcryptHash,
		Exp:                refreshExp,
	}); err != nil {
		tm.logger.Error("can't save refresh token", zap.Error(err))
		return "", "", errors.Join(err, constants.ErrRepository)
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
	oldRefreshBytes, err := base64.StdEncoding.DecodeString(oldRefreshB64)
	if err != nil {
		tm.logger.Error("can't decode refresh token", zap.Error(err))
		return "", "", constants.ErrInvalidToken
	}

	bcryptHash, err := bcryptHashFrom(oldRefreshBytes)
	if err != nil {
		tm.logger.Error("can't hash refresh token", zap.Error(err))
		return "", "", constants.ErrCantHashToken
	}

	guid, refToken, err := tm.repository.GetRefTokenAndGUID(ctx, bcryptHash)
	if err != nil {
		tm.logger.Error("can't get refresh token exp and id", zap.Error(err))
		return "", "", errors.Join(err, constants.ErrRepository)
	}

	if refToken.Exp < time.Now().Unix() {
		tm.logger.Error("refresh token expired")
		return "", "", constants.ErrExpired
	}

	return tm.GetTokens(ctx, guid)
}
