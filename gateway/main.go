package main

import (
	"log"

	handlers "github.com/AmirMohG/container-runner/gateway/handlers"
	services "github.com/AmirMohG/container-runner/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := gin.Default()

	redis := services.NewRedisClient()
	amqpConn, err := services.NewAMQPConnection()
	if err != nil {
		panic(err)
	}
	defer amqpConn.Close()
	ch, err := amqpConn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	amqpConn.CreateQueue("containers")

	//routes
	r.POST("/jobs/add", func(g *gin.Context) {
		handlers.AddHandler(g, amqpConn, redis)
	})
	r.GET("/jobs/stat", func(g *gin.Context) {
		handlers.StatHandler(g, redis)
	})
	r.GET("/jobs/delete", func(g *gin.Context) {
		handlers.DeleteHandler(g, amqpConn, redis)
	})
	//

	if err := r.Run(); err != nil {
		panic(err)
	}
}
