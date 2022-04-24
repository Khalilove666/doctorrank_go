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

var hospitalCollection *mongo.Collection = configs.GetCollection(configs.DB, "hospitals")

func CreateHospital() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var hospital models.Hospital
		defer cancel()

		if err := c.BindJSON(&hospital); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		validationErr := validate.Struct(hospital)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
			return
		}

		hospital.Id = primitive.NewObjectID()

		resultInsertionNumber, insertErr := hospitalCollection.InsertOne(ctx, hospital)
		if insertErr != nil {
			msg := fmt.Sprintf("Error creating hospital item")
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: msg})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: resultInsertionNumber})
	}
}

func AllHospitals() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var hospitals []models.Hospital
		defer cancel()

		queries := c.Request.URL.Query()
		skip, _ := strconv.ParseInt(queries.Get("skip"), 10, 64)
		limit, _ := strconv.ParseInt(queries.Get("limit"), 10, 64)
		term := queries.Get("term")
		opts := options.FindOptions{Skip: &skip, Limit: &limit}
		filter := bson.D{{"name", term}}
		cursor, err := userCollection.Find(ctx, filter, &opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if err = cursor.All(ctx, &hospitals); err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: hospitals})
	}
}
