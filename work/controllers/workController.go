package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"wxProjectDev/work/services"
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

	e.GET("/api/work/test/hello/:name", func(context *gin.Context) {
		time.Sleep(time.Millisecond * 5)
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg": "hello " + context.Param("name"),
			},
		)
	})

	v1 := e.Group("/work")
	{
		v1.POST("/createReport", func(context *gin.Context) {
			userID := context.GetHeader("userID")
			teamIdStr := context.PostForm("teamID")
			done := context.PostForm("done")
			toDo := context.PostForm("toDo")
			problem := context.PostForm("problem")
			repType := context.PostForm("repType")
			services.ServiceCreateReport(userID, teamIdStr, done, toDo, problem, repType, context)
		})

		v1.POST("updateReport", func(context *gin.Context) {
			userID := context.GetHeader("userID")
			reportIdStr := context.PostForm("reportID")
			done := context.PostForm("done")
			toDo := context.PostForm("toDo")
			problem := context.PostForm("problem")
			services.UpdateReport(userID, reportIdStr, done, toDo, problem, context)
		})

		v1.GET("/getReport", func(context *gin.Context) {
			repIdStr := context.Query("repID")
			services.GetReport(repIdStr, context)
		})

		v1.GET("/getReportInfos", func(context *gin.Context) {
			userID := context.GetHeader("userID")
			teamIdStr := context.Query("teamID")
			services.ServiceGetTeamRep(teamIdStr, userID, context)
		})

	}

	return e
}
