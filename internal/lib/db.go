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
	DBName  = "jwt-auth-db"
	_autoUp = "AUTO_UP"
)

func NewDatabase(conf Config) (db Database, err error) {
	ctx := context.Background()
	var client *mongo.Client

	switch conf.Storage.DatabaseDSN {
	case _autoUp:
		vdbConf := dockerdb.EmptyConfig().Vendor("mongo").DBName(DBName).
			NoSQL(func(c dockerdb.Config) (stop bool) {
				dsn := fmt.Sprintf("mongodb://127.0.0.1:%s", c.GetActualPort())
				opt := options.Client()
				opt.ApplyURI(dsn).SetTimeout(1 * time.Second)

				client, err = mongo.Connect(ctx, opt)
				if err != nil {
					log.Println("can't connect to mongodb", zap.Error(err))
					return false
				}

				if err := client.Ping(ctx, nil); err != nil {
					log.Println("can't ping mongodb", zap.Error(err))
					return false
				}

				return true
			}, 30, 2*time.Second).PullImage().StandardDBPort("27017").Build()

		_, err := dockerdb.New(context.Background(), vdbConf)
		if err != nil {
			return db, fmt.Errorf("can't up or connect to dockerdb: %v", err)
		}

		return Database{Database: client.Database(DBName)}, nil
	case "":
		return db, fmt.Errorf("DatabaseDSN shouldn't be empty")
	default:
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(conf.Storage.DatabaseDSN))
		if err != nil {
			return db, err
		}

		return Database{Database: client.Database(DBName)}, nil
	}
}
