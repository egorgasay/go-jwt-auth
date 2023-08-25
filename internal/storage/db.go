package storage

import (
	"context"
	"errors"
	"fmt"
	"go-jwt-auth/internal/constants"
	"go-jwt-auth/internal/domains"
	"go-jwt-auth/internal/lib"
	"go-jwt-auth/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	_user               = "user"
	_refreshTokenBcrypt = "refresh_token.refresh_token_bcrypt"
)

// Database is a struct that contains a database.
type Database struct {
	db lib.Database
}

// NewDatabase creates a new instance of Database.
func NewDatabase(db lib.Database) domains.Database {
	return Database{db: db}
}

// Close closes the database.
func (d Database) Close() error {
	if err := d.db.Client().Disconnect(context.Background()); err != nil {
		return fmt.Errorf("can't disconnect: %v", err)
	}
	return nil
}

// SaveRefresh saves a refresh token to the database.
func (d Database) SaveRefresh(ctx context.Context, guid string, refresh models.RefreshToken) error {
	if _, err := d.db.Collection(_user).InsertOne(ctx, models.User{
		GUID:         guid,
		RefreshToken: refresh,
	}); err != nil {
		return err
	}

	return nil
}

// GetRefTokenAndGUID retrieves the reference token and GUID from the database.
func (d Database) GetRefTokenAndGUID(ctx context.Context, refresh string) (guid string, rt models.RefreshToken, err error) {
	var us models.User
	filter := bson.D{{_refreshTokenBcrypt, refresh}}
	err = d.db.Collection(_user).FindOneAndDelete(ctx, filter).Decode(&us)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", models.RefreshToken{}, constants.ErrNotFound
		}

		return "", rt, err
	}

	return us.GUID, us.RefreshToken, nil
}
