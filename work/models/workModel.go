package models

import "time"

type Report struct {
	ReportID string    `json:"reportID"`
	UserID   string    `json:"userID"`
	TeamID   string    `json:"teamID"`
	Done     string    `json:"done"`
	ToDO     string    `json:"todo"`
	Problem  string    `json:"problem"`
	RepDate  time.Time `json:"repDate"`
}

type ReportInfo struct {
	ReportID string    `json:"reportID"`
	UserID   string    `json:"userID"`
	UserName string    `json:"userName"`
	RepDate  time.Time `json:"repDate"`
}
