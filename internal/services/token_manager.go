package services

import (
	"context"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go-jwt-auth/internal/domains"
	"go-jwt-auth/internal/lib"
	"go-jwt-auth/internal/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// TokenManager is a service for managing tokens.
type TokenManager struct {
	storage    domains.Repository
	logger     lib.Logger
	key        []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewTokenManager creates a new instance of TokenManager.
//
// It takes a repository, a logger, and a config as parameters.
// It returns a TokenManager and an error.
func NewTokenManager(st domains.Repository, logger lib.Logger, conf lib.Config) (domains.TokenManager, error) {
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
		storage:    st,
		logger:     logger,
		key:        []byte(conf.JWT.Key),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
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
	access, _, err = tm.generateAccess(guid)
	if err != nil {
		tm.logger.Error("can't generate access token", zap.Error(err))
		return "", "", err
	}

	refresh, refreshExp, err := tm.generateRefresh()
	if err != nil {
		tm.logger.Error("can't generate refresh token", zap.Error(err))
		return "", "", err
	}

	bcryptHash, err := bcryptHashFrom(refresh)
	if err != nil {
		tm.logger.Error("can't hash refresh token", zap.Error(err))
		return "", "", ErrCantHashToken
	}

	if err := tm.storage.SaveRefresh(ctx, guid, models.RefreshToken{
		RefreshTokenBCrypt: bcryptHash,
		Exp:                refreshExp,
	}); err != nil {
		tm.logger.Error("can't save refresh token", zap.Error(err))
		return "", "", ErrInvalidToken
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
		return "", "", ErrInvalidToken
	}

	bcryptHash, err := bcryptHashFrom(string(oldRefreshBytes))
	if err != nil {
		tm.logger.Error("can't hash refresh token", zap.Error(err))
		return "", "", ErrCantHashToken
	}

	guid, refToken, err := tm.storage.GetRefTokenAndGUID(ctx, bcryptHash)
	if err != nil {
		tm.logger.Error("can't get refresh token exp and id", zap.Error(err))
		return "", "", ErrInvalidToken
	}

	if refToken.Exp < time.Now().Unix() {
		tm.logger.Error("refresh token expired")
		return "", "", ErrExpired
	}

	return tm.GetTokens(ctx, guid)
}

// bcryptHashFrom generates a bcrypt hash from a given token.
//
// It takes a string parameter named 'token' and returns a string representing the bcrypt hash and an error.
func bcryptHashFrom(token string) (string, error) {
	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bcryptHash), nil
}

// generateAccess generates an access token for a given GUID.
func (tm *TokenManager) generateAccess(guid string) (string, int64, error) {
	if guid == "" {
		return "", 0, ErrInvalidGUID
	}

	exp := time.Now().Add(tm.accessTTL).Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			_guid: guid,
			_exp:  exp,
		})

	access, err := t.SignedString(tm.key)
	if err != nil {
		tm.logger.Error("can't sign token", zap.Error(err))
		return "", exp, ErrSign
	}

	return access, exp, nil
}

// generateRefresh generates a refresh token.
func (tm *TokenManager) generateRefresh() (string, int64, error) {
	uuidObj, err := uuid.NewUUID()
	if err != nil {
		tm.logger.Error("can't generate uuid for refresh token", zap.Error(err))
		return "", 0, ErrGenerateUUID
	}
	refreshExp := time.Now().Add(tm.refreshTTL).Unix()

	return uuidObj.String(), refreshExp, nil
}
