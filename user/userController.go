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
		v1.POST("/login/:code", func(context *gin.Context) {
			code := context.Param("code")
			userName := context.PostForm("userName")
			log.Println("code is ", code)
			ServiceLogin(code, userName, context)
		})

		// 更新用户名
		v1.POST("/userName/:openid", func(context *gin.Context) {
			openid := context.Param("openid")
			userName := context.PostForm("userName")
			ServiceUpdateUserName(openid, userName, context)
		})
		// 获得用户名
		v1.GET("/userName/:openid", func(context *gin.Context) {
			openid := context.Param("openid")
			ServiceSelectUserName(openid, context)
		})
	}

	v2 := e.Group("/team")
	{
		v2.POST("/createTeam/:openid", func(context *gin.Context) {
			openid := context.Param("openid")
			teamName := context.PostForm("teamName")
			ServiceCreateTeam(openid, teamName, context)
		})

		v2.POST("/updateTeam/:teamID", func(context *gin.Context) {
			teamID := context.Param("teamID")
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

		v2.GET("/selectTeam/:teamID", func(context *gin.Context) {
			teamID := context.Param("teamID")
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

		v2.GET("/selectMembers/:teamID", func(context *gin.Context) {
			teamID := context.Param("teamID")
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
