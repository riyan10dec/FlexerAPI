package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type User struct {
	UserID            int    `json:"-"`
	ClientID          string `json:"-"`
	SuperiorID        int    `json:"superiorID"`
	UserName          string `json:"userName"`
	UserLogin         string `json:"userLogin"`
	UserPassword      string `json:"UserPassword"`
	Role              string `json:"Role"`
	ActiveStart       string `json:"ActiveStart"`
	ActiveEnd         string `json:"ActiveEnd"`
	Email             string `json:"Email"`
	ResultCode        int
	ResultDescription string
	EntryUser         string `json:"EntryUser"`
}

//DoLogin : Login Func
func (u *User) GetUserID(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("getUserID"),
		u.UserLogin).Scan(&u.UserID)
}

//AddEmployee : AddEmployee Func
func (u *User) AddEmployee(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("cmsAddUser"),
		u.ClientID,
		u.UserName,
		u.Role,
		u.SuperiorID,
		u.Email,
		u.UserPassword,
		u.ActiveStart,
		u.ActiveEnd,
	).Scan(&u.ResultCode, &u.ResultDescription)
}
