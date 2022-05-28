package work

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"wxProjectDev/user"
)

// ServiceCreateReport 创建一个日报
func ServiceCreateReport(userID, teamIdStr, done, toDo, problem string, context *gin.Context) {
	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "teamID不为 int64",
				"repID": "",
			},
		)
		return
	}

	_, err = user.SelectOneMember(teamID, userID)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "userID 和 teamID无法对应",
				"repID": "",
			},
		)
		return
	}

	repID, err := createReport(userID, teamID, done, toDo, problem)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "error",
				"repID": strconv.FormatInt(repID, 10),
			},
		)
	} else {
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg":   "ok",
				"repID": strconv.FormatInt(repID, 10),
			},
		)
	}
}

// ServiceGetReport 获取单个日报的详细内容
func ServiceGetReport(repIDStr string, context *gin.Context) {
	// 查询repID是否为int64
	repID, err := strconv.ParseInt(repIDStr, 10, 64)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "repID不为 int64",
				"repID": "",
			},
		)
		return
	}

	// 查询该report
	var msg string
	var resStatus int
	rep, err := selectReport(repID)
	if err == sql.ErrNoRows {
		msg = "未找到该reportID"
		resStatus = http.StatusBadRequest
	} else if err != nil {
		msg = "数据库查询错误"
		resStatus = http.StatusServiceUnavailable
	} else {
		msg = "ok"
		resStatus = http.StatusOK
	}
	context.ShouldBindJSON(rep)
	context.JSON(
		resStatus,
		gin.H{
			"msg":    msg,
			"report": rep,
		},
	)
}

func ServiceGetTeamRep(teamIdStr, userID string, context *gin.Context) {
	teamID, err := strconv.ParseInt(teamIdStr, 10, 64)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "teamID不为 int64",
				"repID": "",
			},
		)
		return
	}

	_, err = user.SelectOneMember(teamID, userID)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "userID 和 teamID无法对应",
				"repID": "",
			},
		)
		return
	}

	reports, err := selectAllRep(teamID)
	if err != nil {
		context.JSON(
			http.StatusServiceUnavailable,
			gin.H{
				"msg":     err.Error(),
				"reports": nil,
			},
		)
		return
	}
	repInfos := make([]ReportInfo, 0, len(reports))
	for _, rep := range reports {
		var repInfoTemp = ReportInfo{
			rep.ReportID,
			rep.UserID,
			"",
			rep.RepDate,
		}
		repInfoTemp.UserName, _ = getUserName(rep.UserID)
		repInfos = append(repInfos, repInfoTemp)
	}
	context.ShouldBindJSON(repInfos)
	context.JSON(
		http.StatusOK,
		gin.H{
			"repInfos": repInfos,
		},
	)
}
