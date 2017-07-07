package model

import (
	query "Flexer/Query"
	"database/sql"
)

type Login struct {
	ClientID  int    `json:"-"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Userlogin string `json:"userlogin"`
}

//DoLogin : Login Func
func (l *Login) DoLogin(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("loginQuery"),
		l.Userlogin, l.Password).Scan(&l.Username, &l.ClientID)
}
