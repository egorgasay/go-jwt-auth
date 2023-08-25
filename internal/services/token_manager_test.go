package services

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"go-jwt-auth/internal/domains/mocks"
	"go-jwt-auth/internal/lib"
	"testing"
	"time"
)

type (
	repoMock func(c *mocks.Repository)
	genMock  func(c *mocks.GeneratorService)
)

var _contextType = mock.AnythingOfType("context.backgroundCtx")

func TestTokenManager_GetTokens(t *testing.T) {

	const (
		accessTTL  = time.Minute
		refreshTTL = time.Hour
	)

	var rtokenType = mock.AnythingOfType("models.RefreshToken")

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
				c.On("Access", _contextType, "123", []byte("123"), accessTTL).
					Return("123", int64(123), nil)
				c.On("Refresh", _contextType, refreshTTL).
					Return("123", int64(123), nil)
			},
			repoMock: func(c *mocks.Repository) {
				c.On("SaveRefresh", _contextType, "123", rtokenType).
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
				c.On("Access", _contextType, "fkbhq34btyu1g4yug13ur", []byte("123"), accessTTL).
					Return("11hg1f1f3v13rv1vf1hbu3rg13rjh11vkh1h", int64(345542514), nil)
				c.On("Refresh", _contextType, refreshTTL+1).
					Return("134briu1g3ryg13ry13yurv1uovr", int64(5243141), nil)
			},
			repoMock: func(c *mocks.Repository) {
				c.On("SaveRefresh", _contextType, "fkbhq34btyu1g4yug13ur", rtokenType).
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
				c.On("Access", _contextType, "", []byte("123"), accessTTL).
					Return("", int64(0), ErrInvalidGUID)
			},
			repoMock: func(c *mocks.Repository) {
			},
			wantErr: ErrInvalidGUID,
		},
		{
			name: "RefreshError",
			args: args{
				ctx:  context.Background(),
				guid: "qkefkq",
			},
			genMock: func(c *mocks.GeneratorService) {
				c.On("Access", _contextType, "qkefkq", []byte("123"), accessTTL).
					Return("11hg1f1f3v13rv1vf1hbu3rg13rjh11vkh1h", int64(345542514), nil)
				c.On("Refresh", _contextType, refreshTTL+6).
					Return("", int64(0), ErrGenerateUUID)
			},
			repoMock: func(c *mocks.Repository) {
			},
			wantErr: ErrGenerateUUID,
		},
		{
			name: "RepoError",
			args: args{
				ctx:  context.Background(),
				guid: "kl21rlk",
			},
			genMock: func(c *mocks.GeneratorService) {
				c.On("Access", _contextType, "kl21rlk", []byte("123"), accessTTL).
					Return("11hg1f1f3v13rv1vf1hbu3rg13rjh11vkh1h", int64(345542514), nil)
				c.On("Refresh", _contextType, refreshTTL+10).
					Return("134briu1g3ryg13ry13yurv1uovr", int64(5243141), nil)
			},
			repoMock: func(c *mocks.Repository) {
				c.On("SaveRefresh", _contextType, "kl21rlk", rtokenType).
					Return(errors.New("repo error"))
			},
			wantErr: ErrRepository,
		},
	}

	gen := mocks.NewGeneratorService(t)
	repo := mocks.NewRepository(t)

	logger, err := lib.NewLogger()
	if err != nil {
		t.Fatalf("can't create Logger instance: %v", err)
	}

	tm := &TokenManager{
		repository: repo,
		logger:     logger,
		key:        []byte("123"),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		generator:  gen,
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.genMock(gen)
			tt.repoMock(repo)

			tm.refreshTTL += time.Duration(i)

			gotAccess, gotRefresh, err := tm.GetTokens(tt.args.ctx, tt.args.guid)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetTokens() error = %v, wantErr %v", err, tt.wantErr)
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
