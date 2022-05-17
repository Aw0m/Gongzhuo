package user

import (
	"database/sql"
	"fmt"
	"log"
	"wxProjectDev/utils"
)

func createTeam(creatorID, teamName string) (int64, error) {
	utils.SetMachineId(0)
	teamID := utils.GetSnowflakeId()

	trans, _ := db.Begin()
	_, err := trans.Exec("INSERT INTO team (teamID, teamName, creatorID) VALUE (?, ?, ?)", teamID, teamName, creatorID)
	if err != nil {
		log.Println("create team，出现错误！")
		trans.Rollback()
		return -1, fmt.Errorf("add: %v", err)
	}
	creator, _ := selectUser(creatorID)
	err = createMember(teamID, creatorID, creator.userName, true)
	if err != nil {
		trans.Rollback()
		return -2, fmt.Errorf("add: %v", err)
	}

	trans.Commit()
	return teamID, nil
}

func updateTeam(teamID int64, teamName string) error {
	_, err := db.Exec("UPDATE team SET teamName = ? WHERE teamID = ?", teamName, teamID)
	if err != nil {
		log.Println("update team，出现错误！")
		return fmt.Errorf("add: %v", err)
	}
	return nil
}

func selectTeam(teamID int64) (Team, error) {
	var team Team
	row := db.QueryRow("SELECT * FROM team WHERE teamID = ?", teamID)
	if err := row.Scan(&team.TeamID, &team.TeamName, &team.CreatorID); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no such teamID: %d\n", teamID)
			return team, err
		}
		return team, fmt.Errorf("select team %d: %v", teamID, err)
	}
	return team, nil
}

func createMember(teamID int64, userID string, userName string, admin bool) error {
	_, err := db.Exec("INSERT INTO member (teamID, userID, userName, admin) VALUE (?, ?, ?, ?)", teamID, userID, userName, admin)
	if err != nil {
		log.Println("create member，出现错误！")
		return fmt.Errorf("add: %v", err)
	}
	return nil
}

func selectMember(teamID int64) ([]Member, error) {
	rows, err := db.Query("SELECT * FROM member WHERE teamID = ?", teamID)
	if err != nil {
		log.Println("select Meber 出现错误")
		return nil, fmt.Errorf("select: %v", err)
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var member Member
		if err := rows.Scan(&member.TeamID, &member.UserID, &member.UserName, &member.Admin); err != nil {
			log.Fatal(err)
		}
		members = append(members, member)
	}
	return members, nil
}
