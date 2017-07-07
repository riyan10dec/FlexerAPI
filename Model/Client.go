package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type Client struct {
	ClientName string `json:"clientName"`
	ClientID   int    `json:"-"`
}

//GetClient : Get Client Name by ClientID Func
func (c *Client) GetClient(db *sql.DB) error {
	return db.QueryRow(query.SearchQuery("getClientQuery"),
		c.ClientID).Scan(&c.ClientName)
}
