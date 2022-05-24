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
	v1 := e.Group("/work")
	{
		v1.POST("/createReport", func(context *gin.Context) {
			//TODO 需要先验证一下该user是否是属于该team的
			userID := context.GetHeader("userID")
			teamID := context.PostForm("teamID")
			done := context.PostForm("done")
			toDo := context.PostForm("toDo")
			problem := context.PostForm("problem")
			ServiceCreateReport(userID, teamID, done, toDo, problem, context)
		})
	}

	return e
}
