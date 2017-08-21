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
	Tasks             []Task        `json:"tasks"`
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

func (s *Session) GetTasks(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	rows, err := db.Query(query.SearchQuery("getTask"),
		s.SessionID)
	if err != nil {
		return err
	}
	for rows.Next() {
		var t Task
		err := rows.Scan(t.TaskID,
			t.TaskName,
			t.TaskComplexity,
			t.IsDaily,
			t.TaskSource,
			t.TargetDate,
			t.TaskPriority,
			t.IsNew,
			t.IsInProgress)
		if err != nil {
			return err
		}
		s.Tasks = append(s.Tasks, t)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
