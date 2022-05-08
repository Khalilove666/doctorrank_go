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
	"log"
	"net/http"
	"strconv"
	"time"
)

var commentCollection *mongo.Collection = configs.GetCollection(configs.DB, "comments")

func CreateOrUpdateComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var comment models.Comment
		defer cancel()

		userId := c.GetString("_id")

		if err := c.BindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		validationErr := validate.Struct(comment)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
			return
		}

		count, err := commentCollection.CountDocuments(ctx, bson.M{"user_id": userId})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		if count > 0 {
			var doctor bson.M
			if err := c.BindJSON(&comment); err != nil {
				c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
				return
			}
			updatedAt := time.Now().Unix()
			doctor["updated_at"] = updatedAt
			result, err := doctorCollection.UpdateOne(
				ctx,
				bson.M{"user_id": userId},
				bson.M{"$set": doctor},
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
				return
			}

			c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: result})
			return

		} else {
			var doctor models.Doctor
			if err := c.BindJSON(&doctor); err != nil {
				c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
				return
			}

			validationErr := validate.Struct(doctor)
			if validationErr != nil {
				c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
				return
			}

			doctor.Id = primitive.NewObjectID()
			doctor.UserId = userId
			doctor.CreatedAt = time.Now().Unix()
			doctor.UpdatedAt = time.Now().Unix()

			resultInsertionNumber, insertErr := doctorCollection.InsertOne(ctx, doctor)
			if insertErr != nil {
				msg := fmt.Sprintf("Error creating doctor item")
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: msg})
				return
			}

			c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: resultInsertionNumber})
		}
	}
}

func AllComments() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var comments []models.Comment
		defer cancel()

		queries := c.Request.URL.Query()
		skip, _ := strconv.ParseInt(queries.Get("skip"), 10, 64)
		limit, _ := strconv.ParseInt(queries.Get("limit"), 10, 64)
		doctorId := queries.Get("doctor_id")

		opts := options.FindOptions{Skip: &skip, Limit: &limit}
		filter := bson.M{"doctor_id": primitive.ObjectIDFromHex(doctorId)}

		cursor, err := commentCollection.Find(ctx, filter, &opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if err = cursor.All(ctx, &comments); err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: comments})
	}
}
