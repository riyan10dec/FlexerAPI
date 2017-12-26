package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type Performance struct {
	UserID                    int            `json:"userID"`
	EmployeeID                string         `json:"employeeID"`
	UserName                  string         `json:"userName"`
	PositionName              string         `json:"positionName"`
	DepartmentName            string         `json:"departmentName"`
	WorkDays                  int            `json:"workDays"`
	SessionDuration           float32        `json:"sessionDuration"`
	SessionDurationDaily      float32        `json:"sessionDurationDaily"`
	ActivityDuration          float32        `json:"activityDuration"`
	ActivityDurationDaily     float32        `json:"activityDurationDaily"`
	ActivityName              string         `json:"activityName"`
	ActivityType              string         `json:"activityType"`
	ActivityCategory          string         `json:"activityCategory"`
	ActivityClassification    string         `json:"activityClassification"`
	ActivityHour              int            `json:"activityHour"`
	ActivityHour2             int            `json:"activityHour2"`
	AMPM                      string         `json:"ampm"`
	ActivityHourLabel         string         `json:"activityHourLabel"`
	ProductiveDuration        float32        `json:"productiveDuration"`
	ProductiveDurationDaily   float32        `json:"productiveDurationDaily"`
	UnproductiveDuration      float32        `json:"unproductiveDuration"`
	UnproductiveDurationDaily float32        `json:"unproductiveDurationDaily"`
	UnclassifiedDuration      float32        `json:"unclassifiedDuration"`
	UnclassifiedDurationDaily float32        `json:"unclassifiedDurationDaily"`
	Keystroke                 float32        `json:"keystroke"`
	KeystrokePerHour          float32        `json:"keystrokePerHour"`
	MaxKeystrokePerHour       float32        `json:"maxKeystrokePerHour"`
	MouseClick                float32        `json:"mouseClick"`
	MouseClickPerHour         float32        `json:"mouseClickPerHour"`
	MaxMouseClickPerHour      float32        `json:"maxMouseClickPerHour"`
	PeriodStart               sql.NullString `json:"periodStart"`
	PeriodEnd                 sql.NullString `json:"periodEnd"`
	SessionDate               string         `json:"sessionDate"`
	FirstLoginDate            string         `json:"firstLoginDate"`
	NumOfResult               int            `json:"numOfResult"`
	Performances              []Performance
}

func (p *Performance) GetUserPerformance(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	var periodStart interface{}
	var periodEnd interface{}
	if p.PeriodStart.Valid == false {
		setNilIf(&periodStart)
	} else {
		periodStart = p.PeriodStart.String
	}
	if p.PeriodEnd.Valid == false {
		setNilIf(&periodEnd)
	} else {
		periodEnd = p.PeriodEnd.String
	}
	return db.QueryRow(query.SearchQuery("cmsGetUserPerformance"),
		p.UserID, periodStart, periodEnd).Scan(
		&p.UserID,
		&p.EmployeeID,
		&p.UserName,
		&p.PositionName,
		&p.DepartmentName,
		&p.WorkDays,
		&p.SessionDuration,
		&p.SessionDurationDaily,
		&p.ActivityDuration,
		&p.ActivityDurationDaily,
		&p.ProductiveDuration,
		&p.ProductiveDurationDaily,
		&p.UnproductiveDuration,
		&p.UnproductiveDurationDaily,
		&p.UnclassifiedDuration,
		&p.UnclassifiedDurationDaily,
		&p.Keystroke,
		&p.KeystrokePerHour,
		&p.MaxKeystrokePerHour,
		&p.MouseClick,
		&p.MouseClickPerHour,
		&p.MaxMouseClickPerHour,
	)
}

func (p *Performance) GetUserDaily(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	var periodStart interface{}
	var periodEnd interface{}
	if p.PeriodStart.Valid == false {
		setNilIf(&periodStart)
	} else {
		periodStart = p.PeriodStart.String
	}
	if p.PeriodEnd.Valid == false {
		setNilIf(&periodEnd)
	} else {
		periodEnd = p.PeriodEnd.String
	}

	rows, err := db.Query(query.SearchQuery("cmsGetUserDaily"),
		p.UserID, periodStart, periodEnd, p.NumOfResult,
	)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var p2 Performance
		err := rows.Scan(
			&p2.UserID,
			&p2.EmployeeID,
			&p2.UserName,
			&p2.PositionName,
			&p2.DepartmentName,
			&p2.SessionDate,
			&p2.FirstLoginDate,
			&p2.SessionDuration,
			&p2.ActivityDuration,
			&p2.ProductiveDuration,
			&p2.UnproductiveDuration,
			&p2.UnclassifiedDuration,
		)
		if err != nil {
			return err
		}
		p.Performances = append(p.Performances, p2)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
func (p *Performance) GetUserDailyActivity(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	var periodStart interface{}
	var periodEnd interface{}
	if p.PeriodStart.Valid == false {
		setNilIf(&periodStart)
	} else {
		periodStart = p.PeriodStart.String
	}
	if p.PeriodEnd.Valid == false {
		setNilIf(&periodEnd)
	} else {
		periodEnd = p.PeriodEnd.String
	}

	rows, err := db.Query(query.SearchQuery("cmsGetUserDailyActivity"),
		p.UserID, periodStart, periodEnd,
	)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var p2 Performance
		err := rows.Scan(
			&p2.SessionDate,
			&p2.ActivityName,
			&p2.ActivityType,
			&p2.ActivityCategory,
			&p2.ActivityClassification,
			&p2.Keystroke,
			&p2.MouseClick,
			&p2.ActivityDuration,
		)
		if err != nil {
			return err
		}
		p.Performances = append(p.Performances, p2)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func (p *Performance) GetUserDailyTimeline(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)

	rows, err := db.Query(query.SearchQuery("cmsGetUserDailyTimeline"),
		p.UserID, p.SessionDate,
	)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var p2 Performance
		err := rows.Scan(
			&p2.ActivityHour,
			&p2.ActivityHour2,
			&p2.AMPM,
			&p2.ActivityHourLabel,
			&p2.ProductiveDuration,
			&p2.UnproductiveDuration,
			&p2.UnclassifiedDuration,
			&p2.ActivityCategory,
		)
		if err != nil {
			return err
		}
		p.Performances = append(p.Performances, p2)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
func setNilIf(v *interface{}) {
	*v = nil
}
