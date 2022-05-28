package user

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"net/http"
	"strconv"
)

func ServiceCreateTeam(openid, teamName string, context *gin.Context) {
	if teamID, err := createTeam(openid, teamName); err == nil {
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
		teamStr := TeamStr{
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
	team, err1 := selectTeam(teamID)
	membersStr, err2 := selectTeamMembers(teamID)
	if err1 == nil && err2 == nil {
		teamStr := TeamStr{
			TeamIdStr: strconv.FormatInt(team.TeamID, 10),
			TeamName:  team.TeamName,
			CreatorID: team.CreatorID,
		}
		for i := range membersStr {
			userTemp, _ := selectUser(membersStr[i].UserID)
			membersStr[i].UserName = userTemp.userName
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
	member, err := SelectOneMember(teamID, userID)
	if err != nil || !member.Admin {
		context.Status(http.StatusBadRequest)
		return
	}

	// 是否已经存在code，不存在则需要先生成
	code, err := getTeamCode(teamID)
	if err == redis.Nil {
		code = strconv.Itoa(10000 + rand.Intn(9999))
		if err = setTeamCode(teamID, code); err != nil {
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
	member, err := SelectOneMember(teamID, userID)
	if err != nil || !member.Admin {
		context.Status(http.StatusBadRequest)
		return
	}

	// 生成code
	code := strconv.Itoa(10000 + rand.Intn(9999))
	if err = setTeamCode(teamID, code); err != nil {
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
	code, err := getTeamCode(teamID)
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
	err = createMember(teamID, userID, userName, false)
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
	if members, err := selectUserMembers(userID); err != nil {
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg": err.Error(),
			},
		)
	} else {
		teams := make([]TeamStr, 0, len(members))
		for _, mem := range members {
			teamTemp, _ := selectTeam(mem.TeamID)
			teamStrTemp := TeamStr{
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
