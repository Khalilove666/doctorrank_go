package main

import (
	"doctorrank_go/configs"
	"doctorrank_go/middlewares"
	"doctorrank_go/routes"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"net/http"
)

func main() {

	port := configs.Env("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.Default()
	configs.ConnectDB()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{configs.Env("CLIENT")},
		AllowedMethods:   []string{http.MethodHead, http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Origin", "Authorization", "Content-Type"},
	})
	router.Use(c)

	routes.UserRoute(router)
	routes.DoctorRoute(router)
	routes.CommentRoute(router)
	routes.HospitalRoute(router)
	routes.ProfessionRoute(router)

	router.Use(middlewares.Authentication())

	router.Run("localhost:" + port)
}
