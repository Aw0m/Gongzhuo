package user

type User struct {
	userID   string
	userName string
}

func (u *User) UserID() string {
	return u.userID
}

func (u *User) SetUserID(userID string) {
	u.userID = userID
}

func (u *User) UserName() string {
	return u.userName
}

func (u *User) SetUserName(userName string) {
	u.userName = userName
}
