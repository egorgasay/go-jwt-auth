package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"go-jwt-auth/internal/constants"
	"go-jwt-auth/internal/domains/mocks"
	"go-jwt-auth/internal/lib"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type tmMock func(c *mocks.TokenManager)

func TestTokenHandler_GetTokens(t *testing.T) {
	type args struct {
		guid string
	}
	tests := []struct {
		name     string
		wantJSON string
		tmMock   tmMock
		args     args
	}{
		{
			name: "ok#1",
			wantJSON: `{
	"access_token": "MTIz",
	"refresh_token": "MTIz"
}`,
			tmMock: func(c *mocks.TokenManager) {
				c.On("GetTokens", mock.Anything, "123").Return("MTIz", "MTIz", nil)
			},
			args: args{
				guid: "123",
			},
		},
		{
			name: "ok#2",
			wantJSON: `{
	"access_token": "qqhqbw18hfqd183hqwvdlgvqgwvdjqvgd",
	"refresh_token": "1ybu3fg178fo26f6ig2d1"
}`,
			tmMock: func(c *mocks.TokenManager) {
				c.On("GetTokens", mock.Anything, "huhqfhqi").
					Return("qqhqbw18hfqd183hqwvdlgvqgwvdjqvgd",
						"1ybu3fg178fo26f6ig2d1", nil)
			},
			args: args{
				guid: "huhqfhqi",
			},
		},
		{
			name: "ErrGenerate",
			wantJSON: `{
	"error": "can't generate\ncan't generate token"
}`,
			tmMock: func(c *mocks.TokenManager) {
				c.On("GetTokens", mock.Anything, "huhqfhqi").
					Return("", "", errors.Join(constants.ErrGenerate, constants.ErrGenerateToken))
			},
			args: args{
				guid: "huhqfhqi",
			},
		},
	}

	logger, err := lib.NewLogger()
	if err != nil {
		t.Fatalf("can't create Logger instance: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := mocks.NewTokenManager(t)
			h := &TokenHandler{
				tokens: tokens,
				logger: logger,
			}
			tt.tmMock(tokens)

			path := "/t"

			r := gin.Default()
			r.GET(path, h.GetTokens)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, path, nil)
			q := req.URL.Query()
			q.Set("guid", tt.args.guid)

			req.URL.RawQuery = q.Encode()

			r.ServeHTTP(w, req)

			if !cmpJSON(tt.wantJSON, w.Body.String()) {
				t.Errorf("want:\n%v\ngot:\n%v", tt.wantJSON, w.Body.String())
				return
			}
		})
	}
}

func cmpJSON(want, got string) bool {
	var m1 = make(map[rune]int)
	var m2 = make(map[rune]int)

	for _, v := range want {
		if v == ' ' || v == '\n' || v == '\t' {
			continue
		}
		m1[v]++
	}
	for _, v := range got {
		if v == ' ' || v == '\n' || v == '\t' {
			continue
		}
		m2[v]++
	}

	return reflect.DeepEqual(m1, m2)
}

func TestTokenHandler_RefreshTokens(t *testing.T) {
	type args struct {
		body   string
		access string
	}
	tests := []struct {
		name     string
		wantJSON string
		tmMock   tmMock
		args     args
	}{
		{
			name: "ok#1",
			wantJSON: `{
	"access_token": "nqbkfbqkjf",
	"refresh_token": "jkqwbfjkqbj"
}`,
			tmMock: func(c *mocks.TokenManager) {
				c.On("RefreshTokens", mock.Anything, "qwmdq", "jqnkfjnq").
					Return("nqbkfbqkjf", "jkqwbfjkqbj", nil)
			},
			args: args{
				access: "qwmdq",
				body: `{
					"refresh_token": "jqnkfjnq"
				}`,
			},
		},
		{
			name: "ErrExpired",
			wantJSON: `{
	"error": "token expired"
}`,
			tmMock: func(c *mocks.TokenManager) {
				c.On("RefreshTokens", mock.Anything, "huhqfhqi", "jqnkfjnq").
					Return("", "", constants.ErrTokenExpired)
			},
			args: args{
				access: "huhqfhqi",
				body: `{
					"refresh_token": "jqnkfjnq"
				}`,
			},
		},
		{
			name: "ErrInvalid",
			wantJSON: `{
	"error": "invalid token"
}`,
			tmMock: func(c *mocks.TokenManager) {
				c.On("RefreshTokens", mock.Anything, "huhqfhqi", "jqnkfjnq").
					Return("", "", constants.ErrInvalidToken)
			},
			args: args{
				access: "huhqfhqi",
				body: `{
					"refresh_token": "jqnkfjnq"
				}`,
			},
		},
	}

	logger, err := lib.NewLogger()
	if err != nil {
		t.Fatalf("can't create Logger instance: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := mocks.NewTokenManager(t)
			h := &TokenHandler{
				tokens: tokens,
				logger: logger,
			}
			tt.tmMock(tokens)

			path := "/t"

			r := gin.Default()
			r.POST(path, h.RefreshTokens)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(tt.args.body))
			q := req.URL.Query()
			req.Header.Set("Authorization", tt.args.access)

			req.URL.RawQuery = q.Encode()

			r.ServeHTTP(w, req)

			if !cmpJSON(tt.wantJSON, w.Body.String()) {
				t.Errorf("want:\n%v\ngot:\n%v", tt.wantJSON, w.Body.String())
				return
			}
		})
	}
}
