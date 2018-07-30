package middleware

import (
	"boo-blog/http/session"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	"net/http"
)

type authUser struct{}

func (au *authUser) Before(w http.ResponseWriter, req *httprouter.Request, p *helpers.P) bool {
	user, logged := session.User(req.Header.Get("id"))
	if !logged {
		return true
	}
	p.Set("visitor_id", user.Id)
	p.Set("visitor", user)

	return true
}

func (au *authUser) After(_ http.ResponseWriter, _ *httprouter.Request, _ *helpers.P) bool {
	return true
}

var AuthUser authUser

type mustAuthUser struct{}

func (au *mustAuthUser) Before(w http.ResponseWriter, req *httprouter.Request, p *helpers.P) bool {
	if p.Get("visitor_id") == nil {
		w.WriteHeader(401)
		return false
	}
	return true
}

func (au *mustAuthUser) After(_ http.ResponseWriter, _ *httprouter.Request, _ *helpers.P) bool {
	return true
}

var MustAuthUser mustAuthUser
