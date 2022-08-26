package database

import (
	"context"
	"fmt"
	"time"

	"github.com/zrwaite/OweYeah/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetInvoice(id string) (invoice *model.InvoiceOrPayment, status int) {
	invoice = &model.InvoiceOrPayment{}
	if id == "" {
		status = 400
		return
	}
	filter, filterSuccess := CreateIdFilter(id)
	if !filterSuccess {
		status = 400
		return
	}
	status = Get("invoices", filter, invoice)
	return
}

func CreateInvoice(ctx context.Context, input model.InvoiceOrPaymentInput) *model.InvoiceResult {
	invoiceResult := &model.InvoiceResult{}
	user, userStatus := GetUser(input.CreatedByUsername)
	connection, connectionStatus := GetConnection(input.ConnectionID)
	if userStatus == 404 {
		invoiceResult.Errors = append(invoiceResult.Errors, "User not found")
	} else if userStatus == 400 {
		invoiceResult.Errors = append(invoiceResult.Errors, "Something went wrong getting user")
	}
	if connectionStatus == 404 {
		invoiceResult.Errors = append(invoiceResult.Errors, "Connection not found")
	} else if connectionStatus == 400 {
		invoiceResult.Errors = append(invoiceResult.Errors, "Something went wrong getting connection")
	}

	if len(invoiceResult.Errors) != 0 {
		return invoiceResult
	}
	userConnection, validConnection := ParseUserConnection(input.CreatedByUsername, connection)

	if !validConnection {
		invoiceResult.Errors = append(invoiceResult.Errors, "Something went wrong parsing the connection")
		return invoiceResult
	}

	invoice := &model.InvoiceOrPayment{
		CreatedByUsername: input.CreatedByUsername,
		ConnectionID:      input.ConnectionID,
		CreatedAt:         time.Now().Format("2006-01-02"),
		Amount:            input.Amount,
	}
	newInvoice, insertErr := mongoDatabase.Collection("invoices").InsertOne(context.TODO(), invoice)
	if insertErr != nil {
		invoiceResult.Errors = append(invoiceResult.Errors, "Failed to create invoice; "+insertErr.Error())
		return invoiceResult
	} else {
		invoice.ID = newInvoice.InsertedID.(primitive.ObjectID).Hex()
		if !UpdateInvoice(invoice) {
			invoiceResult.Errors = append(invoiceResult.Errors, "Failed to update invoice")
			return invoiceResult
		}
		if !AddInvoiceToUser(user, invoice) {
			invoiceResult.Errors = append(invoiceResult.Errors, "Failed to add invoice to user")
			return invoiceResult
		}
		if !AddInvoiceToUser(userConnection.Contact, invoice) {
			invoiceResult.Errors = append(invoiceResult.Errors, "Failed to add invoice to contact user")
			return invoiceResult
		}
		if !SettleInvoiceConnectionDebt(user, invoice, connection) {
			invoiceResult.Errors = append(invoiceResult.Errors, "Failed to settle invoice connection debt")
			return invoiceResult
		}
		if err := ResolveCycles(user, userConnection); err != nil {
			invoiceResult.Errors = append(invoiceResult.Errors, "Failed to resolve cycles")
			return invoiceResult
		}
		return &model.InvoiceResult{
			Success: true,
			Invoice: invoice,
		}
	}
}

func AddInvoiceToUser(user *model.User, invoice *model.InvoiceOrPayment) bool {
	if invoice.ID == "" {
		return false
	}
	user.InvoiceIds = append(user.InvoiceIds, invoice.ID)
	return UpdateUser(user)
}

func SettleInvoiceConnectionDebt(user *model.User, invoice *model.InvoiceOrPayment, connection *model.Connection) bool {
	if connection.Username1 == user.Username {
		connection.Debt -= invoice.Amount
	} else if connection.Username2 == user.Username {
		connection.Debt += invoice.Amount
	} else {
		return false
	}
	return UpdateConnection(connection)
}

func UpdateInvoice(invoice *model.InvoiceOrPayment) bool {
	update := bson.D{{Key: "$set", Value: invoice}}
	filter, filterSuccess := CreateIdFilter(invoice.ID)
	if !filterSuccess {
		return false
	}
	res, err := mongoDatabase.Collection("invoices").UpdateOne(context.TODO(), filter, update)
	// fmt.Printf("%+v\n", res)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if res.MatchedCount == 0 {
		return false
	}
	return true
}
