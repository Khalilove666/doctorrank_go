package routes

import (
	"doctorrank_go/controllers"
	"doctorrank_go/middlewares"
	"github.com/gin-gonic/gin"
)

func DoctorRoute(router *gin.Engine) {
	router.POST("/doctors", middlewares.Authentication(), controllers.CreateDoctor())
	router.GET("/doctors", controllers.AllDoctors())
	router.GET("/doctors/:user_id", controllers.DoctorById())
}
