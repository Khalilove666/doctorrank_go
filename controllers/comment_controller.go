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
	"log"
	"net/http"
	"strconv"
	"time"
)

var commentCollection *mongo.Collection = configs.GetCollection(configs.DB, "comments")

func CreateOrUpdateComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		queries := c.Request.URL.Query()
		userId, _ := primitive.ObjectIDFromHex(c.GetString("_id"))
		doctorId, _ := primitive.ObjectIDFromHex(queries.Get("doctorId"))

		doctorCount, err := doctorCollection.CountDocuments(ctx, bson.M{"_id": doctorId})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if doctorCount <= 0 {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: "doctorId not found"})
			return
		}
		count, err := commentCollection.CountDocuments(ctx, bson.M{"user_id": userId, "doctor_id": doctorId})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		if count > 0 {
			var comment bson.M
			if err := c.BindJSON(&comment); err != nil {
				c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
				return
			}
			comment["updated_at"] = time.Now().Unix()
			result, err := commentCollection.UpdateOne(
				ctx,
				bson.M{"user_id": userId, "doctor_id": doctorId},
				bson.M{"$set": comment},
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
				return
			}

			c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: result})
			return

		} else {
			var comment models.Comment
			if err := c.BindJSON(&comment); err != nil {
				c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
				return
			}

			validationErr := validate.Struct(comment)
			if validationErr != nil {
				c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
				return
			}

			comment.Id = primitive.NewObjectID()
			comment.UserId = userId
			comment.DoctorId = doctorId
			comment.Likes = []models.Like{}
			comment.CreatedAt = time.Now().Unix()
			comment.UpdatedAt = time.Now().Unix()

			resultInsertionNumber, insertErr := commentCollection.InsertOne(ctx, comment)
			if insertErr != nil {
				msg := fmt.Sprintf("Error creating comment item")
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
		var comments []bson.M
		defer cancel()

		queries := c.Request.URL.Query()
		skip, _ := strconv.ParseInt(queries.Get("skip"), 10, 64)
		limit, _ := strconv.ParseInt(queries.Get("limit"), 10, 64)
		if limit <= 0 {
			limit = 12
		}
		doctorId, _ := primitive.ObjectIDFromHex(queries.Get("doctorId"))

		pipeline := []bson.M{
			{
				"$match": bson.M{"doctor_id": doctorId},
			},
			{
				"$lookup": bson.M{
					"from":         "users",
					"localField":   "user_id",
					"foreignField": "_id",
					"as":           "user",
				},
			},
			{"$unwind": "$user"},
			{
				"$project": bson.M{
					"_id":             1,
					"text":            1,
					"doctor_id":       1,
					"rate":            1,
					"likes":           1,
					"user._id":        1,
					"user.first_name": 1,
					"user.last_name":  1,
					"user.username":   1,
					"user.img":        1,
				},
			},
			{"$skip": skip},
			{"$limit": limit},
		}
		cursor, err := commentCollection.Aggregate(ctx, pipeline)

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

func LikeOrDislikeComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var reqBody bson.M
		var update bson.M = nil
		defer cancel()

		userId, _ := primitive.ObjectIDFromHex(c.GetString("_id"))
		commentId, _ := primitive.ObjectIDFromHex(c.Param("comment_id"))

		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		filter := bson.M{
			"_id":           commentId,
			"likes.user_id": userId,
		}

		count, err := commentCollection.CountDocuments(ctx, filter)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		likeStatusFloat, _ := reqBody["like_status"].(float64)
		likeStatus := int(likeStatusFloat)

		if count >= 1 {
			if likeStatus == 1 {
				update = bson.M{"$set": bson.M{"likes.$.status": true}}
			} else if likeStatus == -1 {
				update = bson.M{"$set": bson.M{"likes.$.status": false}}
			} else if likeStatus == 0 {
				filter = bson.M{"_id": commentId}
				update = bson.M{"$pull": bson.M{"likes": bson.M{"user_id": userId}}}
			}
		} else {
			filter = bson.M{"_id": commentId}
			if likeStatus == 1 {
				update = bson.M{"$push": bson.M{"likes": bson.M{"user_id": userId, "status": true}}}
			} else if likeStatus == -1 {
				update = bson.M{"$push": bson.M{"likes": bson.M{"user_id": userId, "status": false}}}
			}
		}

		result, err := commentCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: result})
	}
}
