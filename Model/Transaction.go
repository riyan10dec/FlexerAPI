package model

type Transaction struct {
	TransactionID     string `json:"transactionID"`
	Keystroke         int    `json:"keystroke"`
	Mouseclick        int    `json:"mouseclick"`
	ActivityName      string `json:"activityName"`
	ActivityType      string `json:"activityType"`
	StartDate         string `json:"startDate"`
	EndDate           string `json:"endDate"`
	ResultCode        int
	ResultDescription string
}
