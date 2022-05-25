package work

import "time"

type Report struct {
	ReportID int64     `json:"reportID"`
	UserID   string    `json:"userID"`
	TeamID   int64     `json:"teamID"`
	Done     string    `json:"done"`
	ToDO     string    `json:"toDO"`
	Problem  string    `json:"problem"`
	RepDate  time.Time `json:"repDate"`
}
