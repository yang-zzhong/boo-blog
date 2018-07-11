package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	m "github.com/yang-zzhong/go-model"
	query "github.com/yang-zzhong/go-querybuilder"
	"os"
	"reflect"
	"strings"
)

type IdMaker interface {
	NewId() interface{}
}

type Config struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	Database string
	ImageDir string
	BlogDir  string
}

var conf *Config

func InitDriver(config *Config) {
	conf = config
	sureDir(conf.ImageDir)
	sureDir(conf.BlogDir)
}

func OpenDB() error {
	if db, err := sql.Open(conf.Driver, dsn()); err != nil {
		return err
	} else {
		m.Config(db, &query.MysqlModifier{})
	}
	return nil
}

func CloseDB() {
	m.Conn.DB.Close()
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
	dsn := conf.Username + ":" + conf.Password + "@"
	if conf.Host != "" {
		dsn += "tcp(" + conf.Host
		if conf.Port != "" {
			dsn += ":" + conf.Port
		}
		dsn += ")"
	}
	dsn += "/" + conf.Database + "?parseTime=true"
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
