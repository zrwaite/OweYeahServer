package database

import (
	"context"
	"fmt"
	"time"

	"github.com/zrwaite/OweMate/auth"
	"github.com/zrwaite/OweMate/auth/tokens"
	"github.com/zrwaite/OweMate/graph/model"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUser(username string) (user *model.User, status int) {
	user = &model.User{}
	if username == "" {
		status = 400
		return
	}
	// opts := options.FindOne().SetProjection(projection)
	cursor := mongoDatabase.Collection("users").FindOne(context.TODO(), CreateUsernameFilter(username))
	// fmt.Printf("%+v\n", cursor)
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

func User(ctx context.Context, username string) *model.UserResult {
	user, status := GetUser(username)
	var userResult *model.UserResult
	if status == 200 {
		userResult = &model.UserResult{
			User:    user,
			Success: true,
		}
	} else if status == 404 {
		userResult = &model.UserResult{
			Success: false,
			Errors:  []string{"User not found"},
		}
	} else if status == 400 {
		userResult = &model.UserResult{
			Success: false,
			Errors:  []string{"Something went wrong"},
		}
	}
	return userResult
}

func CreateUser(ctx context.Context, input model.UserInput) (userAuthResult *model.UserAuthResult) {
	userAuthResult = &model.UserAuthResult{}
	user := &model.User{}
	_, foundStatus := GetUser(input.Username)
	if foundStatus == 200 {
		userAuthResult.Errors = append(userAuthResult.Errors, "User already exists")
		return
	} else if foundStatus == 400 {
		userAuthResult.Errors = append(userAuthResult.Errors, "Something went wrong")
		return
	}

	err := user.CreateHash(input.Password)
	user.Username = input.Username
	user.CreatedAt = time.Now().Format("2006-01-02")
	user.InvoiceIds = []string{}
	user.PaymentIds = []string{}
	user.DisplayName = ""

	if err != nil {
		userAuthResult.Errors = append(userAuthResult.Errors, "Failed to create hash for user "+user.Username+" ; "+err.Error())
		return
	}
	newUser, insertErr := mongoDatabase.Collection("users").InsertOne(context.TODO(), user)
	if insertErr != nil {
		userAuthResult.Errors = append(userAuthResult.Errors, "Failed to create user "+user.Username+" ; "+insertErr.Error())
		return
	} else {
		user.ID = newUser.InsertedID.(string)
		token, success := tokens.EncodeToken(user.Username)
		if !success {
			userAuthResult.Errors = append(userAuthResult.Errors, "Failed to create token for user "+user.Username+" ; "+err.Error())
			return
		}
		return &model.UserAuthResult{
			Success: true,
			Token:   token,
			User:    user,
		}
	}
}

func Login(ctx context.Context, input model.UserInput) *model.UserAuthResult {
	user, status := GetUser(input.Username)
	errors := []string{}
	if status == 404 {
		errors = append(errors, "User not found")
	} else if status == 400 {
		errors = append(errors, "Something went wrong")
	} else {
		if !auth.CheckPasswordHash(input.Password, user.Hash) {
			errors = append(errors, "Incorrect password")
		} else {
			token, tokenSuccess := tokens.EncodeToken(input.Username)
			if tokenSuccess {
				return &model.UserAuthResult{
					Success: true,
					Token:   token,
					User:    user,
				}
			} else {
				errors = append(errors, "Failed to create token")
			}
		}
	}
	return &model.UserAuthResult{
		Success: false,
		Errors:  errors,
	}
}

func DeleteUser(ctx context.Context, username string) *model.Result {
	result := &model.Result{}
	_, status := GetUser(username)
	if status == 404 {
		result.Errors = append(result.Errors, "User not found")
		return result
	} else if status == 400 {
		result.Errors = append(result.Errors, "Something went wrong")
		return result
	} else if status == 200 {
		_, err := mongoDatabase.Collection("users").DeleteOne(context.TODO(), CreateUsernameFilter(username))
		if err != nil {
			result.Errors = append(result.Errors, "Failed to delete user "+username+" ; "+err.Error())
			return result
		}
	}
	return &model.Result{
		Success: true,
	}
}

func UpdateUser(user *model.User) bool {
	update := bson.D{{"$set", user}}
	_, err := mongoDatabase.Collection("users").UpdateOne(context.TODO(), CreateUsernameFilter(user.Username), update)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
