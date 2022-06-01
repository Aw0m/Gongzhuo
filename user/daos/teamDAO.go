package daos

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"wxProjectDev/public"
	"wxProjectDev/public/utils"
	"wxProjectDev/user/models"
)

func CreateTeam(creatorID, teamName string) (int64, error) {
	utils.SetMachineId(0)
	teamID := utils.GetSnowflakeId()

	trans, _ := public.DB.Begin()
	_, err := trans.Exec("INSERT INTO team (teamID, teamName, creatorID) VALUE (?, ?, ?)", teamID, teamName, creatorID)
	if err != nil {
		log.Println("create team，出现错误！")
		trans.Rollback()
		return -1, fmt.Errorf("add: %v", err)
	}
	creator, _ := SelectUser(creatorID)
	err = CreateMember(teamID, creatorID, creator.UserName, true)
	if err != nil {
		trans.Rollback()
		log.Println(err)
		return -2, fmt.Errorf("add: %v", err)
	}

	trans.Commit()
	return teamID, nil
}

func UpdateTeam(teamID int64, teamName string) error {
	_, err := public.DB.Exec("UPDATE team SET teamName = ? WHERE teamID = ?", teamName, teamID)
	if err != nil {
		log.Println("update team，出现错误！")
		return fmt.Errorf("add: %v", err)
	}
	return nil
}

func SelectTeam(teamID int64) (models.Team, error) {
	var team models.Team
	row := public.DB.QueryRow("SELECT * FROM team WHERE teamID = ?", teamID)
	if err := row.Scan(&team.TeamID, &team.TeamName, &team.CreatorID); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no such teamID: %d\n", teamID)
			return team, err
		}
		return team, fmt.Errorf("select team %d: %v", teamID, err)
	}
	return team, nil
}

func CreateMember(teamID int64, userID string, userName string, admin bool) error {
	_, err := public.DB.Exec("INSERT INTO member (teamID, userID, userName, admin) VALUE (?, ?, ?, ?)", teamID, userID, userName, admin)
	if err != nil {
		log.Println("create member，出现错误！")
		return fmt.Errorf("add: %v", err)
	}
	return nil
}

// SelectTeamMembers TODO 需要修改一下表结构，删除 userName。 使得selectTeamMembers依然返回正常的 []Member
func SelectTeamMembers(teamID int64) ([]models.MemberStr, error) {
	rows, err := public.DB.Query("SELECT teamID, userID, admin FROM member WHERE teamID = ?", teamID)
	if err != nil {
		log.Println("select member 出现错误", err.Error())
		return nil, fmt.Errorf("select: %v", err)
	}
	defer rows.Close()

	var members []models.MemberStr
	for rows.Next() {
		var member models.MemberStr
		if err := rows.Scan(&member.TeamID, &member.UserID, &member.Admin); err != nil {
			log.Fatal(err)
		}
		members = append(members, member)
	}
	return members, nil
}

func SelectOneMember(teamID int64, userID string) (models.Member, error) {
	var member models.Member
	row := public.DB.QueryRow("SELECT * FROM member WHERE teamID = ? AND userID = ?", teamID, userID)
	if err := row.Scan(&member.TeamID, &member.UserID, &member.UserName, &member.Admin); err != nil {
		log.Println("select one member，出现错误！")
		return member, fmt.Errorf("add: %v", err)
	}
	return member, nil
}

func GetTeamCode(teamID int64) (string, error) {
	//ctx := context.Background()
	val, err := public.RDS.Get(public.CTX, strconv.FormatInt(teamID, 10)).Result()
	if err == redis.Nil {
		log.Printf("teamID: %d 还没有验证码\n", teamID)
		return "", redis.Nil
	} else if err != nil {
		log.Printf("teamID: %d 尝试获取验证码失败\n", teamID)
		return "", err
	}
	return val, nil
}

func SetTeamCode(teamID int64, code string) error {
	//ctx := context.Background()
	_, err := public.RDS.Set(public.CTX, strconv.FormatInt(teamID, 10), code, 0).Result()
	if err != nil {
		return fmt.Errorf("设置team验证码出错: %v", err)
	}
	return nil
}

func SelectUserMembers(userID string) ([]models.Member, error) {
	rows, err := public.DB.Query("SELECT * FROM member WHERE userID = ?", userID)
	if err != nil {
		log.Println("select member 出现错误", err.Error())
		return nil, fmt.Errorf("select: %v", err)
	}
	defer rows.Close()

	var members []models.Member
	for rows.Next() {
		var member models.Member
		if err := rows.Scan(&member.TeamID, &member.UserID, &member.UserName, &member.Admin); err != nil {
			log.Fatal(err)
		}
		members = append(members, member)
	}
	return members, nil
}

func SetAdmin(userID string, teamID int64) error {
	_, err := public.DB.Exec("UPDATE member SET admin = true WHERE userID = ? AND teamID = ?", userID, teamID)
	return err
}
