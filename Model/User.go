package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type User struct {
	UserID            int    `json:"userID"`
	EmployeeID        string `json:"employeeID"`
	ClientID          int    `json:"clientID"`
	SuperiorID        int    `json:"superiorID"`
	SuperiorName      string `json:"superiorName"`
	UserName          string `json:"userName"`
	UserLogin         string `json:"userLogin"`
	UserPassword      string
	OldPassword       string
	NewPassword       string
	ActiveStart       string `json:"activeStart"`
	ActiveEnd         string `json:"activeEnd"`
	ActiveStatus      string `json:"activeStatus"`
	Email             string `json:"email"`
	ResultCode        int
	ResultDescription string
	EntryUser         int            `json:"entryUser"`
	ModifiedBy        int            `json:"modifiedBy"`
	PositionName      string         `json:"positionName"`
	DepartmentName    string         `json:"departmentName"`
	IPAddress         sql.NullString `json:"ipAddress"`
	LoginDate         string         `json:"loginDate"`
	SubsCount         int            `json:"subsCount"`
	LastActivity      sql.NullString `json:"lastActivity"`
	ReferenceUser     []User
	ActiveOnly        bool `json:"activeOnly"`
	Features          []Feature
	Activities        []Activity
	SubscriptionID    int     `json:"subscriptionID"`
	GMTDiff           float32 `json:"gmtDiff"`
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
		u.EmployeeID,
		u.UserName,
		u.PositionName,
		u.DepartmentName,
		u.SuperiorID,
		u.Email,
		u.UserPassword,
		u.EntryUser,
		u.GMTDiff,
	).Scan(&u.ResultCode, &u.ResultDescription)
}

//EditEmployee : EditEmployee Func
func (u *User) EditEmployee(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("cmsEditUser"),
		u.UserID,
		u.EmployeeID,
		u.UserName,
		u.PositionName,
		u.DepartmentName,
		u.SuperiorID,
		u.Email,
		u.UserPassword,
		u.ActiveStatus,
		u.ActiveEnd,
		u.ModifiedBy,
		u.GMTDiff,
	).Scan(&u.ResultCode, &u.ResultDescription)
}

func (u *User) GetActiveSubs(db *sql.DB) error {
	rows, err := db.Query(query.SearchQuery("cmsGetActiveSubs"),
		u.UserID)
	defer rows.Close()
	//log.Fatal(u.UserID)
	if err != nil {
		return err
	}
	for rows.Next() {
		var u2 User
		err := rows.Scan(&u2.UserID, &u.EmployeeID, &u2.UserName, &u2.PositionName, &u2.DepartmentName, &u2.IPAddress, &u2.LoginDate)
		if err != nil {
			return err
		}
		u.ReferenceUser = append(u.ReferenceUser, u2)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) EmployeeTreeFirstLevel(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	var q string
	if u.ActiveOnly == true {
		q = query.SearchQuery("cmsEmployeeTreeFirstLevelActive")
	} else {
		q = query.SearchQuery("cmsEmployeeTreeFirstLevelAll")
	}

	rows, err := db.Query(q,
		u.ClientID, u.ClientID,
	)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var u2 User
		err := rows.Scan(&u2.UserID, &u2.UserName, &u2.SubsCount, &u2.ActiveStatus)
		if err != nil {
			return err
		}
		u.ReferenceUser = append(u.ReferenceUser, u2)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) EmployeeTreeSubs(db *sql.DB) error {
	var q string
	if u.ActiveOnly == true {
		q = query.SearchQuery("cmsEmployeeTreeSubsActive")
	} else {
		q = query.SearchQuery("cmsEmployeeTreeSubsAll")
	}

	rows, err := db.Query(q,
		u.UserID,
	)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var u2 User
		err := rows.Scan(&u2.UserID, &u2.UserName, &u2.SubsCount, &u2.ActiveStatus)
		if err != nil {
			return err
		}
		u.ReferenceUser = append(u.ReferenceUser, u2)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) EmployeeTreeChangeSuperior(db *sql.DB) (sql.Result, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := db.Prepare(query.SearchQuery("cmsEmployeeTreeChangeSuperior"))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	//start, _ := time.Parse(time.RFC3339, t.StartTime)
	//finish, _ := time.Parse(time.RFC3339, t.EndTime)
	res, err := stmt.Exec(u.SuperiorID, u.UserID)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return res, err
}

func (u *User) GetEmployees(db *sql.DB) error {
	rows, err := db.Query(query.SearchQuery("cmsEmployeeGrid"),
		u.UserID,
	)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var u2 User
		err := rows.Scan(&u2.UserID,
			&u2.EmployeeID,
			&u2.UserName,
			&u2.PositionName,
			&u2.DepartmentName,
			&u2.ActiveStatus,
			&u2.LastActivity)
		if err != nil {
			return err
		}
		u.ReferenceUser = append(u.ReferenceUser, u2)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) EmailValidation(db *sql.DB) error {
	return db.QueryRow(query.SearchQuery("cmsEmailValidation"),
		u.ClientID,
		u.Email,
		u.UserID).Scan(&u.ResultDescription)
}

func (u *User) ChangePassword(db *sql.DB) error {
	return db.QueryRow(query.SearchQuery("cmsChangePassword"),
		u.UserID,
		u.OldPassword,
		u.NewPassword).Scan(&u.ResultCode, &u.ResultDescription)
}

func (u *User) GetFeatures(db *sql.DB) error {
	rows, err := db.Query(query.SearchQuery("cmsGetFeatures"),
		u.UserID, u.PositionName, u.SubscriptionID,
	)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var f Feature
		err := rows.Scan(&f.FeatureID, &f.FeatureName, &f.FeatureType, &f.FeatureDescription)
		if err != nil {
			return err
		}
		u.Features = append(u.Features, f)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetSubs(db *sql.DB) error {
	rows, err := db.Query(query.SearchQuery("cmsGetSubs"),
		u.UserID, u.GMTDiff, u.ActiveOnly,
	)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var u2 User
		err := rows.Scan(&u2.UserID,
			&u2.EmployeeID,
			&u2.UserName,
			&u2.PositionName,
			&u2.DepartmentName,
			&u2.ActiveStatus,
			&u2.LastActivity,
			&u2.IPAddress,
			&u2.SuperiorID,
			&u2.SuperiorName,
			&u2.Email,
			&u2.ActiveStart,
			&u2.ActiveEnd)
		if err != nil {
			return err
		}
		u.ReferenceUser = append(u.ReferenceUser, u2)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetAllActivities(db *sql.DB) error {
	rows, err := db.Query(query.SearchQuery("cmsGetActivities"),
		u.UserID,
	)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var a Activity
		err := rows.Scan(&a.ActivityName, &a.ActivityType, &a.Category, &a.Classification, &a.Utilization)
		if err != nil {
			return err
		}
		u.Activities = append(u.Activities, a)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
