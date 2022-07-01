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

var hospitalCollection *mongo.Collection = configs.GetCollection(configs.DB, "hospitals")

func CreateHospital() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var body dto.HospitalDTO
		var hospital models.Hospital
		defer cancel()

		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		if validationErr := validate.Struct(body); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
			return
		}

		count, err := hospitalCollection.CountDocuments(ctx, bson.M{"name": body.Name})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: "this hospital name already exists"})
			return
		}
		hospital.Id = primitive.NewObjectID()
		hospital.Name = body.Name

		_, insertErr := hospitalCollection.InsertOne(ctx, hospital)
		if insertErr != nil {
			msg := fmt.Sprintf("Error creating hospital item")
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: msg})
			return
		}

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: hospital.Id})
	}
}

func AllHospitals() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var hospitals []models.Hospital
		var filter bson.M
		defer cancel()

		queries := c.Request.URL.Query()
		skip, _ := strconv.ParseInt(queries.Get("skip"), 10, 64)
		limit, _ := strconv.ParseInt(queries.Get("limit"), 10, 64)
		term := queries.Get("term")
		opts := options.FindOptions{Skip: &skip, Limit: &limit}
		if term == "" {
			filter = bson.M{}
		} else {
			filter = bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: term, Options: "i"}}}
		}
		cursor, err := hospitalCollection.Find(ctx, filter, &opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if err = cursor.All(ctx, &hospitals); err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: hospitals})
	}
}

func UploadHospitalAvatar() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var image dto.ImageDTO
		defer cancel()

		hospitalId, _ := primitive.ObjectIDFromHex(c.Param("hospitalId"))

		count, err := hospitalCollection.CountDocuments(ctx, bson.M{"_id": hospitalId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		if count < 1 {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: "unknown hospital id"})
			return
		}

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
			hospitalId.Hex(),
			helpers.Folders.Hospital,
			image.Coordinates.Top,
			image.Coordinates.Left,
			image.Coordinates.Width,
			image.Coordinates.Height,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		_, err = hospitalCollection.UpdateOne(
			ctx,
			bson.M{"_id": hospitalId},
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
