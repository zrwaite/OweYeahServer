package database

import "go.mongodb.org/mongo-driver/bson"

var UsersValidator = bson.M{
	"$jsonSchema": bson.M{
		"bsonType": "object",
		"required": []string{"id", "username", "hash", "display_name", "created_at", "invoice_ids", "payment_ids"},
		"properties": bson.M{
			"id":           bson.M{"bsonType": "string"},
			"username":     bson.M{"bsonType": "string"},
			"hash":         bson.M{"bsonType": "string"},
			"display_name": bson.M{"bsonType": "string"},
			"created_at":   bson.M{"bsonType": "string"},
			"invoice_ids": bson.M{
				"bsonType":    "array",
				"uniqueItems": true,
				"items":       bson.M{"bsonType": "string"},
			},
			"payment_ids": bson.M{
				"bsonType":    "array",
				"uniqueItems": true,
				"items":       bson.M{"bsonType": "string"},
			},
		},
	},
}
