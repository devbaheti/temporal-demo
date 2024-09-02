package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"temporal-demo/workers"
	"temporal-demo/workflows"
)

func main() {
	go workers.StartTemporalWorker()

	r := gin.Default()

	r.POST("/workflows", workflows.CreateWorkflow)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
