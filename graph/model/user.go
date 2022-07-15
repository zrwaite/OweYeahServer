package model

import "github.com/zrwaite/OweMate/auth"

type User struct {
	ID            string              `bson:"id"`
	Username      string              `bson:"username"`
	Hash          string              `bson:"hash"`
	DisplayName   string              `bson:"display_name"`
	CreatedAt     string              `bson:"created_at"`
	InvoiceIds    []string            `bson:"invoice_ids"`
	Invoices      []*InvoiceOrPayment `bson:"invoices"`
	PaymentIds    []string            `bson:"payment_ids"`
	Payments      []*InvoiceOrPayment `bson:"payments"`
	ConnectionIds []string            `bson:"connection_ids"`
	Connections   []*UserConnection   `bson:"connections"`
}

func (user *User) CreateHash(password string) error {
	hash, err := auth.HashPassword(password)
	user.Hash = hash
	return err
}
