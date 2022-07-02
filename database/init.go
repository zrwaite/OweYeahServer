package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitializeDatabase() {
	mongoDatabase.Collection("users").Drop(context.TODO())
	var options = new(options.CreateCollectionOptions)
	options.Validator = UsersValidator
	error := mongoDatabase.CreateCollection(context.TODO(), "users", options)
	if error != nil {
		log.Fatal("Failed to create users collection" + error.Error())
	}
}
