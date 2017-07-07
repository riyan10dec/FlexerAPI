package model

import (
	query "FlexerAPI/Query"
	"database/sql"
	"time"
)

type Transaction struct {
	TransactionID   int    `json:"transactionID"`
	UserID          int    `json:"-"`
	Keystroke       int    `json:"keystroke"`
	Mouseclick      int    `json:"mouseclick"`
	ApplicationID   int    `json:"-"`
	URL             string `json:"url"`
	ApplicationName string `json:"applicationName"`
	Userlogin       string `json:"userlogin"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
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
	start, _ := time.Parse(time.RFC3339, t.StartTime)
	finish, _ := time.Parse(time.RFC3339, t.EndTime)
	res, err := stmt.Exec(t.UserID, t.ApplicationID, t.URL, t.Keystroke, t.Mouseclick, start, finish)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return res, err
}
