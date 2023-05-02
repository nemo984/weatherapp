package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepositary struct {
	collection *mongo.Collection
}

func NewRepositary(mongodb *mongo.Database) *UserRepositary {
	return &UserRepositary{collection: mongodb.Collection("users")}
}

func (repo *UserRepositary) UpsertUser(ctx context.Context, user *User) (User, error) {
	if user.ID == primitive.NilObjectID {
		user.ID = primitive.NewObjectID()
	}

	filter := bson.M{
		"device_id": user.DeviceID,
	}
	update := bson.M{
		"$set": user,
	}

	var usr User
	err := repo.collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetUpsert(true), options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&usr)
	if err != nil {
		return usr, err
	}
	return usr, nil
}
