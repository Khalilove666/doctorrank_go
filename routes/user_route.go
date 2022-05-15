package routes

import (
	"doctorrank_go/controllers"
	"doctorrank_go/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.POST("/register", controllers.Register())
	router.POST("/login", controllers.Login())
	router.POST("/logout", controllers.Logout())
	router.GET("/refresh", controllers.Refresh())
	router.PUT("/role", middlewares.Authentication(), controllers.ChangeUserRole())
}
