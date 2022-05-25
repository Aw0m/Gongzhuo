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
func selectAllRep(teamID int64) ([]Report, error) {
	rows, err := public.DB.Query("SELECT reportID, userID, repDate FROM report WHERE teamID = ?", teamID)
	if err != nil {
		log.Println("查询指定teamID的report出现错误：", err.Error())
		return nil, err
	}
	defer rows.Close()

	var reports []Report
	var timeStrTemp string
	for rows.Next() {
		var rep Report
		if err := rows.Scan(&rep.ReportID, &rep.UserID, &timeStrTemp); err != nil {
			log.Fatal(err)
		}
		rep.RepDate, _ = time.Parse("2006-01-02 15:04:05", timeStrTemp)
		reports = append(reports, rep)
	}
	return reports, nil
}

func getUserName(userID string) (string, error) {
	var userName string
	row := public.DB.QueryRow("SELECT userName FROM user WHERE userID = ?", userID)
	if err := row.Scan(&userName); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no such userID: %s\n", userID)
			return userName, err
		}
		return userName, fmt.Errorf("selectUser %s: %v", userID, err)
	}
	return userName, nil
}
