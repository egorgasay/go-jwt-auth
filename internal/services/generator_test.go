package services

import (
	"context"
	"errors"
	"go-jwt-auth/internal/constants"
	"go-jwt-auth/internal/lib"
	"testing"
	"time"
)

func TestGeneratorService_AccessToken(t *testing.T) {
	type res struct {
		Access string
		Exp    int64
		Err    error
	}
	type args struct {
		ctx       context.Context
		guid      string
		key       []byte
		accessTTL time.Duration
	}
	tests := []struct {
		name string
		args args
		want res
	}{
		{
			name: "ok",
			args: args{
				ctx:       context.Background(),
				guid:      "123",
				key:       []byte("123"),
				accessTTL: time.Minute,
			},
			want: res{},
		},
		{
			name: "ok#2",
			args: args{
				ctx:       context.Background(),
				guid:      "g28f123gvud1vuy31vry3rv3",
				key:       []byte("1u4fy1vyv1uv1ey"),
				accessTTL: time.Hour,
			},
			want: res{},
		},
		{
			name: "ErrInvalidGUID",
			args: args{
				ctx:       context.Background(),
				guid:      "",
				key:       []byte("qfeqjfkj"),
				accessTTL: time.Hour,
			},
			want: res{
				Err: constants.ErrInvalidGUID,
			},
		},
	}

	logger, err := lib.NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GeneratorService{
				logger: logger,
			}
			var got res
			_, _, got.Err = g.AccessToken(tt.args.ctx, tt.args.guid, tt.args.key, tt.args.accessTTL)
			if !errors.Is(got.Err, tt.want.Err) {
				t.Errorf("AccessToken() error = %v, wantErr %v", got.Err, tt.want.Err)
			}
		})
	}
}

func TestGeneratorService_RefreshToken(t *testing.T) {
	type res struct {
		Access string
		Exp    int64
		Err    error
	}
	type args struct {
		ctx        context.Context
		refreshTTL time.Duration
	}
	tests := []struct {
		name string
		args args
		want res
	}{
		{
			name: "ok",
			args: args{
				ctx:        context.Background(),
				refreshTTL: time.Minute,
			},
			want: res{},
		},
		{
			name: "ok#2",
			args: args{
				ctx:        context.Background(),
				refreshTTL: time.Hour,
			},
			want: res{},
		},
	}

	logger, err := lib.NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GeneratorService{
				logger: logger,
			}
			var got res
			_, _, got.Err = g.RefreshToken(tt.args.ctx, tt.args.refreshTTL)
			if !errors.Is(got.Err, tt.want.Err) {
				t.Errorf("AccessToken() error = %v, wantErr %v", got.Err, tt.want.Err)
			}
		})
	}
}
