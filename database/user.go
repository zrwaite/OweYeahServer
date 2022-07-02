package database

import (
	"context"
	"fmt"

	"github.com/zrwaite/OweMate/graph/model"
)

func GetUser(username string) (user *model.User, status int) {
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
