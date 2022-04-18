package gosql

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
)

var (
	data *Data
	once sync.Once
)

type Data struct {
	DB *gorm.DB
}

type QueryBuilder struct{}

type JopitDBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

// Close closes the resources used by data.
func (d *Data) Close() error {
	if data == nil {
		return nil
	}

	return data.DB.Close()
}

func NewSQL(jopitDBConfig JopitDBConfig) *Data {
	once.Do(func() {
		InitSQL(jopitDBConfig)
	})

	return data
}

func GetConnection(jopitDBConfig JopitDBConfig) (*gorm.DB, error) {
	username := jopitDBConfig.Username
	password := jopitDBConfig.Password
	host := jopitDBConfig.Host
	port := jopitDBConfig.Port
	database := jopitDBConfig.Database
	connection := username + ":" + password + "@" + "tcp(" + host + ":" + port + ")/" + database + "?charset=utf8&parseTime=true"
	return gorm.Open("mysql", connection)
}

func InitSQL(jopitDBConfig JopitDBConfig) {
	db, err := GetConnection(jopitDBConfig)
	if err != nil {
		fmt.Printf("Error MySQL connection : %s", err)
	}
	data = &Data{
		DB: db,
	}
}
