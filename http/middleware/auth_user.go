package middleware

import (
	"boo-blog/http/session"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	"net/http"
)

func AuthUser(w http.ResponseWriter, req *httprouter.Request, p *helpers.P) bool {
	user, logged := session.User(req.Header.Get("id"))
	// s, _ := session.Store.Get(req.Request, "auth")
	// s.Options.Domain = "192.168.3.206:8081"
	if !logged {
		w.WriteHeader(401)
		return false
	}
	p.Set("visitor_id", user.Id)
	p.Set("visitor", user)

	return true
}
