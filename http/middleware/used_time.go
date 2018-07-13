package middleware

import (
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	"log"
	"net/http"
	"time"
)

type usedTime struct {
	begin time.Time
}

func (ut *usedTime) Before(_ http.ResponseWriter, _ *httprouter.Request, _ *helpers.P) bool {
	ut.begin = time.Now()
	return true
}

func (ut *usedTime) After(_ http.ResponseWriter, _ *httprouter.Request, _ *helpers.P) bool {
	log.Printf("used: %.5fs", time.Now().Sub(ut.begin).Seconds())
	return true
}

var UsedTime usedTime
