package model

type Article struct {
	Id     string `db:"id char(32) pk"`
	Title  string `db:"title varchar(64)"`
	UserId string `db:"user_id char(32)"`
}
