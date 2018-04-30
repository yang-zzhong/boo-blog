package model

import (
	"database/sql"
	"github.com/go-ini/ini"
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
var conn *sql.DB
var connected bool

func InitDriver(config *ini.Section) {
	conf.driver = config.Key("driver").String()
	conf.host = config.Key("host").String()
	conf.port = config.Key("port").String()
	conf.username = config.Key("username").String()
	conf.password = config.Key("password").String()
	conf.database = config.Key("database").String()
}

func init() {
	connected = false
}

func driver() *sql.DB {
	if connected {
		return conn
	}
	var err error
	if conn, err = sql.Open(conf.driver, dsn()); err != nil {
		panic(err)
	}
	connected = true
	return conn
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
	mm := NewModelMapper(model)
	if item, ok := mm.FnFds[model.(Model).PK()]; ok {
		IdValue := mValue.FieldByName(item.Name)
		IdValue.SetString(model.(IdMaker).NewId().(string))
	}

	return model
}

func CreateRepo(model interface{}) (repo *Repo, err error) {
	repo = NewRepo(model, driver(), &MysqlModifier{})
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
