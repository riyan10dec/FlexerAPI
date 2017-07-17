package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type User struct {
	UserID            int    `json:"userID"`
	ClientID          string `json:"clientID"`
	SuperiorID        int    `json:"superiorID"`
	UserName          string `json:"userName"`
	UserLogin         string `json:"userLogin"`
	UserPassword      string `json:"userPassword"`
	Role              string `json:"role"`
	ActiveStart       string `json:"activeStart"`
	ActiveEnd         string `json:"activeEnd"`
	Email             string `json:"email"`
	ResultCode        int
	ResultDescription string
	EntryUser         int `json:"entryUser"`
	ModifiedBy        int `json:"modifiedBy"`
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
		u.EntryUser,
	).Scan(&u.ResultCode, &u.ResultDescription)
}

//EditEmployee : EditEmployee Func
func (u *User) EditEmployee(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("cmsEditUser"),
		u.ClientID,
		u.UserName,
		u.Role,
		u.SuperiorID,
		u.Email,
		u.UserPassword,
		u.ActiveStart,
		u.ActiveEnd,
		u.ModifiedBy,
	).Scan(&u.ResultCode, &u.ResultDescription)
}
