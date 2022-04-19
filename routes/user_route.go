package routes

import (
	"doctorrank_go/controllers"
	"doctorrank_go/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.POST("/register", controllers.Register())
	router.GET("/login", controllers.Login())
	router.GET("/refresh", controllers.Refresh())
	router.PUT("/role", middlewares.Authentication(), controllers.ChangeUserRole())
}
