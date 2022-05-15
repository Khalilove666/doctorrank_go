package routes

import (
	"doctorrank_go/controllers"
	"doctorrank_go/middlewares"
	"github.com/gin-gonic/gin"
)

func DoctorRoute(router *gin.Engine) {
	router.PUT("/doctors", middlewares.Authentication(), controllers.CreateOrUpdateDoctor())
	router.GET("/doctors", controllers.AllDoctors())
	router.GET("/doctors/:doctorId", controllers.DoctorById())
}
