package middleware

import (
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	"net/http"
)

type acrossDomain struct{}

func (ad *acrossDomain) Before(w http.ResponseWriter, _ *httprouter.Request, _ *helpers.P) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "id")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	return true
}

func (ad *acrossDomain) After(w http.ResponseWriter, _ *httprouter.Request, _ *helpers.P) bool {
	return true
}

var AcrossDomain acrossDomain
