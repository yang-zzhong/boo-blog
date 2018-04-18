package model

import (
	"database/sql"
)

type UserLogin struct {
	Id       string    `db:"id varcher(128) pk"`
	UserId   string    `db:"user_id varcher(128)"`
	City     string    `db:"city varchar(64)"`
	LoginAt  time.Time `db:"login_at datatime"`
	LogoutAt NullTime  `db:"logout_at datetime nil"`
}
