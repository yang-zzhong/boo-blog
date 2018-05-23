package middleware

import (
	"boo-blog/http/session"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	"net/http"
)

func AuthUser(w http.ResponseWriter, req *httprouter.Request, p *helpers.P) bool {
	s, _ := session.Store.Get(req.Request, "auth")
	s.Options.Domain = "192.168.3.206:8081"
	var userId interface{}
	var ok bool
	if userId, ok = s.Values["user_id"]; !ok || userId == nil {
		w.WriteHeader(401)
		return false
	}
	p.Set("visitor_id", userId)

	return true
}
