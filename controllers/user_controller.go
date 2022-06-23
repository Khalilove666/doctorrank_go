package controllers

import (
	"context"
	"doctorrank_go/configs"
	"doctorrank_go/dto"
	"doctorrank_go/helpers"
	"doctorrank_go/models"
	"doctorrank_go/responses"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"time"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var register dto.RegisterDTO
		defer cancel()

		if err := c.BindJSON(&register); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		validationErr := validate.Struct(register)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": register.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: "this email already exists"})
			return
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"username": register.Username})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: "this username already exists"})
			return
		}

		user.FirstName = register.FirstName
		user.LastName = register.LastName
		user.Username = register.Username
		user.Email = register.Email
		user.Password = helpers.HashPassword(register.Password)
		user.Role = "user"
		user.CreatedAt = time.Now().Unix()
		user.UpdatedAt = time.Now().Unix()
		user.Id = primitive.NewObjectID()

		signedActivationToken, err := helpers.GenerateEmailToken(user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		if err = helpers.SendConfirmationMail(user.FirstName, user.Email, configs.CLIENT+"/activation?activationToken="+signedActivationToken); err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User could not be not created")
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: msg})
			return
		}

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: resultInsertionNumber})
	}
}

func ActivateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		queries := c.Request.URL.Query()
		activationToken := queries.Get("activationToken")
		claims, msg := helpers.ValidateToken(activationToken)
		if msg != "" {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: msg})
			return
		}

		filter := bson.M{"email": claims.Email}
		update := bson.M{"$set": bson.M{"email_confirmed": true}}
		updateResult, err := userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: updateResult})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var loginCredentials dto.LoginDTO
		var loginRes dto.LoginResDTO
		var foundUser models.User
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

		if foundUser.EmailConfirmed != true {
			c.JSON(http.StatusForbidden, responses.Response{Status: http.StatusForbidden, Message: "error", Data: "email not confirmed"})
			return
		}

		token, _ := helpers.GenerateToken(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Id.Hex())
		refreshToken, _ := helpers.GenerateRefreshToken(foundUser.Id.Hex())

		maxCookieAge := 0
		if loginCredentials.RememberMe {
			maxCookieAge = 60 * int(helpers.RefreshTokenMinutes)
		}
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refreshToken",
			Value:    refreshToken,
			Path:     "/",
			Domain:   configs.Env("DOMAIN"),
			MaxAge:   maxCookieAge,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			HttpOnly: true,
		})

		bsonBytes, _ := bson.Marshal(foundUser)
		bson.Unmarshal(bsonBytes, &loginRes)
		loginRes.Token = token

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: loginRes})
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
			Domain:   configs.Env("DOMAIN"),
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
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}
		userDetails, msg := helpers.ValidateToken(cookie)
		if msg != "" {
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

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user dto.UserUpdateDTO
		defer cancel()

		userId, _ := primitive.ObjectIDFromHex(c.GetString("_id"))

		if err := c.Bind(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		if validationErr := validate.Struct(user); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
			return
		}
		updatedAt := time.Now().Unix()
		update := bson.M{}
		update["updated_at"] = updatedAt
		if user.FirstName != nil {
			update["first_name"] = user.FirstName
		}
		if user.LastName != nil {
			update["last_name"] = user.LastName
		}

		updateResult, err := userCollection.UpdateOne(
			ctx,
			bson.M{"_id": userId},
			bson.M{"$set": update},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: updateResult})
	}
}

func PasswordResetEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var foundUser models.User
		defer cancel()

		queries := c.Request.URL.Query()
		login := queries.Get("login")

		filter := bson.M{"$or": bson.A{bson.M{"email": login}, bson.M{"username": login}}}
		err := userCollection.FindOne(ctx, filter).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: "user not found"})
			return
		}

		signedEmailToken, err := helpers.GenerateEmailToken(foundUser.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		if err = helpers.SendPasswordResetEmail(foundUser.Email, configs.CLIENT+"/reset-password?pswResetToken="+signedEmailToken); err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: "email sent"})
	}
}

func ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var body dto.PasswordResetDTO
		defer cancel()

		queries := c.Request.URL.Query()
		pswResetToken := queries.Get("pswResetToken")
		claims, msg := helpers.ValidateToken(pswResetToken)
		if msg != "" {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: msg})
			return
		}

		if err := c.Bind(&body); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		if validationErr := validate.Struct(body); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
			return
		}

		newPassword := helpers.HashPassword(body.NewPassword)
		updatedAt := time.Now().Unix()
		updateResult, err := userCollection.UpdateOne(
			ctx,
			bson.M{"email": claims.Email},
			bson.M{"$set": bson.M{"password": newPassword, "updated_at": updatedAt}},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: updateResult})
	}
}

func ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var body dto.PasswordDTO
		var user models.User
		defer cancel()

		userId, _ := primitive.ObjectIDFromHex(c.GetString("_id"))

		if err := c.Bind(&body); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		if validationErr := validate.Struct(body); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		passwordIsValid, msg := helpers.VerifyPassword(body.OldPassword, user.Password)
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: msg})
			return
		}

		newPassword := helpers.HashPassword(body.NewPassword)
		updatedAt := time.Now().Unix()
		updateResult, err := userCollection.UpdateOne(
			ctx,
			bson.M{"_id": userId},
			bson.M{"$set": bson.M{"password": newPassword, "updated_at": updatedAt}},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: updateResult})
	}
}

func UploadAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var image dto.ImageDTO
		defer cancel()

		userId := c.GetString("_id")
		userObjId, _ := primitive.ObjectIDFromHex(userId)

		if err := c.Bind(&image); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		validationErr := validate.Struct(image)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
			return
		}

		file, err := image.File.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		defer file.Close()

		buffer, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		fileName, err := helpers.ProcessAndSaveAvatar(
			buffer,
			userId,
			helpers.Folders.User,
			image.Coordinates.Top,
			image.Coordinates.Left,
			image.Coordinates.Width,
			image.Coordinates.Height,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		_, err = userCollection.UpdateOne(
			ctx,
			bson.M{"_id": userObjId},
			bson.D{
				{"$set", bson.D{{"img", fileName}}},
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: fileName})
	}
}
