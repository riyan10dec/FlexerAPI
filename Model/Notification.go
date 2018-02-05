package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type Notification struct {
	NotificationID      int    `json:"notificationID"`
	NotificationMessage string `json:"notificationMessage"`
	ClientID            int    `json:"clientID"`
	UserID              int64  `json:"userID"`
	PageURL             string `json:"pageURL"`
	Seen                int    `json:"seen"`
	Notifications       []Notification
}

//DoLogin : Login Func
func (n *Notification) GetNotification(db *sql.DB) error {
	//rows, err := db.Query(query.SearchQuery("loginQuery"), l.UserLogin, l.Password)
	rows, err := db.Query(query.SearchQuery("getNotificationQuery"),
		n.UserID)
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var n2 Notification
		err := rows.Scan(&n2.NotificationID,
			&n2.NotificationMessage,
			&n2.PageURL,
			&n2.Seen)
		if err != nil {
			return err
		}
		n.Notifications = append(n.Notifications, n2)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
