package work

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Controller() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery(), gin.Logger())
	e.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 01",
			},
		)
	})

	return e
}
