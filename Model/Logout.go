package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type Logout struct {
	SessionID         int `json:"sessionID"`
	ResultCode        int
	ResultDescription string
}

//DoLogout : Logout Func
func (l *Logout) DoLogout(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("logoutQuery"),
		l.SessionID).Scan(&l.ResultCode, &l.ResultDescription)
}
