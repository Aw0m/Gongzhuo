package services

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"net/http"
	"strconv"
	"wxProjectDev/user/daos"
	"wxProjectDev/user/models"
)

func ServiceCreateTeam(openid, teamName string, context *gin.Context) {
	if teamID, err := daos.CreateTeam(openid, teamName); err == nil {
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg":    "ok",
				"teamID": strconv.FormatInt(teamID, 10),
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
	if err := daos.UpdateTeam(teamID, teamName); err == nil {
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
	if team, err := daos.SelectTeam(teamID); err == nil {
		teamStr := models.TeamStr{
			TeamIdStr: strconv.FormatInt(team.TeamID, 10),
			TeamName:  team.TeamName,
			CreatorID: team.CreatorID,
		}
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg":  "ok",
				"team": teamStr,
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
	team, err1 := daos.SelectTeam(teamID)
	membersStr, err2 := daos.SelectTeamMembers(teamID)
	if err1 == nil && err2 == nil {
		teamStr := models.TeamStr{
			TeamIdStr: strconv.FormatInt(team.TeamID, 10),
			TeamName:  team.TeamName,
			CreatorID: team.CreatorID,
		}
		for i := range membersStr {
			userTemp, _ := daos.SelectUser(membersStr[i].UserID)
			membersStr[i].UserName = userTemp.UserName
		}
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg":     "ok",
				"team":    teamStr,
				"members": membersStr,
			},
		)
	} else {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":     "error",
				"team":    team,
				"members": membersStr,
			},
		)
	}
}

// ServiceGetTeamCode  获得指定team的验证码
func ServiceGetTeamCode(teamIdStr string, userID string, context *gin.Context) {
	// teamID是否合格
	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	// 指定member是否存在
	member, err := daos.SelectOneMember(teamID, userID)
	if err != nil || !member.Admin {
		context.Status(http.StatusBadRequest)
		return
	}

	// 是否已经存在code，不存在则需要先生成
	code, err := daos.GetTeamCode(teamID)
	if err == redis.Nil {
		code = strconv.Itoa(10000 + rand.Intn(9999))
		if err = daos.SetTeamCode(teamID, code); err != nil {
			context.Status(http.StatusServiceUnavailable)
			return
		}
	} else if err != nil {
		context.Status(http.StatusServiceUnavailable)
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"msg":      "ok",
			"teamCode": code,
		},
	)
}

// ServiceUpdateTeamCode  更新指定team的验证码
func ServiceUpdateTeamCode(teamIdStr string, userID string, context *gin.Context) {
	// teamID是否合格
	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	// 指定member是否存在
	member, err := daos.SelectOneMember(teamID, userID)
	if err != nil || !member.Admin {
		context.Status(http.StatusBadRequest)
		return
	}

	// 生成code
	code := strconv.Itoa(10000 + rand.Intn(9999))
	if err = daos.SetTeamCode(teamID, code); err != nil {
		context.Status(http.StatusServiceUnavailable)
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"msg":      "ok",
			"teamCode": code,
		},
	)
}

// ServiceJoinTeam 根据验证码加入指定的team
func ServiceJoinTeam(userID, userName, teamIdStr, teamCode string, context *gin.Context) {
	// teamID是否合格
	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	// 是否已经存在code，不存在则需要先生成
	code, err := daos.GetTeamCode(teamID)
	if err == redis.Nil {
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg": "wrong code",
			},
		)
		return
	} else if err != nil {
		context.Status(http.StatusServiceUnavailable)
		return
	}

	// 判断输出code是否等于实际code
	if code != teamCode {
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg": "wrong code",
			},
		)
		return
	}

	// 验证通过则create Member
	err = daos.CreateMember(teamID, userID, userName, false)
	if err != nil {
		context.Status(http.StatusServiceUnavailable)
	} else {
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg": "ok",
			},
		)
	}
}

func ServiceSelectAllTeams(userID string, context *gin.Context) {
	if members, err := daos.SelectUserMembers(userID); err != nil {
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg": err.Error(),
			},
		)
	} else {
		teams := make([]models.TeamStr, 0, len(members))
		for _, mem := range members {
			teamTemp, _ := daos.SelectTeam(mem.TeamID)
			teamStrTemp := models.TeamStr{
				TeamIdStr: strconv.FormatInt(teamTemp.TeamID, 10),
				TeamName:  teamTemp.TeamName,
				CreatorID: teamTemp.CreatorID,
			}
			teams = append(teams, teamStrTemp)
		}

		context.JSON(
			http.StatusOK,
			gin.H{
				"msg":   "ok",
				"teams": teams,
			},
		)
	}
}

func ServiceAddAdmin(userID, memberID, teamIdStr string, context *gin.Context) {
	var httpCode int
	var msg string
	if teamID, err := strconv.ParseInt(teamIdStr, 10, 64); err != nil {
		httpCode = http.StatusBadRequest
		msg = "teamID 格式错误！"
	} else if team, err := daos.SelectTeam(teamID); err != nil {
		httpCode = http.StatusBadRequest
		msg = "不存在该teamID"
	} else if team.CreatorID != userID {
		httpCode = http.StatusBadRequest
		msg = "操作者不是Creator"
	} else if _, err := daos.SelectOneMember(teamID, memberID); err != nil {
		httpCode = http.StatusBadRequest
		msg = "该member不是该team的成员"
	} else if err := daos.SetAdmin(memberID, teamID); err != nil {
		httpCode = http.StatusServiceUnavailable
		msg = "无法设置该成员为admin"
	} else {
		httpCode = http.StatusOK
		msg = "ok"
	}
	context.JSON(
		httpCode,
		gin.H{
			"msg": msg,
		},
	)
}
