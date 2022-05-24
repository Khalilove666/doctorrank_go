package helpers

import (
	"doctorrank_go/configs"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"strconv"
	"time"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Id        string
	jwt.StandardClaims
}

var SecretKey = configs.Env("SECRET_KEY")
var TokenMinutes, _ = strconv.ParseInt(configs.Env("TOKEN_MINUTES"), 10, 64)
var RefreshTokenMinutes, _ = strconv.ParseInt(configs.Env("REFRESH_TOKEN_MINUTES"), 10, 64)

func GenerateToken(email string, firstName string, lastName string, uid string) (signedToken string, err error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Id:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(TokenMinutes)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SecretKey))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, err
}
func GenerateRefreshToken(id string) (signedRefreshToken string, err error) {
	refreshClaims := &SignedDetails{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(RefreshTokenMinutes)).Unix(),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SecretKey))

	if err != nil {
		log.Panic(err)
		return
	}

	return refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}

	return claims, msg
}
