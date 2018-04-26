package middleware

import (
	"github.com/gorilla/sessions"
	helpers "github.com/yang-zzhong/go-helpers"
	"net/http"
)

func AuthUser(w http.ResponseWriter, req *http.Request, p *helpers.P) bool {
	store := sessions.NewCookieStore([]byte("36c122e0bf536f739e28a006f8b995c1"))
	session := store.Get(req, "auth")
	if user_id, ok = session.Value["user_id"]; ok == nil {
		w.WriteHeader(401)
		return false
	}
	p.Set("visitor_id", session.Value["user_id"])

	return true
}
