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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"net/http"
	"strconv"
	"time"
)

var doctorCollection *mongo.Collection = configs.GetCollection(configs.DB, "doctors")

func CreateOrUpdateDoctor() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId, _ := primitive.ObjectIDFromHex(c.GetString("_id"))

		count, err := doctorCollection.CountDocuments(ctx, bson.M{"user_id": userId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		if count > 0 {
			var doctor bson.M
			if err := c.BindJSON(&doctor); err != nil {
				c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
				return
			}

			if hospitalId, exist := doctor["hospital_id"].(string); exist {
				doctor["hospital_id"], _ = primitive.ObjectIDFromHex(hospitalId)
			}
			if professionId, exist := doctor["profession_id"].(string); exist {
				doctor["profession_id"], _ = primitive.ObjectIDFromHex(professionId)
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

func AllDoctors() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var doctors []bson.M
		defer cancel()

		queries := c.Request.URL.Query()
		term := queries.Get("term")
		skip, _ := strconv.ParseInt(queries.Get("skip"), 10, 64)
		limit, _ := strconv.ParseInt(queries.Get("limit"), 10, 64)
		if limit <= 0 {
			limit = 12
		}

		pipeline := []bson.M{
			{
				"$lookup": bson.M{
					"from": "comments",
					"let":  bson.M{"id": "$_id"},
					"pipeline": []bson.M{
						{"$match": bson.M{"$expr": bson.M{"$eq": []string{"$doctor_id", "$$id"}}}},
						{"$group": bson.M{"_id": nil, "value": bson.M{"$avg": "$rate"}, "count": bson.M{"$sum": 1}}},
						{"$project": bson.M{"_id": 0}},
					},
					"as": "rating",
				},
			},
			{
				"$unwind": bson.M{"path": "$rating", "preserveNullAndEmptyArrays": true},
			},
			{
				"$project": bson.M{
					"rating": bson.M{"$ifNull": []interface{}{"$rating", bson.M{"value": 0, "count": 0}}},
				},
			},
			{
				"$group": bson.M{
					"_id":   nil,
					"value": bson.M{"$avg": "$rating.value"},
					"count": bson.M{"$sum": "$rating.count"},
				},
			},
			{
				"$project": bson.M{"_id": 0, "genAvg": "$value", "genCount": "$count"},
			},
			{
				"$lookup": bson.M{
					"from":         "doctors",
					"localField":   "null",
					"foreignField": "null",
					"as":           "doctor",
				},
			},
			{
				"$unwind": bson.M{"path": "$doctor", "preserveNullAndEmptyArrays": true},
			},
			{
				"$lookup": bson.M{
					"from": "comments",
					"let":  bson.M{"id": "$doctor._id"},
					"pipeline": []bson.M{
						{
							"$match": bson.M{"$expr": bson.M{"$eq": []string{"$doctor_id", "$$id"}}},
						},
						{
							"$group": bson.M{"_id": nil, "value": bson.M{"$avg": "$rate"}, "count": bson.M{"$sum": 1}},
						},
						{"$project": bson.M{"_id": 0}},
					},
					"as": "rating",
				},
			},
			{
				"$unwind": bson.M{"path": "$rating", "preserveNullAndEmptyArrays": true},
			},
			{
				"$project": bson.M{
					"doctor":   1,
					"genAvg":   1,
					"genCount": 1,
					"rating":   bson.M{"$ifNull": []interface{}{"$rating", bson.M{"value": 0, "count": 0}}},
				},
			},
			{
				"$project": bson.M{
					"genAvg":   1,
					"genCount": 1,
					"doctor":   1,
					"rating":   1,
					"rank": bson.M{
						"$divide": []bson.M{
							{
								"$sum": []bson.M{
									{"$multiply": []interface{}{"$rating.value", "$rating.count"}},
									{"$multiply": []interface{}{"$genAvg", "$genCount"}},
								},
							},
							{"$sum": []interface{}{"$genCount", "$rating.count"}},
						},
					},
				},
			},
			{
				"$replaceRoot": bson.M{
					"newRoot": bson.M{"$mergeObjects": []interface{}{bson.M{"rank": "$rank", "rating": "$rating"}, "$doctor"}},
				},
			},
			{"$lookup": bson.M{
				"from":         "professions",
				"localField":   "profession_id",
				"foreignField": "_id",
				"as":           "profession",
			}},
			{"$lookup": bson.M{
				"from":         "hospitals",
				"localField":   "hospital_id",
				"foreignField": "_id",
				"as":           "hospital",
			}},
			{"$unwind": "$profession"},
			{"$unwind": "$hospital"},
			{
				"$project": bson.M{
					"full_name":  bson.M{"$concat": []string{"$first_name", " ", "$last_name"}},
					"title":      1,
					"user_id":    1,
					"first_name": 1,
					"last_name":  1,
					"img":        1,
					"rank":       1,
					"rating":     1,
					"profession": 1,
					"hospital":   1,
				},
			},
			{"$match": bson.M{"full_name": bson.M{"$regex": primitive.Regex{Pattern: term, Options: "i"}}}},
			{"$sort": bson.M{"rank": -1}},
			{"$skip": skip},
			{"$limit": limit},
		}

		cursor, err := doctorCollection.Aggregate(ctx, pipeline)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if err = cursor.All(ctx, &doctors); err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: doctors})
	}
}

func DoctorById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var doctors []bson.M
		var rating []bson.M
		var result bson.M

		doctorId, _ := primitive.ObjectIDFromHex(c.Param("doctorId"))

		pipeline := []bson.M{
			{"$lookup": bson.M{
				"from":         "professions",
				"localField":   "profession_id",
				"foreignField": "_id",
				"as":           "profession",
			}},
			{"$lookup": bson.M{
				"from":         "hospitals",
				"localField":   "hospital_id",
				"foreignField": "_id",
				"as":           "hospital",
			}},
			{"$unwind": "$profession"},
			{"$unwind": "$hospital"},
			{"$project": bson.M{
				"full_name":       bson.M{"$concat": []string{"$first_name", " ", "$last_name"}},
				"title":           1,
				"user_id":         1,
				"first_name":      1,
				"last_name":       1,
				"img":             1,
				"about":           1,
				"experience":      1,
				"education":       1,
				"contact":         1,
				"created_at":      1,
				"updated_at":      1,
				"profession._id":  1,
				"profession.name": 1,
				"hospital._id":    1,
				"hospital.name":   1,
				"hospital.img":    1,
			}},
			{"$match": bson.M{"_id": doctorId}},
		}
		cursor, err := doctorCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if err = cursor.All(ctx, &doctors); err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		if len(doctors) > 0 {
			pipeline = []bson.M{
				{"$match": bson.M{"doctor_id": doctorId}},
				{"$group": bson.M{"_id": doctorId, "rate": bson.M{"$avg": "$rate"}, "reviews": bson.M{"$sum": 1}}},
			}
			cursor, err = commentCollection.Aggregate(ctx, pipeline)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
				return
			}
			if err = cursor.All(ctx, &rating); err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
				return
			}

			result = doctors[0]
			if len(rating) > 0 {
				result["rate"] = rating[0]["rate"].(float64)
				result["reviews"] = rating[0]["reviews"].(int32)
			} else {
				result["rate"] = -1
				result["reviews"] = 0
			}
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: result})
	}
}

func DoctorBySelf() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var doctors []bson.M
		var result bson.M

		userId, _ := primitive.ObjectIDFromHex(c.GetString("_id"))

		pipeline := []bson.M{
			{"$lookup": bson.M{
				"from":         "professions",
				"localField":   "profession_id",
				"foreignField": "_id",
				"as":           "profession",
			}},
			{"$lookup": bson.M{
				"from":         "hospitals",
				"localField":   "hospital_id",
				"foreignField": "_id",
				"as":           "hospital",
			}},
			{"$unwind": "$profession"},
			{"$unwind": "$hospital"},
			{"$project": bson.M{
				"title":           1,
				"user_id":         1,
				"first_name":      1,
				"last_name":       1,
				"img":             1,
				"about":           1,
				"experience":      1,
				"education":       1,
				"contact":         1,
				"created_at":      1,
				"updated_at":      1,
				"profession._id":  1,
				"profession.name": 1,
				"hospital._id":    1,
				"hospital.name":   1,
				"hospital.img":    1,
			}},
			{"$match": bson.M{"user_id": userId}},
		}
		cursor, err := doctorCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if err = cursor.All(ctx, &doctors); err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if len(doctors) > 0 {
			result = doctors[0]
		}
		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: result})
	}
}

func UploadDoctorAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var image dto.ImageDTO
		var doctor bson.M
		defer cancel()

		userId, _ := primitive.ObjectIDFromHex(c.GetString("_id"))

		opts := options.FindOne().SetProjection(bson.M{"_id": 1})
		if err := doctorCollection.FindOne(ctx, bson.M{"user_id": userId}, opts).Decode(doctor); err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		doctorId := doctor["_id"].(primitive.ObjectID)

		if err := c.Bind(&image); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		if validationErr := validate.Struct(image); validationErr != nil {
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
			doctorId.Hex(),
			helpers.Folders.Doctor,
			image.Coordinates.Top,
			image.Coordinates.Left,
			image.Coordinates.Width,
			image.Coordinates.Height,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		_, err = doctorCollection.UpdateOne(
			ctx,
			bson.M{"_id": doctorId},
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
