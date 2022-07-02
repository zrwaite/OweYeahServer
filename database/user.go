package database

import (
	"context"
	"fmt"
	"time"

	"github.com/zrwaite/OweMate/auth/tokens"
	"github.com/zrwaite/OweMate/graph/model"
)

func GetUser(username string) (user *model.User, status int) {
	user = &model.User{}
	if username == "" {
		status = 400
		return
	}
	// opts := options.FindOne().SetProjection(projection)
	cursor := mongoDatabase.Collection("users").FindOne(context.TODO(), CreateUsernameFilter(username))
	if err := cursor.Decode(user); err != nil {
		if err.Error() == "mongo: no documents in result" {
			status = 404
		} else {
			fmt.Println("Failed to get user " + username + " ; " + err.Error())
			status = 400
		}
		return
	} else {
		status = 200
	}
	return
}

func CreateUser(userInput *model.UserInput) (errors []string, token string) {
	user := &model.User{}
	_, foundStatus := GetUser(userInput.Username)
	if foundStatus == 200 {
		errors = append(errors, "User already exists")
		return
	} else if foundStatus == 400 {
		errors = append(errors, "Something went wrong")
		return
	}

	err := user.CreateHash(userInput.Password)
	user.Username = userInput.Username
	user.CreatedAt = time.Now().Format("2006-01-02")
	user.InvoiceIds = []string{}
	user.PaymentIds = []string{}
	user.DisplayName = ""

	if err != nil {
		errors = append(errors, "Failed to create hash for user "+user.Username+" ; "+err.Error())
		return
	}
	newUser, insertErr := mongoDatabase.Collection("users").InsertOne(context.TODO(), user)
	if insertErr != nil {
		errors = append(errors, "Failed to create user "+user.Username+" ; "+insertErr.Error())
		return
	} else {
		fmt.Println(newUser)
		success := false
		token, success = tokens.EncodeToken(user.Username)
		if !success {
			errors = append(errors, "Failed to create token for user "+user.Username+" ; "+err.Error())
			return
		}
	}
	return
}
