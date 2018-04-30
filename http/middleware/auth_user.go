package middleware

import (
	"boo-blog/http/session"
	"github.com/gorilla/sessions"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	"net/http"
)

func AuthUser(w http.ResponseWriter, req *httprouter.Request, p *helpers.P) bool {
	s, _ := session.Store.Get(req.Request, "auth")
	var userId interface{}
	var ok bool
	if userId, ok = s.Values["user_id"]; !ok || userId == nil {
		w.WriteHeader(401)
		return false
	}
	p.Set("visitor_id", userId)

	return true
}
