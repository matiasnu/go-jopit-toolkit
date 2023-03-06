package gonosql

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/matiasnu/go-jopit-toolkit/goutils/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI = "mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority"
)

var (
	data *Data
	once sync.Once
)

type Data struct {
	DB       *mongo.Client
	Database *mongo.Database
	Error    error
}

type Config struct {
	Username string
	Password string
	Host     string
	Database string
}

// Close closes the resources used by data.
func (d *Data) Close() {
	if data == nil {
		return
	}

	if err := data.DB.Disconnect(context.Background()); err != nil {
		logger.Errorf("Error disconect DB", err)
	}
	logger.Debugf("Connection close sucessfully")
}

func (d *Data) NewCollection(collection string) *mongo.Collection {
	return d.Database.Collection(collection)
}

func NewNoSQL(jopitDBConfig Config) *Data {
	once.Do(func() {
		InitNoSQL(jopitDBConfig)
	})

	return data
}

func GetConnection(jopitDBConfig Config) (*mongo.Client, error) {
	host := jopitDBConfig.Host
	username := jopitDBConfig.Username
	password := jopitDBConfig.Password

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts := options.Client().
		ApplyURI(fmt.Sprintf(mongoURI, username, password, host)).
		SetServerAPIOptions(serverAPIOptions)
	// TODO pass context? (Analyze)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mongo.Connect(ctx, clientOpts)
}

func InitNoSQL(jopitDBConfig Config) {
	var (
		errDB    error
		database *mongo.Database
	)
	db, err := GetConnection(jopitDBConfig)
	if err != nil {
		errDB = fmt.Errorf("Error NoSQL connection: %s", err)
	} else {
		// Check the connections
		if err = db.Ping(context.TODO(), nil); err != nil {
			errDB = fmt.Errorf("Error NoSQL connection: %s", err)
		}
		database = db.Database(jopitDBConfig.Database)
	}

	data = &Data{
		DB:       db,
		Error:    errDB,
		Database: database,
	}
}
