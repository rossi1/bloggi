package users

import (

	"net/http"
	"strings"
	"errors"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/dgrijalva/jwt-go"

)


const appSECRET = "Golang"

func ReturnDBInstance() *gorm.DB {
	db, err := gorm.Open("sqlite3", "blog.db")
	if err != nil {
		panic("Error occured")
	}
	return db
}

func CreateToken(authD User) (string, error) {
	claims := jwt.MapClaims{}
	claims["user"] = authD.Username
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] =  time.Now().Unix()	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(appSECRET))
}

func fromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Invalid Authorization header")
	}
	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

func ValidateJwtToken(r *http.Request) (string, error) {
	tokenString, err := fromAuthHeader(r)
	if err != nil {
		return "", err
	}

	token, tokenErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSECRET), nil
	})
	if tokenErr != nil {
		return "", tokenErr
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	
	if token.Valid && ok {
		username, isOk := claims["user"].(string)
		
		if !isOk {
			return "", errors.New("Unexpected Error occured")
		}
		
		return username, nil
	}
	return "", errors.New("Token not Valid")

}
