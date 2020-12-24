package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
)

func main() {
	log := logrus.New()
	// hooks, config,...

	r := gin.New()
	r.Use(ginlogrus.Logger(log), gin.Recovery())

	// pingpong
	r.GET("/ping", func(c *gin.Context) {
		c.Data(200, "text/plain", []byte("pong"))
	})

	r.Run("127.0.0.1:8081")
}