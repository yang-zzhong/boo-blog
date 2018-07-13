package session

import "github.com/gorilla/sessions"

var Store sessions.Store

func InitStore(sessionSecret string) {
	Store = sessions.NewCookieStore([]byte(sessionSecret))
}
