package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/AmirMohG/container-runner/services"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func AddHandler(c *gin.Context, amqpConn services.AMQPConnection, redisClient services.RedisClient) {

	val, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "when": "extracting request body"})
		return
	}

	cid := randStringBytes(16)

	ch, err := amqpConn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	err = ch.Publish("", "containers", false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Headers: amqp.Table{
				"reqType": "add",
				"cid":     cid,
			},
			Body: []byte(val),
		})
	if err != nil {
		panic(err)
	}

	err = redisClient.Set(cid, "requested")
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"success": "added to queue", "cid": cid})
}
func StatHandler(c *gin.Context, redisClient services.RedisClient) {
	cid := c.Query("cid")
	val, err := redisClient.Get(cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": val})
}

func DeleteHandler(c *gin.Context, amqpConn services.AMQPConnection, redisClient services.RedisClient) {

	cid := c.Query("cid")
	ch, err := amqpConn.Channel()
	ch.Publish("", "containers", false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Headers: amqp.Table{
				"reqType": "delete",
			},
			Body: []byte(cid),
		})
	err = redisClient.Set(cid, "requestedDelete")
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"success": "requested deletion", "cid": cid})
}
func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
