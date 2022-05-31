package routes

import (
	"doctorrank_go/configs"
	"doctorrank_go/controllers"
	"doctorrank_go/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	path := configs.Env("FILESYSTEM_PATH")

	router.POST("/register", controllers.Register())
	router.POST("/login", controllers.Login())
	router.POST("/logout", controllers.Logout())
	router.GET("/refresh", controllers.Refresh())
	router.PUT("/role", middlewares.Authentication(), controllers.ChangeUserRole())
	router.PUT("/update", middlewares.Authentication(), controllers.UpdateUser())
	router.PUT("/password", middlewares.Authentication(), controllers.ChangePassword())
	router.PUT("/avatar", middlewares.Authentication(), controllers.UploadAvatar())
	router.Static("/user/avatar", path+"/user/avatar/")
	router.Static("/user/thumbnail", path+"/user/thumbnail/")
}
