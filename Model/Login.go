package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type Login struct {
	ClientID int    `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Session  Session
}

//DoLogin : Login Func
func (l *Login) DoLogin(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("loginQuery"),
		l.Email, l.Password).Scan(&l.Session.SessionID) //.Scan(&l.Username, &l.ClientID)
}

//DoLoginCMS : CMS Login Func
func (l *Login) DoLoginCMS(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("loginCMSQuery"),
		l.Email, l.Password).Scan(&l.ClientID)
}
