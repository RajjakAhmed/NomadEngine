package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/RajjakAhmed/NomadEngine/internal/store"
)

func main() {
	err := store.InitDB()
	if err != nil {
		log.Fatalf("DB Init failed: %v", err)
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "Nomad Engine running 🚀",
		})
	})

	log.Println("Nomad Engine started on :8080")
	router.Run(":8080")
}