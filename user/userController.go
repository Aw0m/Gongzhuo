package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
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
					"error": "Welcome Service User!",
				},
			)
		})

		v1.POST("/login/:code", func(context *gin.Context) {
			code := context.Param("code")
			log.Println("code is ", code)
			ServiceLogin(code, context)
		})
	}

	v2 := e.Group("work")
	{
		v2.GET("/", func(context *gin.Context) {

		})
	}

	return e
}

func parseResponse(response *http.Response) (map[string]interface{}, error) {
	var result map[string]interface{}
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(body, &result)
	}

	return result, err
}
