package database

import (
	"context"
	"fmt"
	"time"

	"github.com/zrwaite/OweYeah/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPayment(id string) (payment *model.InvoiceOrPayment, status int) {
	payment = &model.InvoiceOrPayment{}
	if id == "" {
		status = 400
		return
	}
	filter, filterSuccess := CreateIdFilter(id)
	if !filterSuccess {
		status = 400
		return
	}
	status = Get("payments", filter, payment)
	return
}

func CreatePayment(ctx context.Context, input model.InvoiceOrPaymentInput) *model.PaymentResult {
	paymentResult := &model.PaymentResult{}
	user, userStatus := GetUser(input.CreatedByUsername)
	connection, connectionStatus := GetConnection(input.ConnectionID)
	if userStatus == 404 {
		paymentResult.Errors = append(paymentResult.Errors, "User not found")
	} else if userStatus == 400 {
		paymentResult.Errors = append(paymentResult.Errors, "Something went wrong getting user")
	}
	if connectionStatus == 404 {
		paymentResult.Errors = append(paymentResult.Errors, "Connection not found")
	} else if connectionStatus == 400 {
		paymentResult.Errors = append(paymentResult.Errors, "Something went wrong getting connection")
	}

	if len(paymentResult.Errors) != 0 {
		return paymentResult
	}
	userConnection, validConnection := ParseUserConnection(input.CreatedByUsername, connection)

	if !validConnection {
		paymentResult.Errors = append(paymentResult.Errors, "Something went wrong parsing the connection")
		return paymentResult
	}

	payment := &model.InvoiceOrPayment{
		CreatedByUsername: input.CreatedByUsername,
		ConnectionID:      input.ConnectionID,
		CreatedAt:         time.Now().Format("2006-01-02"),
		Amount:            input.Amount,
	}
	newPayment, insertErr := mongoDatabase.Collection("payments").InsertOne(context.TODO(), payment)
	if insertErr != nil {
		paymentResult.Errors = append(paymentResult.Errors, "Failed to create payment; "+insertErr.Error())
		return paymentResult
	} else {
		payment.ID = newPayment.InsertedID.(primitive.ObjectID).Hex()
		if !UpdatePayment(payment) {
			paymentResult.Errors = append(paymentResult.Errors, "Failed to update payment")
			return paymentResult
		}
		if !AddPaymentToUser(user, payment) {
			paymentResult.Errors = append(paymentResult.Errors, "Failed to add payment to user")
			return paymentResult
		}
		if !AddPaymentToUser(userConnection.Contact, payment) {
			paymentResult.Errors = append(paymentResult.Errors, "Failed to add payment to contact user")
			return paymentResult
		}
		if !SettlePaymentConnectionDebt(user, payment, connection) {
			paymentResult.Errors = append(paymentResult.Errors, "Failed to settle payment connection debt")
			return paymentResult
		}
		return &model.PaymentResult{
			Success: true,
			Payment: payment,
		}
	}
}

func AddPaymentToUser(user *model.User, payment *model.InvoiceOrPayment) bool {
	if payment.ID == "" {
		return false
	}
	user.PaymentIds = append(user.PaymentIds, payment.ID)
	return UpdateUser(user)
}

func SettlePaymentConnectionDebt(user *model.User, payment *model.InvoiceOrPayment, connection *model.Connection) bool {
	if connection.Username1 == user.Username {
		connection.Debt += payment.Amount
	} else if connection.Username2 == user.Username {
		connection.Debt -= payment.Amount
	} else {
		return false
	}
	return UpdateConnection(connection)
}

func UpdatePayment(payment *model.InvoiceOrPayment) bool {
	update := bson.D{{Key: "$set", Value: payment}}
	filter, filterSuccess := CreateIdFilter(payment.ID)
	if !filterSuccess {
		return false
	}
	res, err := mongoDatabase.Collection("payments").UpdateOne(context.TODO(), filter, update)
	// fmt.Printf("%+v\n", res)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if res.MatchedCount == 0 {
		return false
	}
	return true
}
