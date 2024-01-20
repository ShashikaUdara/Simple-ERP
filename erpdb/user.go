package erpdb

import "time"

type User struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type UserSession struct {
	ID               int       `db:"id"`
	UserID           string    `db:"user_id"`
	SessionToken     string    `db:"session_token"`
	IPAddress        string    `db:"ip_address"`
	UserAgent        string    `db:"user_agent"`
	CreationTime     time.Time `db:"creation_time"`
	LastActivityTime time.Time `db:"last_activity_time"`
}
