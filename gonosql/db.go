package gonosql

import (
	"context"
	"fmt"
	"sync"

	"github.com/matiasnu/go-jopit-toolkit/goutils/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	data       *Data
	once       sync.Once
	collection *mongo.Collection
)

type Data struct {
	DB    *mongo.Client
	Error error
}

type JopitDBConfig struct {
	Username   string
	Password   string
	Host       string
	Port       int
	Database   string
	Collection string
}

// Close closes the resources used by data.
func (d *Data) Close(ctx context.Context) {
	if data == nil {
		return
	}

	if err := data.DB.Disconnect(ctx); err != nil {
		logger.Errorf("Error disconect DB", err)
	}
	logger.Debugf("Connection close sucessfully")
}

func NewNoSQL(jopitDBConfig JopitDBConfig) *Data {
	once.Do(func() {
		InitNoSQL(jopitDBConfig)
	})

	return data
}

func GetConnection(jopitDBConfig JopitDBConfig) (*mongo.Client, error) {
	host := jopitDBConfig.Host
	port := jopitDBConfig.Port

	credential := options.Credential{
		Username: jopitDBConfig.Username,
		Password: jopitDBConfig.Password,
	}
	clientOpts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", host, port)).SetAuth(credential)
	return mongo.Connect(context.TODO(), clientOpts)
}

func InitNoSQL(jopitDBConfig JopitDBConfig) {
	var errDB error
	db, err := GetConnection(jopitDBConfig)
	if err != nil {
		errDB = fmt.Errorf("Error NoSQL connection: %s", err)
	}
	// Check the connections
	if err = db.Ping(context.TODO(), nil); err != nil {
		errDB = fmt.Errorf("Error NoSQL connection: %s", err)
	}
	collection = data.DB.Database(jopitDBConfig.Database).Collection(jopitDBConfig.Collection)
	data = &Data{
		DB:    db,
		Error: errDB,
	}
}
