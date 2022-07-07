package database

import (
	"context"
	"time"

	"github.com/zrwaite/OweMate/graph/model"
)

func CreateInvoice(ctx context.Context, input model.InvoiceInput) *model.InvoiceResult {
	invoiceResult := &model.InvoiceResult{}
	_, createdByUserFound := GetUser(input.CreatedByUsername)
	_, invoicedToUserFound := GetUser(input.InvoicedToUsername)
	if createdByUserFound == 404 {
		invoiceResult.Errors = []string{"Created-By-User not found"}
	} else if createdByUserFound == 400 {
		invoiceResult.Errors = []string{"Failed to query for Created-By-User"}
	}
	if invoicedToUserFound == 404 {
		invoiceResult.Errors = append(invoiceResult.Errors, "Invoiced-To-User not found")
	} else if invoicedToUserFound == 400 {
		invoiceResult.Errors = append(invoiceResult.Errors, "Failed to query for Invoiced-To-User")
	}

	if len(invoiceResult.Errors) != 0 {
		return invoiceResult
	}

	invoice := &model.Invoice{
		CreatedByUsername:  input.CreatedByUsername,
		InvoicedToUsername: input.InvoicedToUsername,
		Amount:             input.Amount,
		CreatedAt:          time.Now().Format("2006-01-02"),
	}

	_, insertErr := mongoDatabase.Collection("invoices").InsertOne(context.TODO(), invoice)
	if insertErr != nil {
		invoiceResult.Errors = append(invoiceResult.Errors, "Failed to create invoice; "+insertErr.Error())
		return invoiceResult
	} else {
		return &model.InvoiceResult{
			Success: true,
			Invoice: invoice,
		}
	}
}
