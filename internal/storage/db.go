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
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	_user               = "user"
	_refreshTokenBcrypt = "token.refresh_hash"
	_guid               = "guid"
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
func (d Database) SaveToken(ctx context.Context, guid string, t models.Token) error {
	if _, err := d.db.Collection(_user).InsertOne(ctx, models.User{
		GUID:  guid,
		Token: t,
	}); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return constants.ErrAlreadyExists
		}
		return err
	}

	return nil
}

// GetRefTokenByGUID retrieves a refresh token from the database.
func (d Database) GetRefTokenByGUID(ctx context.Context, guid string) (t models.Token, err error) {
	var us models.User
	filter := bson.D{{_guid, guid}}
	opts := options.FindOneAndUpdate().SetUpsert(false)
	update := bson.D{{"$set", bson.D{{_refreshTokenBcrypt, ""}}}}
	err = d.db.Collection(_user).FindOneAndUpdate(ctx, filter, update, opts).
		Decode(&us)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return t, constants.ErrNotFound
		}

		return t, err
	}

	return us.Token, nil
}
