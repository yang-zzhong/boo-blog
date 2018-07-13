package middleware

import (
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	"log"
	"net/http"
)

type connectDB struct{}

func (cdb *connectDB) Before(_ http.ResponseWriter, _ *httprouter.Request, p *helpers.P) bool {
	if err := model.OpenDB(); err != nil {
		log.Fatal("数据库连接失败")
		return false
	}
	log.Print("数据库连接成功")
	return true
}

func (cdb *connectDB) After(_ http.ResponseWriter, _ *httprouter.Request, _ *helpers.P) bool {
	log.Print("断开数据库连接")
	model.CloseDB()
	return true
}

var DB connectDB
