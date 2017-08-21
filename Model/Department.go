package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type Department struct {
	ClientID          int    `json:"clientID"`
	DepartmentList    string `json:"departmentList"` //list of departments, separate by |
	EntryBy           string `json:"entryBy"`
	ResultCode        int
	ResultDescription string
	Selected          int    `json:selected`
	DepartmentName    string `json:departmentName`
}

func (d *Department) SaveDepartment(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("cmsSaveDepartments"),
		d.ClientID, d.DepartmentList, d.EntryBy).Scan(
		&d.ResultCode,
		&d.ResultDescription)
}

func (d *Department) GetAllDepartments(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("cmsGetAllDepartments"),
		d.ClientID).Scan(
		&d.Selected,
		&d.DepartmentName)
}
func (d *Department) GetActiveDepartments(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("cmsGetActiveDepartments"),
		d.ClientID).Scan(
		&d.DepartmentName)
}
