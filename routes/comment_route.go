package routes

import (
	"doctorrank_go/controllers"
	"doctorrank_go/middlewares"
	"github.com/gin-gonic/gin"
)

func CommentRoute(router *gin.Engine) {
	router.PUT("/comments", middlewares.Authentication(), controllers.CreateOrUpdateComment())
	router.GET("/comments", controllers.AllComments())
	router.PUT("/comments/:comment_id/like", middlewares.Authentication(), controllers.LikeOrDislikeComment())
}
