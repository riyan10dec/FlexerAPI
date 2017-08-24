package model

type Activity struct {
	ActivityName   string  `json:"activityName"`
	ActivityType   int     `json:"activityType"`
	Category       int     `json:"category"`
	Classification string  `json:"classification"`
	Utilization    float64 `json:"utilization"`
}
