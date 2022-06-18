package daos

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"wxProjectDev/public"
	"wxProjectDev/public/utils"
	"wxProjectDev/work/models"
)

func CreateReport(userID string, teamID int64, done, toDo, problem, repType string) (int64, error) {
	utils.SetMachineId(0)
	repID := utils.GetSnowflakeId()
	repDate := time.Now().Local()

	_, err := public.DB.Exec(
		"INSERT INTO report (reportID, userID, teamID, done, toDo, problem, repDate, repType) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		repID, userID, teamID, done, toDo, problem, repDate, repType)
	if err != nil {
		log.Println("创建report时出现问题")
		return -1, err
	}

	return repID, nil
}

func SelectReport(repID int64) (models.Report, error) {
	var report models.Report
	var repDateStr string
	row := public.DB.QueryRow("SELECT * FROM report WHERE reportID = ?", repID)
	if err := row.Scan(&report.ReportID, &report.UserID, &report.TeamID, &report.Done, &report.ToDO, &report.Problem, &repDateStr, &report.Type); err != nil {
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
func SelectAllRep(teamID int64) ([]models.Report, error) {
	rows, err := public.DB.Query("SELECT reportID, userID, repDate, repType FROM report WHERE teamID = ?", teamID)
	if err != nil {
		log.Println("查询指定teamID的report出现错误：", err.Error())
		return nil, err
	}
	defer rows.Close()

	var reports []models.Report
	var timeStrTemp string
	for rows.Next() {
		var rep models.Report
		if err := rows.Scan(&rep.ReportID, &rep.UserID, &timeStrTemp, &rep.Type); err != nil {
			log.Fatal(err)
		}
		rep.RepDate, _ = time.Parse("2006-01-02 15:04:05", timeStrTemp)
		reports = append(reports, rep)
	}
	return reports, nil
}

func GetUserName(userID string) (string, error) {
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

func UpdateRep(repID int64, done, toDo, problem string) error {
	repDate := time.Now().Local()
	_, err := public.DB.Exec("UPDATE report SET done = ?, toDo = ?, problem = ?, repDate = ? WHERE reportID = ?",
		done, toDo, problem, repDate, repID)
	if err != nil {
		log.Printf("更新rep错误!, 错误：%v", err.Error())
		return err
	}
	return nil
}
