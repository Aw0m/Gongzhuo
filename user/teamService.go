package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ServiceCreateTeam(openid, teamName string, context *gin.Context) {
	if teamID, err := createTeam(openid, teamName); err == nil {
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg":    "ok",
				"teamID": teamID,
			},
		)
	} else {
		context.JSON(
			http.StatusServiceUnavailable,
			gin.H{
				"msg":    "error",
				"teamID": -1,
			},
		)
	}
}

func ServiceUpdateTeam(teamID int64, teamName string, context *gin.Context) {
	if err := updateTeam(teamID, teamName); err == nil {
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg": "ok",
			},
		)
	} else {
		context.JSON(
			http.StatusServiceUnavailable,
			gin.H{
				"msg": "error",
			},
		)
	}
}

func ServiceSelectTeam(teamID int64, context *gin.Context) {
	if team, err := selectTeam(teamID); err == nil {
		context.ShouldBindJSON(&team)
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg":  "ok",
				"team": team,
			},
		)
	} else {
		context.JSON(
			http.StatusServiceUnavailable,
			gin.H{
				"msg":  "error",
				"team": nil,
			},
		)
	}
}

func ServiceSelectMember(teamID int64, context *gin.Context) {
	team, err1 := selectTeam(teamID)
	members, err2 := selectMember(teamID)
	if err1 == nil && err2 == nil {
		context.ShouldBindJSON(&team)
		context.ShouldBindJSON(&members)
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg":     "ok",
				"team":    team,
				"members": members,
			},
		)
	} else {
		context.JSON(
			http.StatusServiceUnavailable,
			gin.H{
				"msg":     "error",
				"team":    team,
				"members": members,
			},
		)
	}
}
