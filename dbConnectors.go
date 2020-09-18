package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id, omitempty"`
	Email    string
	Password string
	Data     interface{}
}

func loadUserByEmail(email string, collection *mongo.Collection) User {
	filter := bson.M{"email": email}

	var result User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		panic(err)
	}

	return result
}
func loadUserById(id string, collection *mongo.Collection) User {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}

	var result User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		panic(err)
	}

	return result
}
