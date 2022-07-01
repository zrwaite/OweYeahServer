package model

import "github.com/zrwaite/OweMate/auth"

func (user *User) CreateHash(password string) error {
	hash, err := auth.HashPassword(password)
	user.Hash = &hash
	return err
}
