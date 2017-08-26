package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type Department struct {
	ClientID             int      `json:"clientID"`
	DepartmentList       []string `json:"departmentList"` //list of departments, separate by |
	DepartmentsSeparator string   `json:"departmentSeparator"`
	EntryBy              string   `json:"entryBy"`
	ResultCode           int
	ResultDescription    string
	Selected             int    `json:"selected"`
	DepartmentName       string `json:"departmentName"`
	EmployeeCount        int    `json:"employeeCount"`
}

func (d *Department) SaveDepartment(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("cmsSaveDepartments"),
		d.ClientID, d.DepartmentList, d.EntryBy).Scan(
		&d.ResultCode,
		&d.ResultDescription)
}

func (d *Department) GetAllDepartments(db *sql.DB) (error, []Department) {
	rows, err := db.Query(query.SearchQuery("cmsGetAllDepartments"),
		d.ClientID)
	if err != nil {
		return err, nil
	}
	var ds []Department

	for rows.Next() {
		var d Department
		err := rows.Scan(&d.Selected,
			&d.DepartmentName, &d.EmployeeCount)
		if err != nil {
			return err, nil
		}
		ds = append(ds, d)
	}
	err = rows.Err()
	if err != nil {
		return err, nil
	}
	return nil, ds
}
func (d *Department) GetActiveDepartments(db *sql.DB) (error, []Department) {
	rows, err := db.Query(query.SearchQuery("cmsGetActiveDepartments"),
		d.ClientID)
	if err != nil {
		return err, nil
	}
	var ds []Department

	for rows.Next() {
		var d Department
		err := rows.Scan(&d.DepartmentName)
		if err != nil {
			return err, nil
		}
		ds = append(ds, d)
	}
	err = rows.Err()
	if err != nil {
		return err, nil
	}
	return nil, ds
}
