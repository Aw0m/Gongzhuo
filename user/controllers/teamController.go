package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"wxProjectDev/user/services"
)

func ConCreateTeam(context *gin.Context) {
	openid := context.GetHeader("userID")
	teamName := context.PostForm("teamName")
	services.ServiceCreateTeam(openid, teamName, context)
}

func ConUpdateTeam(context *gin.Context) {
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
		services.ServiceUpdateTeam(teamID, teamName, context)
	}
}

func ConSelectTeam(context *gin.Context) {
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
		services.ServiceSelectTeam(teamID, context)
	}
}

func ConSelectMembers(context *gin.Context) {
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
		services.ServiceSelectMember(teamID, context)
	}
}

func ConSelectTeamCode(context *gin.Context) {
	teamID := context.Query("teamID")
	userID := context.GetHeader("userID")
	services.ServiceGetTeamCode(teamID, userID, context)
}

func ConUpdateTeamCode(context *gin.Context) {
	teamID := context.PostForm("teamID")
	userID := context.GetHeader("userID")
	services.ServiceUpdateTeamCode(teamID, userID, context)
}

func ConJoinTeam(context *gin.Context) {
	userID := context.GetHeader("userID")
	userName := context.PostForm("userName")
	teamIdStr := context.PostForm("teamID")
	teamCode := context.PostForm("teamCode")
	services.ServiceJoinTeam(userID, userName, teamIdStr, teamCode, context)
}

func ConSelectAllTeams(context *gin.Context) {
	userID := context.GetHeader("userID")
	services.ServiceSelectAllTeams(userID, context)
}

func ConAddAdmin(context *gin.Context) {
	userID := context.GetHeader("userID")
	memberID := context.PostForm("memberID")
	teamID := context.PostForm("teamID")
	services.ServiceAddAdmin(userID, memberID, teamID, context)
}
