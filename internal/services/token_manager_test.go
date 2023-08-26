package services

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"go-jwt-auth/internal/constants"
	"go-jwt-auth/internal/domains/mocks"
	"go-jwt-auth/internal/lib"
	"go-jwt-auth/internal/models"
	"testing"
	"time"
)

type (
	repoMock func(c *mocks.Repository)
	genMock  func(c *mocks.GeneratorService)
)

var (
	_contextType = mock.AnythingOfType("context.backgroundCtx")
	_rtokenType  = mock.AnythingOfType("models.RefreshToken")
	_stringType  = mock.AnythingOfType("string")
)

func TestTokenManager_GetTokens(t *testing.T) {

	const (
		accessTTL  = time.Minute
		refreshTTL = time.Hour
	)

	type args struct {
		ctx  context.Context
		guid string
	}
	tests := []struct {
		name        string
		args        args
		wantAccess  string
		wantRefresh string
		genMock     genMock
		repoMock    repoMock
		wantErr     error
	}{
		{
			name: "ok",
			args: args{
				ctx:  context.Background(),
				guid: "123",
			},
			wantAccess:  "123",
			wantRefresh: "MTIz",
			genMock: func(c *mocks.GeneratorService) {
				c.On("AccessToken", _contextType, "123", []byte("123"), accessTTL).
					Return("123", int64(123), nil)
				c.On("RefreshToken", _contextType, refreshTTL).
					Return("123", int64(123), nil)
			},
			repoMock: func(c *mocks.Repository) {
				c.On("SaveRefresh", _contextType, "123", _rtokenType).
					Return(nil)
			},
		},
		{
			name: "ok#2",
			args: args{
				ctx:  context.Background(),
				guid: "fkbhq34btyu1g4yug13ur",
			},
			wantAccess:  "11hg1f1f3v13rv1vf1hbu3rg13rjh11vkh1h",
			wantRefresh: "MTM0YnJpdTFnM3J5ZzEzcnkxM3l1cnYxdW92cg==",
			genMock: func(c *mocks.GeneratorService) {
				c.On("AccessToken", _contextType, "fkbhq34btyu1g4yug13ur", []byte("123"), accessTTL).
					Return("11hg1f1f3v13rv1vf1hbu3rg13rjh11vkh1h", int64(345542514), nil)
				c.On("RefreshToken", _contextType, refreshTTL).
					Return("134briu1g3ryg13ry13yurv1uovr", int64(5243141), nil)
			},
			repoMock: func(c *mocks.Repository) {
				c.On("SaveRefresh", _contextType, "fkbhq34btyu1g4yug13ur", _rtokenType).
					Return(nil)
			},
		},
		{
			name: "AccessError",
			args: args{
				ctx:  context.Background(),
				guid: "",
			},
			genMock: func(c *mocks.GeneratorService) {
				c.On("AccessToken", _contextType, "", []byte("123"), accessTTL).
					Return("", int64(0), constants.ErrInvalidGUID)
			},
			repoMock: func(c *mocks.Repository) {
			},
			wantErr: constants.ErrInvalidGUID,
		},
		{
			name: "RefreshError",
			args: args{
				ctx:  context.Background(),
				guid: "qkefkq",
			},
			genMock: func(c *mocks.GeneratorService) {
				c.On("AccessToken", _contextType, "qkefkq", []byte("123"), accessTTL).
					Return("11hg1f1f3v13rv1vf1hbu3rg13rjh11vkh1h", int64(345542514), nil)
				c.On("RefreshToken", _contextType, refreshTTL).
					Return("", int64(0), constants.ErrGenerateUUID)
			},
			repoMock: func(c *mocks.Repository) {
			},
			wantErr: constants.ErrGenerateUUID,
		},
		{
			name: "RepoError",
			args: args{
				ctx:  context.Background(),
				guid: "kl21rlk",
			},
			genMock: func(c *mocks.GeneratorService) {
				c.On("AccessToken", _contextType, "kl21rlk", []byte("123"), accessTTL).
					Return("11hg1f1f3v13rv1vf1hbu3rg13rjh11vkh1h", int64(345542514), nil)
				c.On("RefreshToken", _contextType, refreshTTL).
					Return("134briu1g3ryg13ry13yurv1uovr", int64(5243141), nil)
			},
			repoMock: func(c *mocks.Repository) {
				c.On("SaveRefresh", _contextType, "kl21rlk", _rtokenType).
					Return(errors.New("repoconstants.Error"))
			},
			wantErr: constants.ErrRepository,
		},
	}

	logger, err := lib.NewLogger()
	if err != nil {
		t.Fatalf("can't create Logger instance: %v", err)
	}

	tm := &TokenManager{
		logger:     logger,
		key:        []byte("123"),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := mocks.NewGeneratorService(t)
			repo := mocks.NewRepository(t)
			tm.repository = repo
			tm.generator = gen
			tt.genMock(gen)
			tt.repoMock(repo)

			gotAccess, gotRefresh, err := tm.GetTokens(tt.args.ctx, tt.args.guid)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetTokens() err %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAccess != tt.wantAccess {
				t.Errorf("GetTokens() gotAccess = %v, want %v", gotAccess, tt.wantAccess)
			}
			if gotRefresh != tt.wantRefresh {
				t.Errorf("GetTokens() gotRefresh = %v, want %v", gotRefresh, tt.wantRefresh)
			}
		})
	}
}

func TestTokenManager_RefreshTokens(t *testing.T) {

	const (
		accessTTL  = time.Minute
		refreshTTL = time.Hour
	)

	type args struct {
		ctx  context.Context
		guid string
	}
	tests := []struct {
		name        string
		args        args
		wantAccess  string
		wantRefresh string
		genMock     genMock
		repoMock    repoMock
		wantErr     error
	}{
		{
			name: "ok",
			args: args{
				ctx:  context.Background(),
				guid: "MTIz",
			},
			wantAccess:  "123",
			wantRefresh: "MTIz",
			genMock: func(c *mocks.GeneratorService) {
				c.On("AccessToken", _contextType, "123", []byte("123"), accessTTL).
					Return("123", int64(123), nil)
				c.On("RefreshToken", _contextType, refreshTTL).
					Return("123", int64(123), nil)
			},
			repoMock: func(c *mocks.Repository) {
				c.On("GetRefTokenAndGUID", _contextType, _stringType).
					Return("123", models.RefreshToken{
						RefreshTokenBCrypt: "xxx",
						Exp:                time.Now().Add(refreshTTL).Unix(),
					}, nil)

				c.On("SaveRefresh", _contextType, "123", _rtokenType).
					Return(nil)
			},
		},
		{
			name: "ok#2",
			args: args{
				ctx:  context.Background(),
				guid: "ZmtiaHEzNGJ0eXUxZzR5dWcxM3Vy",
			},
			wantAccess:  "11hg1f1f3v13rv1vf1hbu3rg13rjh11vkh1h",
			wantRefresh: "MTM0YnJpdTFnM3J5ZzEzcnkxM3l1cnYxdW92cg==",
			genMock: func(c *mocks.GeneratorService) {
				c.On("AccessToken", _contextType, "fkbhq34btyu1g4yug13ur", []byte("123"), accessTTL).
					Return("11hg1f1f3v13rv1vf1hbu3rg13rjh11vkh1h", int64(345542514), nil)
				c.On("RefreshToken", _contextType, refreshTTL).
					Return("134briu1g3ryg13ry13yurv1uovr", int64(5243141), nil)
			},
			repoMock: func(c *mocks.Repository) {
				c.On("GetRefTokenAndGUID", _contextType, _stringType).
					Return("fkbhq34btyu1g4yug13ur", models.RefreshToken{
						RefreshTokenBCrypt: "xxx",
						Exp:                time.Now().Add(refreshTTL).Unix(),
					}, nil)

				c.On("SaveRefresh", _contextType, "fkbhq34btyu1g4yug13ur", _rtokenType).
					Return(nil)
			},
		},
		{
			name: "expired",
			args: args{
				ctx:  context.Background(),
				guid: "ZmtiaHEzNGJ0eXUxZzR5dWcxM3Vy",
			},
			genMock: func(c *mocks.GeneratorService) {
			},
			repoMock: func(c *mocks.Repository) {
				c.On("GetRefTokenAndGUID", _contextType, _stringType).
					Return("fkbhq34btyu1g4yug13ur", models.RefreshToken{
						RefreshTokenBCrypt: "xxx",
						Exp:                time.Now().Add(-refreshTTL).Unix(),
					}, nil)
			},
			wantErr: constants.ErrTokenExpired,
		},
		{
			name: "base64Error",
			args: args{
				ctx:  context.Background(),
				guid: "jkb3rn23lrjn23jlrnj23nrj2b3jrb23jrb2jblj2bfjbj4b2j3brj2b3rb2lj3fb23jrb23jrbj23br;j23b;rk32r;fb2kebjb4j2brjf2bh2b4fhb24hbrfh2hd2fhjfvh",
			},
			genMock: func(c *mocks.GeneratorService) {
			},
			repoMock: func(c *mocks.Repository) {
			},
			wantErr: constants.ErrInvalidToken,
		},
		{
			name: "notFound",
			args: args{
				ctx:  context.Background(),
				guid: "MTM0YnJpdTFnM3J5ZzEzcnkxM3l1cnYxdW92cg==",
			},
			genMock: func(c *mocks.GeneratorService) {
			},
			repoMock: func(c *mocks.Repository) {
				c.On("GetRefTokenAndGUID", _contextType, _stringType).
					Return("",
						models.RefreshToken{}, errors.Join(constants.ErrNotFound, constants.ErrRepository))
			},
			wantErr: constants.ErrNotFound,
		},
		{
			name: "RepoError",
			args: args{
				ctx:  context.Background(),
				guid: "MTM0YnJpdTFnM3J5ZzEzcnkxM3l1cnYxdW92cg==",
			},
			genMock: func(c *mocks.GeneratorService) {},
			repoMock: func(c *mocks.Repository) {
				c.On("GetRefTokenAndGUID", _contextType, _stringType).
					Return("",
						models.RefreshToken{}, constants.ErrRepository)
			},
			wantErr: constants.ErrRepository,
		},
	}

	logger, err := lib.NewLogger()
	if err != nil {
		t.Fatalf("can't create Logger instance: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := mocks.NewGeneratorService(t)
			repo := mocks.NewRepository(t)

			tm := &TokenManager{
				repository: repo,
				logger:     logger,
				key:        []byte("123"),
				accessTTL:  accessTTL,
				refreshTTL: refreshTTL,
				generator:  gen,
			}
			tt.genMock(gen)
			tt.repoMock(repo)

			gotAccess, gotRefresh, err := tm.RefreshTokens(tt.args.ctx, tt.args.guid)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("RefreshTokens() err %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAccess != tt.wantAccess {
				t.Errorf("RefreshTokens() gotAccess = %v, want %v", gotAccess, tt.wantAccess)
			}
			if gotRefresh != tt.wantRefresh {
				t.Errorf("RefreshTokens() gotRefresh = %v, want %v", gotRefresh, tt.wantRefresh)
			}
		})
	}
}
