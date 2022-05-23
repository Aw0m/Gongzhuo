package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

		// 登录
		v1.POST("/login", func(context *gin.Context) {
			code := context.PostForm("code")
			userName := context.PostForm("userName")
			log.Println("code is ", code)
			ServiceLogin(code, userName, context)
		})

		// 更新用户名
		v1.POST("/userName", func(context *gin.Context) {
			openid := context.GetHeader("userID")
			userName := context.PostForm("userName")
			ServiceUpdateUserName(openid, userName, context)
		})

		// 获得用户名
		v1.GET("/userName", func(context *gin.Context) {
			openid := context.GetHeader("userID")
			ServiceSelectUserName(openid, context)
		})
	}

	v2 := e.Group("/team")
	{
		v2.POST("/createTeam", func(context *gin.Context) {
			openid := context.GetHeader("userID")
			teamName := context.PostForm("teamName")
			ServiceCreateTeam(openid, teamName, context)
		})

		v2.POST("/updateTeam", func(context *gin.Context) {
			teamID := context.PostForm("teamID")
			teamName := context.PostForm("teamName")
			if teamID, err := strconv.ParseInt(teamID, 10, 64); err != nil {
				context.JSON(
					http.StatusNotFound,
					gin.H{
						"msg": "error",
					},
				)
			} else {
				ServiceUpdateTeam(teamID, teamName, context)
			}

		})

		v2.GET("/selectTeam", func(context *gin.Context) {
			teamID := context.Query("teamID")
			if teamID, err := strconv.ParseInt(teamID, 10, 64); err != nil {
				context.JSON(
					http.StatusNotFound,
					gin.H{
						"msg":  "error",
						"team": nil,
					},
				)
			} else {
				ServiceSelectTeam(teamID, context)
			}
		})

		v2.GET("/selectMembers", func(context *gin.Context) {
			teamID := context.Query("teamID")
			if teamID, err := strconv.ParseInt(teamID, 10, 64); err != nil {
				context.JSON(
					http.StatusNotFound,
					gin.H{
						"msg":     "error",
						"team":    nil,
						"members": nil,
					},
				)
			} else {
				ServiceSelectMember(teamID, context)
			}
		})

		v2.GET("/selectTeamCode", func(context *gin.Context) {
			teamID := context.Query("teamID")
			userID := context.GetHeader("userID")
			ServiceGetTeamCode(teamID, userID, context)
		})

		v2.POST("/updateTeamCode", func(context *gin.Context) {
			teamID := context.PostForm("teamID")
			userID := context.GetHeader("userID")
			ServiceUpdateTeamCode(teamID, userID, context)
		})

		v2.POST("/joinTeam", func(context *gin.Context) {
			userID := context.GetHeader("userID")
			userName := context.PostForm("userName")
			teamIdStr := context.PostForm("teamID")
			teamCode := context.PostForm("teamCode")
			ServiceJoinTeam(userID, userName, teamIdStr, teamCode, context)
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
