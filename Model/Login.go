package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type Login struct {
	ClientID          int     `json:"-"`
	Username          string  `json:"username"`
	Password          string  `json:"password"`
	Email             string  `json:"email"`
	LocationType      string  `json:"locationType"`
	IPAddress         string  `json:"ipAddress"`
	City              string  `json:"city"`
	Lat               float32 `json:"lat"`
	Long              float32 `json:"long"`
	Session           Session
	ResultCode        int
	ResultDescription string
}

//DoLogin : Login Func
func (l *Login) DoLogin(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("loginQuery"),
		l.Email,
		l.Password,
		l.LocationType,
		l.IPAddress,
		l.City,
		l.Lat,
		l.Long).Scan(&l.ResultCode, &l.ResultDescription, &l.Session.SessionID)
}

//DoLoginCMS : CMS Login Func
func (l *Login) DoLoginCMS(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("loginCMSQuery"),
		l.Email, l.Password).Scan(&l.ClientID)
}
