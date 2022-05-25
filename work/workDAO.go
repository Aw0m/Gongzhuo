package work

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"wxProjectDev/public"
	"wxProjectDev/utils"
)

func createReport(userID string, teamID int64, done, toDo, problem string) (int64, error) {
	utils.SetMachineId(0)
	repID := utils.GetSnowflakeId()
	repDate := time.Now().Local()

	_, err := public.DB.Exec(
		"INSERT INTO report (reportID, userID, teamID, done, toDo, problem, repDate) VALUES (?, ?, ?, ?, ?, ?, ?)",
		repID, userID, teamID, done, toDo, problem, repDate)
	if err != nil {
		log.Println("创建report时出现问题")
		return -1, err
	}

	return repID, nil
}

func selectReport(repID int64) (Report, error) {
	var report Report
	var repDateStr string
	row := public.DB.QueryRow("SELECT * FROM report WHERE reportID = ?", repID)
	if err := row.Scan(&report.ReportID, &report.UserID, &report.TeamID, &report.Done, &report.ToDO, &report.Problem, &repDateStr); err != nil {
		log.Printf("error: %v\n", err)
		if err == sql.ErrNoRows {
			log.Printf("no such userID: %d\n", repID)
			return report, err
		}
		return report, fmt.Errorf("selectUser %d: %v", repID, err)
	}
	// layout必须是 "2006-01-02 15:04:05"
	report.RepDate, _ = time.Parse("2006-01-02 15:04:05", repDateStr)
	return report, nil
}
