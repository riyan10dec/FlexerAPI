package model

import (
	query "FlexerAPI/Query"
	"database/sql"
)

type Client struct {
	ClientName         string `json:"clientName"`
	ClientID           int    `json:"-"`
	SubscriptionType   string `json:"subscriptionType"`
	SubscriptionStatus string `json:"subscriptionStatus"`
	SubscriptionStart  string `json:"subscriptionStart"`
	SubscriptionEnd    string `json:"subscriptionEnd"`
	GraceUntil         string `json:"graceUntil"`
	MaxUser            int    `json:"maxUser"`
	RegisteredMember   int    `json:"registeredMember"`
	User               User   `json:"user"`
}

//GetClient : Get Client Name by ClientID Func
func (c *Client) GetClient(db *sql.DB) error {
	return db.QueryRow(query.SearchQuery("getClientQuery"),
		c.ClientID).Scan(&c.ClientName)
}

//CheckSubscription : Checking Client Subscription Staus
func (c *Client) CheckSubscription(db *sql.DB) error {
	return db.QueryRow(query.SearchQuery("cmsCheckSubscription"),
		c.ClientID).Scan(&c.SubscriptionType,
		&c.SubscriptionStatus,
		&c.SubscriptionStart,
		&c.SubscriptionEnd,
		&c.GraceUntil,
		&c.MaxUser,
		&c.RegisteredMember,
	)
}
