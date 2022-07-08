package database

import (
	"context"
	"time"

	"github.com/zrwaite/OweMate/graph/model"
)

func CreateInvoice(ctx context.Context, input model.InvoiceOrPaymentInput) *model.InvoiceOrPaymentResult {
	invoiceResult := &model.InvoiceOrPaymentResult{}
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
	}
	newInvoice, insertErr := mongoDatabase.Collection("invoices").InsertOne(context.TODO(), invoice)
	if insertErr != nil {
		invoiceResult.Errors = append(invoiceResult.Errors, "Failed to create invoice; "+insertErr.Error())
		return invoiceResult
	} else {
		invoice.ID = newInvoice.InsertedID.(string)
		if !AddInvoiceToUser(user, invoice) {
			invoiceResult.Errors = append(invoiceResult.Errors, "Failed to add invoice to user")
			return invoiceResult
		}
		if !AddInvoiceToUser(userConnection.Contact, invoice) {
			invoiceResult.Errors = append(invoiceResult.Errors, "Failed to add invoice to contact user")
			return invoiceResult
		}
		return &model.InvoiceOrPaymentResult{
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

func SettleConnectionDebt(user *model.User, invoice *model.InvoiceOrPayment, connection *model.DatabaseConnection) bool {
	if connection.Username1 == user.Username {
		connection.Debt -= invoice.Amount
	} else if connection.Username2 == user.Username {
		connection.Debt += invoice.Amount
	} else {
		return false
	}
	return true
}
