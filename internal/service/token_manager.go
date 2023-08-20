package service

import (
	"context"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go-jwt-auth/internal/model"
	"go.uber.org/zap"
	"time"
)

type TokenManager struct {
	storage    storage
	logger     *zap.Logger
	key        []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type JWTConfig struct {
	Key        string `json:"key"`
	AccessTTL  string `json:"access_ttl"`
	RefreshTTL string `json:"refresh_ttl"`
}

type storage interface {
	SaveRefresh(ctx context.Context, guid string, refresh model.RefreshToken) error
	GetRefTokenAndGUID(ctx context.Context, refresh string) (guid string, rt model.RefreshToken, err error)
}

func NewTokenManager(st storage, logger *zap.Logger, jwtConf JWTConfig) (*TokenManager, error) {
	accessTTL, err := time.ParseDuration(jwtConf.AccessTTL)
	if err != nil {
		logger.Error("can't parse access_ttl", zap.Error(err))
		return nil, err
	}

	refreshTTL, err := time.ParseDuration(jwtConf.RefreshTTL)
	if err != nil {
		logger.Error("can't parse refresh_ttl", zap.Error(err))
		return nil, err
	}

	return &TokenManager{
		storage:    st,
		logger:     logger,
		key:        []byte(jwtConf.Key),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}, nil
}

const (
	_guid = "guid"
	_exp  = "exp"
)

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

	if err := tm.storage.SaveRefresh(ctx, guid, model.RefreshToken{
		RefreshTokenBCrypt: refresh,
		Exp:                refreshExp,
	}); err != nil {
		tm.logger.Error("can't save refresh token", zap.Error(err))
		return "", "", ErrInvalidToken
	}

	refreshB64 := base64.StdEncoding.EncodeToString([]byte(refresh))

	return access, refreshB64, nil
}

func (tm *TokenManager) RefreshTokens(ctx context.Context, oldRefreshB64 string) (access string, refresh string, err error) {
	oldRefreshBytes, err := base64.StdEncoding.DecodeString(oldRefreshB64)
	if err != nil {
		tm.logger.Error("can't decode refresh token", zap.Error(err))
		return "", "", ErrInvalidToken
	}

	guid, refToken, err := tm.storage.GetRefTokenAndGUID(ctx, string(oldRefreshBytes))
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

func (tm *TokenManager) generateRefresh() (string, int64, error) {
	uuidObj, err := uuid.NewUUID()
	if err != nil {
		tm.logger.Error("can't generate uuid for refresh token", zap.Error(err))
		return "", 0, ErrGenerateUUID
	}
	refreshExp := time.Now().Add(tm.refreshTTL).Unix()

	return uuidObj.String(), refreshExp, nil
}
