package database

import "go.mongodb.org/mongo-driver/bson"

func CreateUsernameFilter(username string) bson.D {
	return bson.D{{
		Key: "username",
		Value: bson.D{{
			Key:   "$eq",
			Value: username,
		}},
	}}
}
