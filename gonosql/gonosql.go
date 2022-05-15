package gonosql

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Repository interface {
	InsertOne(storage Data, models interface{}) (*mongo.InsertOneResult, error)
	GetByFilter(storage Data, filter bson.M) (*mongo.Cursor, error)
}

func InsertOne(storage Data, models interface{}) (*mongo.InsertOneResult, error) {
	result, err := storage.Collection.InsertOne(context.TODO(), models)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetByFilter(storage Data, filter bson.M) (*mongo.Cursor, error) {
	cursor, err := storage.Collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}
