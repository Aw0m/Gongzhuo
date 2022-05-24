package work

import (
	"log"
	"time"
	"wxProjectDev/public"
	"wxProjectDev/utils"
)

func createReport(userID, teamID, done, toDo, problem string) (int64, error) {
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
