package user

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
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
		log.Println(err)
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

func selectMembers(teamID int64) ([]Member, error) {
	rows, err := db.Query("SELECT * FROM member WHERE teamID = ?", teamID)
	if err != nil {
		log.Println("select member 出现错误", err.Error())
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

func selectOneMember(teamID int64, userID string) (Member, error) {
	var member Member
	row := db.QueryRow("SELECT * FROM member WHERE teamID = ? AND userID = ?", teamID, userID)
	if err := row.Scan(&member.TeamID, &member.UserID, &member.UserName, &member.Admin); err != nil {
		log.Println("select one member，出现错误！")
		return member, fmt.Errorf("add: %v", err)
	}
	return member, nil
}

func getTeamCode(teamID int64) (string, error) {
	//ctx := context.Background()
	val, err := rdb.Get(ctx, strconv.FormatInt(teamID, 10)).Result()
	if err == redis.Nil {
		log.Printf("teamID: %d 还没有验证码\n", teamID)
		return "", redis.Nil
	} else if err != nil {
		log.Printf("teamID: %d 尝试获取验证码失败\n", teamID)
		return "", err
	}
	return val, nil
}

func setTeamCode(teamID int64, code string) error {
	//ctx := context.Background()
	_, err := rdb.Set(ctx, strconv.FormatInt(teamID, 10), code, 0).Result()
	if err != nil {
		return fmt.Errorf("设置team验证码出错: %v", err)
	}
	return nil
}
