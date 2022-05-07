package gonosql

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	InsertOne(models interface{}) (*mongo.InsertOneResult, error)
	GetByFilter(filter primitive.M) (*mongo.Cursor, error)
}

func InsertOne(models interface{}) (*mongo.InsertOneResult, error) {
	result, err := collection.InsertOne(context.TODO(), models)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetByFilter(filter primitive.M) (*mongo.Cursor, error) {
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}
