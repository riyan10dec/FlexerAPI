package model

import (
	query "FlexerAPI/Query"
	"database/sql"
	"time"
)

type Task struct {
	TaskID            int            `json:"taskID"`
	UserID            int64          `json:"userID"`
	TaskName          string         `json:"taskName"`
	TaskComplexity    string         `json:"taskComplexity"`
	IsDaily           bool           `json:"isDaily"`
	TargetDate        time.Time      `json:"targetDate"`
	TaskSource        string         `json:"taskSource"`
	TaskStatus        string         `json:"taskStatus"`
	IsInProgress      bool           `json:"isInProgress"`
	IsNew             bool           `json:"isNew"`
	TaskPriority      int            `json:"taskPriority"`
	PeriodStart       sql.NullString `json:"periodStart"`
	PeriodEnd         sql.NullString `json:"periodEnd"`
	AssignmentDate    time.Time      `json:"assignmentDate"`
	CompletedDate     time.Time      `json:"completedDate"`
	Session           Session
	ResultCode        int
	ResultDescription string
	Tasks             []Task
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
func (t *Task) GetUserTasks(db *sql.DB) error {
	var periodStart interface{}
	var periodEnd interface{}
	if t.PeriodStart.Valid == false {
		setNilIf(&periodStart)
	} else {
		periodStart = t.PeriodStart.String
	}
	if t.PeriodEnd.Valid == false {
		setNilIf(&periodEnd)
	} else {
		periodEnd = t.PeriodEnd.String
	}

	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	rows, err := db.Query(query.SearchQuery("getUserTask"),
		t.UserID, periodStart, periodEnd, t.IsInProgress)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var t2 Task
		err := rows.Scan(t2.TaskID,
			t2.TaskName,
			t2.TaskComplexity,
			t2.IsDaily,
			t2.TaskSource,
			t2.AssignmentDate,
			t2.TargetDate,
			t2.TaskPriority,
			t2.TaskStatus,
			t2.CompletedDate)
		if err != nil {
			return err
		}
		t.Tasks = append(t.Tasks, t2)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
