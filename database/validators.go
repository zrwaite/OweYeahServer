package database

import "go.mongodb.org/mongo-driver/bson"

var bsonIDArray = bson.M{
	"bsonType":    "array",
	"uniqueItems": true,
	"items":       bson.M{"bsonType": "string"},
}

var bsonString = bson.M{"bsonType": "string"}

var UsersValidator = bson.M{
	"$jsonSchema": bson.M{
		"bsonType": "object",
		"required": []string{"id", "username", "hash", "display_name", "created_at", "invoice_ids", "payment_ids", "connection_ids"},
		"properties": bson.M{
			"id":             bsonString,
			"username":       bsonString,
			"hash":           bsonString,
			"display_name":   bsonString,
			"created_at":     bsonString,
			"invoice_ids":    bsonIDArray,
			"payment_ids":    bsonIDArray,
			"connection_ids": bsonIDArray,
		},
	},
}

var ConnectionsValidator = bson.M{
	"$jsonSchema": bson.M{
		"bsonType": "object",
		"required": []string{"id", "username1", "username2", "created_at", "debt"},
		"properties": bson.M{
			"id":         bsonString,
			"username1":  bsonString,
			"username2":  bsonString,
			"created_at": bsonString,
			"debt":       bson.M{"bsonType": "double"},
		},
	},
}

var InvoicesAndPaymentsValidator = bson.M{
	"$jsonSchema": bson.M{
		"bsonType": "object",
		"required": []string{"id", "created_by_username", "connection_id", "created_at", "amount"},
		"properties": bson.M{
			"id":                  bsonString,
			"created_by_username": bsonString,
			"connection_id":       bsonString,
			"amount":              bson.M{"bsonType": "double"},
			"created_at":          bsonString,
		},
	},
}
