package routes

import (
	"doctorrank_go/controllers"
	"doctorrank_go/middlewares"
	"github.com/gin-gonic/gin"
)

func ProfessionRoute(router *gin.Engine) {
	router.POST("/professions", middlewares.Authentication(), middlewares.RoleDoctor(), controllers.CreateProfession())
	router.GET("/professions", controllers.AllProfessions())
}
