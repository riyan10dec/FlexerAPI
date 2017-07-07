package model

import (
	query "Flexer/Query"
	"database/sql"
)

type User struct {
	UserID       int    `json:"-"`
	ClientID     string `json:"-"`
	SuperiorID   string `json:"-"`
	UserName     string `json:"userName"`
	UserLogin    string `json:"userLogin"`
	UserPassword string `json:"UserPassword"`
	Role         string `json:"Role"`
}

//DoLogin : Login Func
func (u *User) GetUserID(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("getUserID"),
		u.UserLogin).Scan(&u.UserID)
}
