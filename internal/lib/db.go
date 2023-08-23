package lib

import (
	"context"
	"fmt"
	"github.com/egorgasay/dockerdb/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"log"
	"time"
)

type Database struct {
	*mongo.Database
}

const (
	_dbName = "jwt-auth"
	_autoUp = "AUTO_UP"
)

func NewDatabase(conf Config) (db Database, err error) {
	ctx := context.Background()
	var client *mongo.Client

	switch conf.StorageConfig.DatabaseDSN {
	case _autoUp:
		vdbConf := dockerdb.EmptyConfig().Vendor("mongo").DBName(_dbName).
			NoSQL(func(c dockerdb.Config) (stop bool) {
				client, err = mongo.Connect(ctx, options.Client().ApplyURI(conf.StorageConfig.DatabaseDSN))
				if err != nil {
					log.Println("can't connect to mongodb", zap.Error(err))
				}

				return client.Ping(ctx, nil) == nil
			}, 30, 2*time.Second).Build()

		_, err := dockerdb.New(context.Background(), vdbConf)
		if err != nil {
			return db, fmt.Errorf("can't up or connect to dockerdb: %v", err)
		}

		return Database{Database: client.Database(_dbName)}, nil
	case "":
		return db, fmt.Errorf("DatabaseDSN shouldn't be empty")
	default:
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(conf.StorageConfig.DatabaseDSN))
		if err != nil {
			return db, err
		}

		return Database{Database: client.Database(_dbName)}, nil
	}
}
