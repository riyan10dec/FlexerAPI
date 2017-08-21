package model

import (
	query "FlexerAPI/Query"
	"database/sql"
	"time"
)

type Task struct {
	TaskID            int       `json:"taskID"`
	TaskName          string    `json:"taskName"`
	TaskComplexity    string    `json:"taskComplexity"`
	IsDaily           bool      `json:"isDaily"`
	TargetDate        time.Time `json:"targetDate"`
	TaskSource        string    `json:"taskSource"`
	TaskStatus        string    `json:"taskStatus"`
	IsInProgress      bool      `json:"isInProgress"`
	IsNew             bool      `json:"isNew"`
	TaskPriority      int       `json:"taskPriority"`
	Session           Session
	ResultCode        int
	ResultDescription string
}

//AddTask :
func (t *Task) AddTask(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("AddTask"),
		t.Session.SessionID,
		t.TaskID,
		t.TaskName,
		t.TaskComplexity,
		t.TaskStatus,
	).Scan(&t.ResultCode, &t.ResultDescription)
}
