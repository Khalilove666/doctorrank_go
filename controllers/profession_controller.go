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

var professionCollection *mongo.Collection = configs.GetCollection(configs.DB, "professions")

func CreateProfession() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var profession models.Profession
		defer cancel()

		if err := c.BindJSON(&profession); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
			return
		}

		validationErr := validate.Struct(profession)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
			return
		}
		count, err := professionCollection.CountDocuments(ctx, bson.M{"name": profession.Name})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: "this profession name already exists"})
			return
		}

		profession.Id = primitive.NewObjectID()
		resultInsertionNumber, insertErr := professionCollection.InsertOne(ctx, profession)
		if insertErr != nil {
			msg := fmt.Sprintf("Error creating profession item")
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: msg})
			return
		}

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: resultInsertionNumber})
	}
}

func AllProfessions() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var professions []models.Profession
		var filter primitive.M
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
		cursor, err := professionCollection.Find(ctx, filter, &opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if err = cursor.All(ctx, &professions); err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: professions})
	}
}
