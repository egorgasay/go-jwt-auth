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
func (d Database) SaveToken(ctx context.Context, t models.TokenData) error {
	opts := &options.ReplaceOptions{}
	opts.SetUpsert(true)
	filter := bson.D{{_guid, t.GUID}, {_refreshHash, ""}}
	if _, err := d.db.Collection(_tokens).ReplaceOne(ctx, filter, t, opts); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return constants.ErrAlreadyExists
		}
		return err
	}

	return nil
}

// GetTokensDataByGUID retrieves the access and refresh tokens for a given GUID.
func (d Database) GetTokensDataByGUID(ctx context.Context, guid string) (t []models.TokenData, err error) {
	filter := bson.D{{_guid, guid}}
	cur, err := d.db.Collection(_tokens).Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return t, constants.ErrNotFound
		}

		return t, err
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err = cur.Close(ctx)
	}(cur, ctx)

	for cur.Next(ctx) {
		var td models.TokenData
		err := cur.Decode(&td)
		if err != nil {
			return t, err
		}
		t = append(t, td)
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
