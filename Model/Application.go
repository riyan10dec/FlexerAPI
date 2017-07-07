package model

import (
	query "Flexer/Query"
	"database/sql"
	"log"
)

type Application struct {
	ApplicationName string `json:"applicationName"`
	ApplicationID   int    `json:"-"`
}

//GetApplicationName : Get Application Name by ApplicationID Func
func (a *Application) GetApplicationName(db *sql.DB) error {
	return db.QueryRow(query.SearchQuery("getApplicationName"),
		a.ApplicationID).Scan(&a.ApplicationName)
}

//GetApplicationID : Get Application ID by ApplicationName Func
func (a *Application) GetApplicationID(db *sql.DB) error {
	return db.QueryRow(query.SearchQuery("getApplicationID"),
		a.ApplicationName).Scan(&a.ApplicationID)
}

//CheckApplicationExist : Check if application name exists by applicationName. Return 1 if exists
func (a *Application) CheckApplicationExist(db *sql.DB) (int, error) {
	var isExist bool
	err := db.QueryRow(query.SearchQuery("checkApplicationName"),
		a.ApplicationName).Scan(&isExist)
	switch {
	case err == sql.ErrNoRows:
		return 0, err
	case err != nil:
		return 0, err
	default:
		return 1, nil
	}
}

//CreateTransaction : Insert new Transaction
func (a *Application) CreateApplication(db *sql.DB) (sql.Result, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare(query.SearchQuery("createApplicationQuery"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(a.ApplicationName)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return res, err
}
