package storage

import (
	"context"
	"fmt"
	"go-jwt-auth/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Storage struct {
	db     *mongo.Database
	logger *zap.Logger
}

type Config struct {
	DSN string `json:"dsn"`
}

const (
	_exp                = "exp"
	_guid               = "guid"
	_user               = "user"
	_dbName             = "jwt-auth"
	_refreshTokenBcrypt = "refresh_token.refresh_token_bcrypt"
)

func New(conf Config, logger *zap.Logger) (*Storage, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.DSN))
	if err != nil {
		return nil, err
	}

	return &Storage{logger: logger, db: client.Database(_dbName)}, nil
}

func (s Storage) Close() error {
	if err := s.db.Client().Disconnect(context.Background()); err != nil {
		return fmt.Errorf("can't disconnect: %v", err)
	}
	return nil
}

func (s Storage) SaveRefresh(ctx context.Context, guid string, refresh model.RefreshToken) error {
	if _, err := s.db.Collection(_user).InsertOne(ctx, model.User{
		GUID:         guid,
		RefreshToken: refresh,
	}); err != nil {
		return err
	}

	return nil
}

func (s Storage) GetRefTokenAndGUID(ctx context.Context, refresh string) (guid string, rt model.RefreshToken, err error) {
	var us model.User
	filter := bson.D{{_refreshTokenBcrypt, refresh}}
	err = s.db.Collection(_user).FindOneAndDelete(ctx, filter).Decode(&us)
	if err != nil {
		return "", rt, err
	}

	return us.GUID, us.RefreshToken, nil
}
