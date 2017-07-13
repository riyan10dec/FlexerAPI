package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type Transaction struct {
	TransactionID     string `json:"transactionID"`
	Keystroke         int    `json:"keystroke"`
	Mouseclick        int    `json:"mouseclick"`
	ActivityName      string `json:"activityName"`
	ActivityType      string `json:"activityType"`
	StartDate         string `json:"startDate"`
	EndDate           string `json:"endDate"`
	ResultCode        int
	ResultDescription string
}

//CreateTransaction : Insert new Transaction
func (t *Transaction) CreateTransaction(db *sql.DB) (sql.Result, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := db.Prepare(query.SearchQuery("createTransactionQuery"))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	//start, _ := time.Parse(time.RFC3339, t.StartTime)
	//finish, _ := time.Parse(time.RFC3339, t.EndTime)
	//res, err := stmt.Exec(t.UserID, t.ApplicationID, t.URL, t.Keystroke, t.Mouseclick, start, finish)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return nil, err
}
