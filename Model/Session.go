package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type Session struct {
	SessionID         int64  `json:"sessionID"`
	UserID            int    `json:"userID"`
	Status            string `json:"status"`
	EntryDate         string `json:"entryDate"`
	ModifiedDate      string `json:"modifiedDate"`
	ExpiredDate       string `json:"expiredDate"`
	LocationType      string `json:"locationType"`
	IPAddress         string `json:"ipAddress"`
	City              string `json:"city"`
	Lat               string `json:"lat"`
	Long              string `json:"long"`
	ResultCode        int
	ResultDescription string
	ServerDate        string
	ClientDate        string `json:"clientDate"`
	StartTime         string
	EndTime           sql.NullString
	Transactions      []Transaction `json:"transactions"`
}

//FrontCheckSession : Check if session is valid and get server time
func (s *Session) FrontCheckSession(db *sql.DB) error {
	return db.QueryRow(query.SearchQuery("frontCheckSession"),
		s.SessionID).Scan(
		&s.ResultCode,
		&s.ResultDescription,
		&s.ServerDate,
		&s.StartTime,
		&s.EndTime)
}
