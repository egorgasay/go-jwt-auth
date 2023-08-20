package storage

import (
	"context"
	"fmt"
	"github.com/egorgasay/dockerdb/v3"
	"go-jwt-auth/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

type Storage struct {
	db     *mongo.Database
	logger *zap.Logger
}

type Config struct {
	DSN string `json:"dsn"`
}

const (
	_autoUp             = "AUTO_UP"
	_user               = "user"
	_dbName             = "jwt-auth"
	_refreshTokenBcrypt = "refresh_token.refresh_token_bcrypt"
)

func New(conf Config, logger *zap.Logger) (*Storage, error) {
	ctx := context.Background()
	var err error
	var client *mongo.Client

	switch conf.DSN {
	case _autoUp:
		vdbConf := dockerdb.EmptyConfig().DBName(_dbName).NoSQL(func(c dockerdb.Config) (stop bool) {
			client, err = mongo.Connect(ctx, options.Client().ApplyURI(conf.DSN))
			if err != nil {
				logger.Error("can't connect to mongodb", zap.Error(err))
			}

			return client.Ping(ctx, nil) == nil
		}, 30, 2*time.Second).Build()

		_, err := dockerdb.New(context.Background(), vdbConf)
		if err != nil {
			return nil, fmt.Errorf("can't up or connect to dockerdb: %v", err)
		}

		return &Storage{logger: logger, db: client.Database(_dbName)}, nil
	case "":
		return nil, fmt.Errorf("DSN shouldn't be empty")
	default:
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(conf.DSN))
		if err != nil {
			return nil, err
		}

		return &Storage{logger: logger, db: client.Database(_dbName)}, nil
	}
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
