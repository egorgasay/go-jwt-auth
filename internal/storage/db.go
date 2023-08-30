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
	_refreshHash = "refresh_hash"
	_tokens      = "tokens"
	_guid        = "guid"
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

// SaveToken saves a token.
func (d Database) SaveTokenData(ctx context.Context, t models.TokenData) error {
	if _, err := d.db.Collection(_tokens).InsertOne(ctx, t); err != nil {
		return fmt.Errorf("can't insert token: %v", err)
	}

	return nil
}

// GetTokensDataByGUID retrieves the access and refresh tokens for a given GUID.
func (d Database) GetTokensDataByGUID(ctx context.Context, guid string) (t []models.TokenData, err error) {
	filter := bson.D{{_guid, guid}}
	cur, err := d.db.Collection(_tokens).Find(ctx, filter)
	if err != nil {
		return t, err
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		errClose := cur.Close(ctx)
		if errClose != nil && err == nil {
			err = errClose
		}
	}(cur, ctx)

	if err := cur.All(ctx, &t); err != nil {
		return nil, err
	}

	if len(t) == 0 {
		return t, constants.ErrNotFound
	}

	return t, nil
}

// DeleteTokenData deletes a token.
func (d Database) DeleteTokenData(ctx context.Context, guid, hash string) error {
	filter := bson.D{{_guid, guid}, {_refreshHash, hash}}
	_, err := d.db.Collection(_tokens).DeleteOne(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return constants.ErrNotFound
		}
		return err
	}

	return nil
}
