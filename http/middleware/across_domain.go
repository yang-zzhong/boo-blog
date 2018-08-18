package middleware

import (
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
)

type acrossDomain struct{}

func (ad *acrossDomain) Before(w *httprouter.ResponseWriter, _ *httprouter.Request, _ *helpers.P) bool {
	w.WithHeader("Access-Control-Allow-Origin", "*")
	w.WithHeader("Access-Control-Allow-Headers", "id")
	w.WithHeader("Access-Control-Allow-Methods", "*")
	return true
}

func (ad *acrossDomain) After(_ *httprouter.ResponseWriter, _ *httprouter.Request, _ *helpers.P) bool {
	return true
}

var AcrossDomain acrossDomain
