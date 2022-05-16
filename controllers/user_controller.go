package controllers

import (
	"context"
	"doctorrank_go/configs"
	"doctorrank_go/helpers"
	"doctorrank_go/models"
	"doctorrank_go/responses"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: "this email already exists"})
			return
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: "this username already exists"})
			return
		}

		password := helpers.HashPassword(user.Password)
		user.Password = password

		user.CreatedAt = time.Now().Unix()
		user.UpdatedAt = time.Now().Unix()
		user.Id = primitive.NewObjectID()

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: resultInsertionNumber})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var loginCredentials models.LoginCredentials
		var foundUser models.User
		var loggedUser bson.M
		defer cancel()

		if err := c.BindJSON(&loginCredentials); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		filter := bson.M{"$or": bson.A{bson.M{"email": loginCredentials.Login}, bson.M{"username": loginCredentials.Login}}}
		err := userCollection.FindOne(ctx, filter).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: "username or password is incorrect"})
			return
		}

		passwordIsValid, msg := helpers.VerifyPassword(loginCredentials.Password, foundUser.Password)
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: msg})
			return
		}

		token, _ := helpers.GenerateToken(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Id.Hex())
		if loginCredentials.RememberMe {
			refreshToken, _ := helpers.GenerateRefreshToken(foundUser.Id.Hex())
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "refreshToken",
				Value:    refreshToken,
				Path:     "/",
				Domain:   configs.Env("CLIENT"),
				MaxAge:   60 * int(helpers.RefreshTokenMinutes),
				SameSite: http.SameSiteNoneMode,
				Secure:   true,
				HttpOnly: true,
			})
		}

		bsonBytes, _ := bson.Marshal(foundUser)
		bson.Unmarshal(bsonBytes, &loggedUser)
		loggedUser["token"] = token

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: loggedUser})
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refreshToken",
			Value:    "",
			Path:     "/",
			Domain:   configs.Env("CLIENT"),
			MaxAge:   -1,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			HttpOnly: true,
		})

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: "logged out"})
	}
}

func Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		defer cancel()

		cookie, err := c.Cookie("refreshToken")
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}
		userDetails, msg := helpers.ValidateToken(cookie)
		if msg != "" {
			log.Panic(msg)
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: msg})
			return
		}
		id, err := primitive.ObjectIDFromHex(userDetails.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		err = userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: "incorrect user id"})
			return
		}

		token, _ := helpers.GenerateToken(user.Email, user.FirstName, user.LastName, user.Id.Hex())

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: token})
	}
}

func ChangeUserRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		id := c.GetString("_id")
		objId, _ := primitive.ObjectIDFromHex(id)

		result, err := userCollection.UpdateOne(
			ctx,
			bson.M{"_id": objId},
			bson.D{
				{"$set", bson.D{{"role", "doctor"}}},
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: result})
	}
}
