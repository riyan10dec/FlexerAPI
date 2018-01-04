package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type Activity struct {
	ActivityName      string  `json:"activityName"`
	ActivityType      int     `json:"activityType"`
	Category          int     `json:"category"`
	Classification    string  `json:"classification"`
	Utilization       float64 `json:"utilization"`
	UserID            int     `json:"userID"`
	GMTDiff           float32 `json:"gmtDiff"`
	ResultCode        int
	ResultDescription string
	Activities        []Activity
}

func (u *Activity) GetAllActivities(db *sql.DB) error {
	rows, err := db.Query(query.SearchQuery("cmsGetActivities"),
		u.UserID, u.GMTDiff,
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
func (a *Activity) SaveActivity(db *sql.DB) error {
	return db.QueryRow(query.SearchQuery("cmsSaveActivity"),
		a.UserID, a.ActivityName, a.ActivityType, a.Category, a.Classification).Scan(
		&a.ResultCode,
		&a.ResultDescription)
}
