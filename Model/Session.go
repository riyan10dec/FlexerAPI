package model

type Session struct {
	SessionID    int    `json:"-"`
	UserID       int    `json:"userID"`
	Status       string `json:"status"`
	EntryDate    string `json:"entryDate"`
	ModifiedDate string `json:"modifiedDate"`
	ExpiredDate  string `json:"expiredDate"`
	LocationType string `json:"locationType"`
	IPAddress    string `json:"ipAddress"`
	City         string `json:"city"`
	Lat          string `json:"lat"`
	Long         string `json:"long"`
}
