package database

import (
	"context"
	"fmt"
	"time"

	"github.com/zrwaite/OweMate/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetConnection(id string) (connection *model.DatabaseConnection, status int) {
	connection = &model.DatabaseConnection{}
	if id == "" {
		status = 400
		return
	}
	filter, filterSuccess := CreateIdFilter(connection.ID)
	if !filterSuccess {
		status = 400
		return
	}
	cursor := mongoDatabase.Collection("connections").FindOne(context.TODO(), filter)
	if err := cursor.Decode(connection); err != nil {
		if err.Error() == "mongo: no documents in result" {
			status = 404
		} else {
			fmt.Println("Failed to get connection " + id + " ; " + err.Error())
			status = 400
		}
		return
	} else {
		status = 200
	}
	return
}

func ParseUserConnection(username string, connection *model.DatabaseConnection) (userConnection *model.Connection, success bool) {
	userConnection = &model.Connection{
		ID:        connection.ID,
		Debt:      connection.Debt,
		CreatedAt: connection.CreatedAt,
	}
	if connection.Username1 == username {
		userConnection.ContactUsername = connection.Username2
	} else if connection.Username2 == username {
		userConnection.ContactUsername = connection.Username1
	} else {
		return
	}
	contact, status := GetUser(userConnection.ContactUsername)
	if status == 404 || status == 400 {
		return
	}
	userConnection.Contact = contact
	success = true
	return
}

func CreateConnection(ctx context.Context, username1 string, username2 string) *model.ConnectionResult {
	connectionResult := &model.ConnectionResult{}
	user, userStatus := GetUser(username1)
	contactUser, contactUserStatus := GetUser(username2)
	if userStatus == 404 {
		connectionResult.Errors = append(connectionResult.Errors, "User not found")
	} else if userStatus == 400 {
		connectionResult.Errors = append(connectionResult.Errors, "Something went wrong getting user")
	}
	if contactUserStatus == 404 {
		connectionResult.Errors = append(connectionResult.Errors, "Contact user not found")
	} else if contactUserStatus == 400 {
		connectionResult.Errors = append(connectionResult.Errors, "Something went wrong getting contact user")
	}

	existingConnectionStatus := FindExistingConnection(username1, username2)
	if existingConnectionStatus == 200 {
		connectionResult.Errors = append(connectionResult.Errors, "Connection already exists")
	} else if existingConnectionStatus == 400 {
		connectionResult.Errors = append(connectionResult.Errors, "Something went wrong finding existing connection")
	}

	if len(connectionResult.Errors) != 0 {
		return connectionResult
	}

	connection := &model.DatabaseConnection{
		Username1: username1,
		Username2: username2,
		CreatedAt: time.Now().Format("2006-01-02"),
	}
	newConnection, insertErr := mongoDatabase.Collection("connections").InsertOne(context.TODO(), connection)

	if insertErr != nil {
		connectionResult.Errors = append(connectionResult.Errors, "Failed to create connection; "+insertErr.Error())
		return connectionResult
	} else {
		connection.ID = newConnection.InsertedID.(primitive.ObjectID).Hex()
		if !UpdateConnection(connection) {
			connectionResult.Errors = append(connectionResult.Errors, "Failed to update connection")
			return connectionResult
		}
		if !AddConnectionToUser(user, connection.ID) {
			connectionResult.Errors = append(connectionResult.Errors, "Failed to add connection to user")
			return connectionResult
		}
		if !AddConnectionToUser(contactUser, connection.ID) {
			connectionResult.Errors = append(connectionResult.Errors, "Failed to add connection to contact user")
			return connectionResult
		}
		return &model.ConnectionResult{
			Success: true,
			Connection: &model.Connection{
				ID:              connection.ID,
				Contact:         contactUser,
				ContactUsername: username2,
				CreatedAt:       time.Now().Format("2006-01-02"),
			},
		}
	}
}

func AddConnectionToUser(user *model.User, connectionId string) bool {
	user.ConnectionIds = append(user.ConnectionIds, connectionId)
	return UpdateUser(user)
}

func UpdateConnection(connection *model.DatabaseConnection) bool {
	update := bson.D{{"$set", connection}}
	filter, filterSuccess := CreateIdFilter(connection.ID)
	if !filterSuccess {
		return false
	}
	res, err := mongoDatabase.Collection("connections").UpdateOne(context.TODO(), filter, update)
	// fmt.Printf("%+v\n", res)

	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if res.MatchedCount == 0 {
		return false
	}
	return true
}

func FindExistingConnection(username1 string, username2 string) int {
	filter := bson.D{{"$or", bson.A{bson.D{{"username1", username1}, {"username2", username2}}, bson.D{{"username1", username2}, {"username2", username1}}}}}
	connection := &model.DatabaseConnection{}
	cursor := mongoDatabase.Collection("connections").FindOne(context.TODO(), filter)
	if err := cursor.Decode(connection); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return 404
		} else {
			fmt.Println("Failed to get connection ; " + err.Error())
			return 400
		}
	}
	return 200
}
