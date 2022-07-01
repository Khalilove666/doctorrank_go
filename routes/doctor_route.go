package routes

import (
	"doctorrank_go/configs"
	"doctorrank_go/controllers"
	"doctorrank_go/middlewares"
	"github.com/gin-gonic/gin"
)

func DoctorRoute(router *gin.Engine) {
	path := configs.Env("FILESYSTEM_PATH")

	router.PUT("/doctors/update", middlewares.Authentication(), controllers.UpdateDoctor())
	router.PUT("/doctors/update/experience", middlewares.Authentication(), controllers.UpdateDoctorExperience())
	router.PUT("/doctors/update/education", middlewares.Authentication(), controllers.UpdateDoctorEducation())
	//router.DELETE("/doctors/update/experience", middlewares.Authentication(), controllers.DeleteDoctorExperience())
	//router.DELETE("/doctors/update/education", middlewares.Authentication(), controllers.DeleteDoctorEducation())
	router.GET("/doctors", controllers.AllDoctors())
	router.GET("/doctors/:doctorId", controllers.DoctorById())
	router.GET("/doctors/self", middlewares.Authentication(), controllers.DoctorBySelf())
	router.PUT("/doctors/avatar", middlewares.Authentication(), middlewares.RoleDoctor(), controllers.UploadDoctorAvatar())
	router.Static("/doctor/avatar", path+"/doctor/avatar/")
	router.Static("/doctor/thumbnail", path+"/doctor/thumbnail/")

}
