package middleware

import (
	"github.com/gorilla/sessions"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	"net/http"
)

func AuthUser(w http.ResponseWriter, req *httprouter.Request, p *helpers.P) bool {
	store := sessions.NewCookieStore([]byte("36c122e0bf536f739e28a006f8b995c1"))
	session, _ := store.Get(req.Request, "auth")
	var userId interface{}
	var ok bool
	if userId, ok = session.Values["user_id"]; !ok {
		w.WriteHeader(401)
		return false
	}
	p.Set("visitor_id", userId)

	return true
}
