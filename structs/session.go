package structs

import (
	"database/sql"
	"time"
)

type UserSession struct {
	Userid             uint
	UserIP             string
	UserAgent          string
	UserState          string
	GoogleCode         string
	GoogleAccessToken  string
	GoogleRefreshToken string
	GoogleTokenType    string
	GoogleScope        string
	GoogleIDToken      string
	GoogleExpireIn     time.Time
	SessionID          string
	SessionStart       time.Time
	SessionExpire      time.Time
	SessionSignOut     sql.NullTime
}
