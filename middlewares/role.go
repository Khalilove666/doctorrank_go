package middlewares

import (
	"context"
	"doctorrank_go/configs"
	"doctorrank_go/responses"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func RoleDoctor() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		id := c.GetString("_id")
		objId, _ := primitive.ObjectIDFromHex(id)
		count, err := userCollection.CountDocuments(ctx, bson.M{"$and": bson.A{bson.M{"_id": objId}, bson.M{"role": "doctor"}}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
			return
		}
		if count < 1 {
			c.JSON(http.StatusForbidden, responses.Response{Status: http.StatusForbidden, Message: "error", Data: "only doctors can do this action"})
			return
		}

		c.Next()
	}
}
