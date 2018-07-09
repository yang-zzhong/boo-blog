package http

import (
	"boo-blog/http/route"
	"boo-blog/http/session"
	"log"
	"net/http"
)

type Httpd struct {
	DocRoot       string
	Port          string
	SessionSecret string
	running       bool
}

var Http Httpd

func InitHttp(docRoot, port, sessionSecret string) {
	Http = Httpd{docRoot, port, sessionSecret, false}
}

func Start() error {
	session.InitStore(Http.SessionSecret)
	log.Print("listen on :" + Http.Port)
	Http.running = true
	go log.Fatal(http.ListenAndServe(":"+Http.Port, route.Router(Http.DocRoot)))

	return nil
}
