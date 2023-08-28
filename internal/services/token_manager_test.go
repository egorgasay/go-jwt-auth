package services

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"go-jwt-auth/internal/constants"
	"go-jwt-auth/internal/domains/mocks"
	"go-jwt-auth/internal/lib"
	"go-jwt-auth/internal/models"
	"math"
	"testing"
	"time"
)

type (
	repoMock func(c *mocks.Repository)
	genMock  func(c *mocks.GeneratorService)
)

var (
	_contextType = mock.AnythingOfType("context.backgroundCtx")
	_rtokenType  = mock.AnythingOfType("models.TokenData")
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
			wantAccess:  "MTIz",
			wantRefresh: "MTIz",
			genMock: func(c *mocks.GeneratorService) {
				c.On("AccessToken", _contextType, "123", []byte("123")).
					Return("123", nil)
				c.On("RefreshToken", _contextType).
					Return("123", nil)
			},
			repoMock: func(c *mocks.Repository) {
				c.On("SaveToken", _contextType, _rtokenType).
					Return(nil)
			},
		},
		{
			name: "ok#2",
			args: args{
				ctx:  context.Background(),
				guid: "fkbhq34btyu1g4yug13ur",
			},
			wantAccess:  "MTFoZzFmMWYzdjEzcnYxdmYxaGJ1M3JnMTNyamgxMXZraDFo",
			wantRefresh: "MTM0YnJpdTFnM3J5ZzEzcnkxM3l1cnYxdW92cg==",
			genMock: func(c *mocks.GeneratorService) {
				c.On("AccessToken", _contextType, "fkbhq34btyu1g4yug13ur", []byte("123")).
					Return("11hg1f1f3v13rv1vf1hbu3rg13rjh11vkh1h", nil)
				c.On("RefreshToken", _contextType).
					Return("134briu1g3ryg13ry13yurv1uovr", nil)
			},
			repoMock: func(c *mocks.Repository) {
				c.On("SaveToken", _contextType, _rtokenType).
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
				c.On("AccessToken", _contextType, "", []byte("123")).
					Return("", constants.ErrInvalidGUID)
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
				c.On("AccessToken", _contextType, "qkefkq", []byte("123")).
					Return("11hg1f1f3v13rv1vf1hbu3rg13rjh11vkh1h", nil)
				c.On("RefreshToken", _contextType).
					Return("", constants.ErrGenerateToken)
			},
			repoMock: func(c *mocks.Repository) {
			},
			wantErr: constants.ErrGenerateToken,
		},
		{
			name: "RepoError",
			args: args{
				ctx:  context.Background(),
				guid: "kl21rlk",
			},
			genMock: func(c *mocks.GeneratorService) {
				c.On("AccessToken", _contextType, "kl21rlk", []byte("123")).
					Return("11hg1f1f3v13rv1vf1hbu3rg13rjh11vkh1h", nil)
				c.On("RefreshToken", _contextType).
					Return("134briu1g3ryg13ry13yurv1uovr", nil)
			},
			repoMock: func(c *mocks.Repository) {
				c.On("SaveToken", _contextType, _rtokenType).
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
		ctx     context.Context
		access  string
		refresh string
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
				ctx:     context.Background(),
				access:  "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5Sm5kV2xrSWpvaWFXdHFJbjAuUl95MlAtRHNKQUNZTHBnRG1BLXRBN1FUVnFrZU90MDRKaGxGQ2Z6NjRSbmRRSUlLczVjWW1mTGtFd3MzUW1xWDhSNEc4TkJkaER4T2s4ZVNGZGpvM3c=",
				refresh: "NTA0YmNmMmEtNDVkZi0xMWVlLWE0Y2ItMDYzMGY4YzRkMDRj",
			},
			wantAccess:  "bHdrZW5rZm5xbmZxa3dm",
			wantRefresh: "andybmZid2plYmtmcWh2ZWZxaGo=",
			genMock: func(c *mocks.GeneratorService) {
				c.On("AccessToken", _contextType, "ikj", []byte("123")).
					Return("lwkenkfnqnfqkwf", nil)
				c.On("RefreshToken", _contextType).
					Return("jwrnfbwjebkfqhvefqhj", nil)
			},
			repoMock: func(c *mocks.Repository) {
				c.On("GetTokensDataByGUID", _contextType, "ikj").
					Return([]models.TokenData{
						{
							GUID:        "ikj",
							RefreshHash: "$2a$10$Rct7JqhZDVzFGdRgG0caZurIrkyUe893JhvB0.8eXO.CKOLGppEDy",
							RefreshExp:  math.MaxInt,
							AccessExp:   math.MaxInt,
						},
					}, nil)
				c.On("DeleteTokenData", _contextType, "ikj", "$2a$10$Rct7JqhZDVzFGdRgG0caZurIrkyUe893JhvB0.8eXO.CKOLGppEDy").
					Return(nil)
				c.On("SaveToken", _contextType, _rtokenType).
					Return(nil)
			},
		},
		{
			name: "ok#2",
			args: args{
				ctx:     context.Background(),
				access:  "ZXlKaGJHY2lPaUpJVXpVeE1pSXNJblI1Y0NJNklrcFhWQ0o5LmV5Sm5kV2xrSWpvaWEzZG1kMlVpZlEudWhicFIwRkx3d3VBY3J5eWhRdVJnNlpwVzBNelc1ako2VnhXMlRTZGNxR0o0Vm5oNnBRdk5fN1lBSHlEbEt5eTJyVGg4NXhyZGM0SHlERlJ5elZYNEE=",
				refresh: "YTQxZjIwYjAtNDVlMC0xMWVlLWE0Y2ItMDYzMGY4YzRkMDRj",
			},
			wantAccess:  "andmMzczYjNqaGRiajMxYnJ1",
			wantRefresh: "bjM3Z2ZiMnUzN2Z1MmY=",
			genMock: func(c *mocks.GeneratorService) {
				c.On("AccessToken", _contextType, "kwfwe", []byte("123")).
					Return("jwf373b3jhdbj31bru", nil)
				c.On("RefreshToken", _contextType).
					Return("n37gfb2u37fu2f", nil)
			},
			repoMock: func(c *mocks.Repository) {
				c.On("GetTokensDataByGUID", _contextType, "kwfwe").
					Return([]models.TokenData{
						{
							GUID:        "kwfwe",
							RefreshHash: "$2a$10$VEjOdbltCL7QRByQ1g//4e4KseOMXwvEziIMv2ULi0/8vIuY0394S",
							RefreshExp:  math.MaxInt,
							AccessExp:   math.MaxInt,
						},
					}, nil)
				c.On("DeleteTokenData", _contextType, "kwfwe", "$2a$10$VEjOdbltCL7QRByQ1g//4e4KseOMXwvEziIMv2ULi0/8vIuY0394S").
					Return(nil)
				c.On("SaveToken", _contextType, _rtokenType).
					Return(nil)
			},
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

			gotAccess, gotRefresh, err := tm.RefreshTokens(tt.args.ctx, tt.args.access, tt.args.refresh)
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
