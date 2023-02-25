package gonosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Repository interface {
	InsertOne(ctx context.Context, storage Data, models interface{}) (*mongo.InsertOneResult, error)
	GetByFilter(ctx context.Context, storage Data, key, value string) (*mongo.Cursor, error)
	Get(ctx context.Context, storage Data, id string) *mongo.SingleResult
	Delete(ctx context.Context, storage Data, id string) (*mongo.DeleteResult, error)
	Update(ctx context.Context, storage Data, id string, updateDocument interface{}) (*mongo.UpdateResult, error)
}

func InsertOne(ctx context.Context, storage Data, models interface{}) (*mongo.InsertOneResult, error) {
	result, err := storage.Collection.InsertOne(ctx, models)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetByFilter(ctx context.Context, storage Data, key, value string) (*mongo.Cursor, error) {
	filter := bson.M{key: value}
	cursor, err := storage.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func Get(ctx context.Context, storage Data, id string) *mongo.SingleResult {
	// TODO add findOneOptions in method
	primitiveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// TODO return err? log?
		return nil
	}
	return storage.Collection.FindOne(ctx, bson.M{"_id": primitiveID})
}

func Delete(ctx context.Context, storage Data, id string) (*mongo.DeleteResult, error) {
	primitiveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return storage.Collection.DeleteOne(ctx, bson.M{"_id": primitiveID})
}

func Update(ctx context.Context, storage Data, id string, updateDocument interface{}) (*mongo.UpdateResult, error) {
	primitiveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return storage.Collection.UpdateOne(ctx, bson.M{"_id": primitiveID}, bson.M{"$set": updateDocument})
}
