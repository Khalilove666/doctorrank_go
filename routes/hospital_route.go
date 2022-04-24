package routes

import (
	"doctorrank_go/controllers"
	"doctorrank_go/middlewares"
	"github.com/gin-gonic/gin"
)

func HospitalRoute(router *gin.Engine) {
	router.POST("/hospitals", middlewares.Authentication(), middlewares.RoleDoctor(), controllers.CreateHospital())
	router.GET("/hospitals", controllers.AllHospitals())
}
