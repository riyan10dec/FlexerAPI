package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type Screenshot struct {
	SessionID         int64 `json:"sessionID"`
	ScreenshotID      sql.NullInt64
	ScreenshotDate    string         `json:"screenshotDate"`
	ActivityName      sql.NullString `json:"activityName"`
	ActivityType      sql.NullString `json:"activityType"`
	Filename          sql.NullString
	ResultCode        int
	ResultDescription sql.NullString
}

//GetScreenshotParam : GetScreenshotParam Func
func (s *Screenshot) GetScreenshotParam(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("getScreenshotParamQuery"),
		s.SessionID, s.ActivityName.String, s.ActivityType.String, s.ScreenshotDate).Scan(&s.ResultCode, &s.ResultDescription, &s.ScreenshotID, &s.Filename)
}
func (s *Screenshot) ReportScreenshotStatus(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	return db.QueryRow(query.SearchQuery("reportScreenshotStatusQuery"),
		s.ScreenshotID, s.ResultCode, s.ResultDescription, s.Filename.String).Scan(&s.ResultCode, &s.ResultDescription)
}
