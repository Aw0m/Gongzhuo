package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Controller() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())

	v1 := e.Group("/user")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(
				http.StatusOK,
				gin.H{
					"code":  http.StatusOK,
					"error": "Welcome server 01",
				},
			)
		})

	}

	v2 := e.Group("work")
	{
		v2.GET("/", func(context *gin.Context) {

		})
	}

	return e
}
