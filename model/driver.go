package model

import (
	. "boo-blog/config"
	"database/sql"
	. "database/sql/driver"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/yang-zzhong/go-model"
	. "github.com/yang-zzhong/go-querybuilder"
	"reflect"
	"time"
)

type config struct {
	driver   string
	host     string
	port     string
	username string
	password string
	database string
}

var conf config

func InitDriver() {
	conf.driver = Config.DB.Driver
	conf.host = Config.DB.Host
	conf.port = Config.DB.Port
	conf.username = Config.DB.UserName
	conf.password = Config.DB.Password
	conf.database = Config.DB.Database
}

func driver() (conn *sql.DB, err error) {
	conn, err = sql.Open(conf.driver, dsn())
	return
}

func dsn() string {
	dsn := conf.username + ":" + conf.password + "@"
	if conf.host != "" {
		dsn += "tcp(" + conf.host
		if conf.port != "" {
			dsn += ":" + conf.port
		}
		dsn += ")"
	}
	dsn += "/" + conf.database + "?parseTime=true"
	return dsn
}

func CreateModel(model interface{}) interface{} {
	mValue := reflect.ValueOf(model).Elem()
	IdValue := mValue.FieldByName("Id")
	IdValue.SetString(model.(Model).NewId().(string))
	now := reflect.ValueOf(NullTime{time.Now(), false})
	CreatedAt := mValue.FieldByName("CreatedAt")
	CreatedAt.Set(now)
	UpdatedAt := mValue.FieldByName("UpdatedAt")
	UpdatedAt.Set(now)

	return model
}

func CreateRepo(model interface{}) (repo *Repo, err error) {
	driver, err := driver()
	if err != nil {
		return
	}
	repo = NewRepo(model, driver, &MysqlModifier{})
	return
}

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
