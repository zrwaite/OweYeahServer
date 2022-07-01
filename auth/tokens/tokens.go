package tokens

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenStruct struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func EncodeToken(username string) (tokenString string, success bool) {
	var err error
	jwtKeyString := os.Getenv("JWT_KEY")
	if jwtKeyString == "" {
		fmt.Println("JWT_KEY not found")
		return "", false
	}
	jwtKey := []byte(jwtKeyString)

	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	tokenBody := TokenStruct{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenBody)
	tokenString, err = token.SignedString(jwtKey)
	if err != nil {
		fmt.Println("Failed to create jwt - " + err.Error())
		return "", false
	}
	return tokenString, true
}

func ValidateToken(token string, username string) bool {
	var err error
	jwtKeyString := os.Getenv("JWT_KEY")
	if jwtKeyString == "" {
		fmt.Println("JWT_KEY not found")
		return false
	}
	jwtKey := []byte(jwtKeyString)
	tokenBody := TokenStruct{}

	tkn, err := jwt.ParseWithClaims(token, tokenBody, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return false
	}
	if !tkn.Valid {
		return false
	}
	if tokenBody.Username != username {
		return false
	}
	return true
}
