package controllers

import (
	"context"
	"doctorrank_go/configs"
	"doctorrank_go/models"
	"doctorrank_go/responses"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"time"
)

var doctorCollection *mongo.Collection = configs.GetCollection(configs.DB, "doctors")

func CreateDoctor() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var doctor models.Doctor
		defer cancel()

		userId := c.GetString("_id")

		if err := c.BindJSON(&doctor); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		validationErr := validate.Struct(doctor)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
			return
		}

		doctor.Id = primitive.NewObjectID()
		doctor.UserId = userId
		doctor.CreatedAt = time.Now().Unix()
		doctor.UpdatedAt = time.Now().Unix()

		resultInsertionNumber, insertErr := doctorCollection.InsertOne(ctx, doctor)
		if insertErr != nil {
			msg := fmt.Sprintf("Error creating doctor item")
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: msg})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: resultInsertionNumber})
	}
}

func AllDoctors() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var doctors []models.CompactDoctor
		defer cancel()

		queries := c.Request.URL.Query()
		skip, _ := strconv.ParseInt(queries.Get("skip"), 10, 64)
		limit, _ := strconv.ParseInt(queries.Get("limit"), 10, 64)
		term := queries.Get("term")
		opts := options.FindOptions{Skip: &skip, Limit: &limit}
		filter := bson.D{{"first_name", term}}
		cursor, err := userCollection.Find(ctx, filter, &opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if err = cursor.All(ctx, &doctors); err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: doctors})
	}
}

func DoctorById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var doctor models.Doctor

		userId := c.Param("user_id")
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&doctor)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: "unknown user id"})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: doctor})
	}
}
