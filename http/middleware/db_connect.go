package middleware

import (
	"boo-blog/model"
	"database/sql"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	m "github.com/yang-zzhong/go-model"
	. "github.com/yang-zzhong/go-querybuilder"
	"log"
	"net/http"
)

type connectDB struct {
	db *sql.DB
}

func (cdb *connectDB) Before(_ http.ResponseWriter, _ *httprouter.Request, p *helpers.P) bool {
	if db, err := model.OpenDB(); err != nil {
		log.Fatal("数据库连接失败")
		return false
	} else {
		log.Print("数据库连接成功")
		cdb.db = db
		m.Config(db, &MysqlModifier{})
		return true
	}
}

func (cdb *connectDB) After(_ http.ResponseWriter, _ *httprouter.Request, _ *helpers.P) bool {
	log.Print("断开数据库连接")
	cdb.db.Close()
	return true
}

var DB connectDB
