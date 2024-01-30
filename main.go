package main

import (
	"google-orc/api"
	"log"

	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Find .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	r := gin.Default()
	r.GET("/", api.HealthCheck)
	r.POST("/ocr", api.Ocr)

	r.Run(":9002")
}
