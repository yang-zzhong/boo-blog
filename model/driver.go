package model

import (
	"database/sql"
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
	model "github.com/yang-zzhong/go-model"
	. "github.com/yang-zzhong/go-querybuilder"
	"os"
	"reflect"
	"strings"
)

type IdMaker interface {
	NewId() interface{}
}

type config struct {
	driver    string
	host      string
	port      string
	username  string
	password  string
	database  string
	image_dir string
	blog_dir  string
}

var conf config
var DB *sql.DB

func InitDriver(config *ini.Section) {
	conf.driver = config.Key("driver").String()
	conf.host = config.Key("host").String()
	conf.port = config.Key("port").String()
	conf.username = config.Key("username").String()
	conf.password = config.Key("password").String()
	conf.database = config.Key("database").String()
	conf.image_dir = config.Key("image_dir").String()
	conf.blog_dir = config.Key("blog_dir").String()
	if conn, err := sql.Open(conf.driver, dsn()); err != nil {
		panic(err)
	} else {
		model.Config(conn, &MysqlModifier{})
		DB = conn
	}
	sureDir(conf.image_dir)
	sureDir(conf.blog_dir)
}

func sureDir(dir string) {
	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
		}
		if err != nil {
			panic(err)
		}
	}
	if !fi.IsDir() {
		err = os.MkdirAll(dir, 0755)
	}
	if err != nil {
		panic(err)
	}
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

func nullArrayDBValue(value interface{}) interface{} {
	result := strings.Join(value.([]string), ",")
	return result
}

func nullArrayValue(value interface{}) (result reflect.Value) {
	v := value.(sql.NullString)
	if v.Valid {
		val, _ := v.Value()
		result = reflect.ValueOf(strings.Split(val.(string), ","))
	} else {
		result = reflect.ValueOf([]string{})
	}
	return
}
