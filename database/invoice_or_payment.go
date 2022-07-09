package database

import (
	"context"
	"errors"

	"github.com/zrwaite/OweMate/graph/model"
)

func InvoiceOrPaymentConnection(ctx context.Context, obj *model.InvoiceOrPayment) (*model.Connection, error) {
	connection, status := GetConnection(obj.ConnectionID)
	if status == 404 {
		return nil, errors.New("connection not found")
	} else if status == 400 {
		return nil, errors.New("something went wrong getting connection")
	}
	return connection, nil
}
