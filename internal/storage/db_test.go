package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/egorgasay/dockerdb/v3"
	"go-jwt-auth/internal/lib"
	"go-jwt-auth/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

func TestDatabase_SaveToken(t *testing.T) {
	type args struct {
		ctx context.Context
		t   models.TokenData
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				t: models.TokenData{
					GUID: "123",
				},
			},
			wantErr: nil,
		},
		{
			name: "ok#2",
			args: args{
				ctx: context.Background(),
				t: models.TokenData{
					GUID: "yguf67d7rr7di",
				},
			},
			wantErr: nil,
		},
	}

	ctx := context.Background()
	vdb, client, err := upMongo(ctx, t)
	defer func() {
		err = vdb.Clear(ctx)
		if err != nil {
			t.Errorf("can't clear container, possible container leak and wrong results in the future tests: %v", err)
		}
	}()

	d := Database{
		db: lib.Database{Database: client.Database(lib.DBName)},
	}

	collection := client.Database(lib.DBName).Collection(_tokens)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.SaveTokenData(tt.args.ctx, tt.args.t); !errors.Is(err, tt.wantErr) {
				t.Errorf("SaveToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			r := collection.FindOne(ctx, bson.D{{_guid, tt.args.t.GUID}})
			if err := r.Err(); err != nil {
				t.Fatalf("FindOne error: %v", err)
			}

			var td models.TokenData
			if err := r.Decode(&td); err != nil {
				t.Fatalf("Decode error: %v", err)
			}

			assert.DeepEqual(t, tt.args.t, td)
		})
	}
}

func upMongo(ctx context.Context, t *testing.T) (vdb *dockerdb.VDB, cl *mongo.Client, err error) {
	cfg := dockerdb.EmptyConfig().Vendor("mongo").DBName("SaveTokenData").
		NoSQL(func(c dockerdb.Config) (stop bool) {
			dsn := fmt.Sprintf("mongodb://127.0.0.1:%s", c.GetActualPort())
			opt := options.Client()
			opt.ApplyURI(dsn).SetTimeout(1 * time.Second)

			cl, err = mongo.Connect(ctx, opt)
			if err != nil {
				t.Log("can't connect to mongodb", zap.Error(err))
				return false
			}

			if err := cl.Ping(ctx, nil); err != nil {
				t.Log("can't ping mongodb", zap.Error(err))
				return false
			}

			t.Logf("connected to mongodb at %s", dsn)
			return true
		}, 30, 2*time.Second).PullImage().StandardDBPort("27017").Build()

	vdb, err = dockerdb.New(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	return vdb, cl, err
}
