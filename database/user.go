package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zrwaite/OweMate/auth"
	"github.com/zrwaite/OweMate/auth/tokens"
	"github.com/zrwaite/OweMate/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUser(username string) (user *model.User, status int) {
	user = &model.User{}
	if username == "" {
		status = 400
		return
	}
	status = Get("users", CreateUsernameFilter(username), user)
	return
}

func GetFilteredUsers(ctx context.Context, partialUsername string) (usersResult *model.UsersResult) {
	usersResult = &model.UsersResult{}
	if partialUsername == "" {
		usersResult.Errors = append(usersResult.Errors, "No username provided")
		return
	}
	cursor, findErr := mongoDatabase.Collection("users").Find(context.TODO(), CreatePartialUsernameFilter(partialUsername))
	if findErr != nil {
		usersResult.Errors = append(usersResult.Errors, "Failed to find users ; "+findErr.Error())
		return
	}

	if err := cursor.All(context.TODO(), &usersResult.Users); err != nil {
		if err.Error() == "mongo: no documents in result" {
			usersResult.Success = true
		} else {
			usersResult.Errors = append(usersResult.Errors, "Failed to get users ; "+err.Error())
		}
		return
	}
	usersResult.Success = true
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
	user.ConnectionIds = []string{}
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
		user.ID = newUser.InsertedID.(primitive.ObjectID).Hex()
		UpdateUser(user)
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
	update := bson.D{{Key: "$set", Value: user}}
	_, err := mongoDatabase.Collection("users").UpdateOne(context.TODO(), CreateUsernameFilter(user.Username), update)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func UserInvoices(ctx context.Context, obj *model.User) ([]*model.InvoiceOrPayment, error) {
	var invoices []*model.InvoiceOrPayment
	for _, invoiceId := range obj.InvoiceIds {
		invoice, status := GetInvoice(invoiceId)
		if status == 200 {
			invoices = append(invoices, invoice)
		} else if status == 404 {
			return invoices, errors.New("Invoice " + invoiceId + " not found")
		} else {
			return invoices, errors.New("Something went wrong while getting invoice " + invoiceId)
		}
	}
	return invoices, nil
}

func UserPayments(ctx context.Context, obj *model.User) ([]*model.InvoiceOrPayment, error) {
	var payments []*model.InvoiceOrPayment
	for _, paymentId := range obj.PaymentIds {
		payment, status := GetPayment(paymentId)
		if status == 200 {
			payments = append(payments, payment)
		} else if status == 404 {
			return payments, errors.New("Payment " + paymentId + " not found")
		} else {
			return payments, errors.New("Something went wrong while getting payment " + paymentId)
		}
	}
	return payments, nil
}

func UserConnections(ctx context.Context, obj *model.User) ([]*model.UserConnection, error) {
	var connections []*model.UserConnection
	for _, connectionId := range obj.ConnectionIds {
		databaseConnection, status := GetConnection(connectionId)
		if status == 200 {
			connection, parseSuccess := ParseUserConnection(obj.Username, databaseConnection)
			if parseSuccess {
				connections = append(connections, connection)
			} else {
				return connections, errors.New("Failed to parse connection " + connectionId)
			}
		} else if status == 404 {
			return connections, errors.New("Connection " + connectionId + " not found")
		} else {
			return connections, errors.New("Something went wrong while getting connection " + connectionId)
		}
	}
	return connections, nil
}
