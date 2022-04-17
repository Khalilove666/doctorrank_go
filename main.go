package main

import (
	"doctorrank_go/configs"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	configs.ConnectDB()
	router.Run("localhost:8000")
}
