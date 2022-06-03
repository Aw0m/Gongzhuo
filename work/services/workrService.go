package services

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	userDao "wxProjectDev/user/daos"
	workDao "wxProjectDev/work/daos"
	"wxProjectDev/work/models"
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

	_, err = userDao.SelectOneMember(teamID, userID)
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

	repID, err := workDao.CreateReport(userID, teamID, done, toDo, problem)
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
	rep, err := workDao.SelectReport(repID)
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

	_, err = userDao.SelectOneMember(teamID, userID)
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

	reports, err := workDao.SelectAllRep(teamID)
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
	repInfos := make([]models.ReportInfo, 0, len(reports))
	for _, rep := range reports {
		var repInfoTemp = models.ReportInfo{
			ReportID: rep.ReportID,
			UserID:   rep.UserID,
			RepDate:  rep.RepDate,
		}
		repInfoTemp.UserName, _ = workDao.GetUserName(rep.UserID)
		repInfos = append(repInfos, repInfoTemp)
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"repInfos": repInfos,
		},
	)
}
