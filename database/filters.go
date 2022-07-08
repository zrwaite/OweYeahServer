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

func CreateIdFilter(id string) bson.D {
	return bson.D{{
		Key: "_id",
		Value: bson.D{{
			Key:   "$eq",
			Value: id,
		}},
	}}
}
