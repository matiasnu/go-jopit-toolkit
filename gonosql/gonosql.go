package gonosql

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Repository interface {
	InsertOne(ctx context.Context, storage Data, models interface{}) (*mongo.InsertOneResult, error)
	GetByFilter(ctx context.Context, storage Data, filter bson.M) (*mongo.Cursor, error)
	Get(ctx context.Context, storage Data, id string)
}

func InsertOne(ctx context.Context, storage Data, models interface{}) (*mongo.InsertOneResult, error) {
	result, err := storage.Collection.InsertOne(ctx, models)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetByFilter(ctx context.Context, storage Data, filter bson.M) (*mongo.Cursor, error) {
	cursor, err := storage.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func Get(ctx context.Context, storage Data, id string) *mongo.SingleResult {
	// TODO add findOneOptions in method
	return storage.Collection.FindOne(ctx, bson.M{"_id": id})
}
