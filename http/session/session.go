package session

import "github.com/golrilla/sessions"

var Store sessions.Store

func InitStore(sessionSecret string) {
	Store = sessions.NewCookieStore([]byte(sessionSecret))
}
