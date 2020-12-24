package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/render"  // or "gopkg.in/unrolled/render.v1"
)

func main() {
	r := render.New(render.Options{
		IndentJSON: true,
	})

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		r.JSON(c.Writer, http.StatusOK, map[string]string{"welcome": "This is rendered JSON!"})
	})

	router.Run(":3000")
}