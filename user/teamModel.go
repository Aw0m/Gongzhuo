package user

type Team struct {
	TeamID    int64  `json:"teamID"`
	TeamName  string `json:"teamName"`
	CreatorID string `json:"creatorID"`
}

type Member struct {
	TeamID   int64  `json:"teamID"`
	UserID   string `json:"userID"`
	UserName string `json:"userName"`
	Admin    bool   `json:"admin"`
}
