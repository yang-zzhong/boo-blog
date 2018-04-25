package model

import (
	. "boo-blog/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/yang-zzhong/go-model"
	. "github.com/yang-zzhong/go-querybuilder"
	"reflect"
	"time"
)

type IdMaker interface {
	NewId() interface{}
}

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
	IdValue := mValue.FieldByName(model.(Model).PK())
	IdValue.SetString(model.(IdMaker).NewId().(string))

	return model
}

func CreateRepo(model interface{}) (repo *Repo, err error) {
	driver, err := driver()
	if err != nil {
		return
	}
	repo = NewRepo(model, driver, &MysqlModifier{})
	repo.OnUpdate(func(model interface{}) {
		mValue := reflect.ValueOf(model).Elem()
		mValue.FieldByName("UpdatedAt").Set(reflect.ValueOf(time.Now()))
	})
	repo.OnCreate(func(model interface{}) {
		mValue := reflect.ValueOf(model).Elem()
		now := reflect.ValueOf(time.Now())
		mValue.FieldByName("CreatedAt").Set(now)
		mValue.FieldByName("UpdatedAt").Set(now)
	})
	return
}
