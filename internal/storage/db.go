package storage

import (
	"context"
	"fmt"
	"go-jwt-auth/internal/domains"
	"go-jwt-auth/internal/lib"
	"go-jwt-auth/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	_user               = "user"
	_refreshTokenBcrypt = "refresh_token.refresh_token_bcrypt"
)

type Database struct {
	db lib.Database
}

func NewDatabase(db lib.Database) domains.Database {
	return Database{db: db}
}

func (d Database) Close() error {
	if err := d.db.Client().Disconnect(context.Background()); err != nil {
		return fmt.Errorf("can't disconnect: %v", err)
	}
	return nil
}

func (d Database) SaveRefresh(ctx context.Context, guid string, refresh models.RefreshToken) error {
	if _, err := d.db.Collection(_user).InsertOne(ctx, models.User{
		GUID:         guid,
		RefreshToken: refresh,
	}); err != nil {
		return err
	}

	return nil
}

func (d Database) GetRefTokenAndGUID(ctx context.Context, refresh string) (guid string, rt models.RefreshToken, err error) {
	var us models.User
	filter := bson.D{{_refreshTokenBcrypt, refresh}}
	err = d.db.Collection(_user).FindOneAndDelete(ctx, filter).Decode(&us)
	if err != nil {
		return "", rt, err
	}

	return us.GUID, us.RefreshToken, nil
}
