package http

import (
	"boo-blog/http/route"
	"boo-blog/http/session"
	"log"
	"net"
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
	if l, err := net.Listen("tcp4", ":"+Http.Port); err != nil {
		return err
	} else {
		log.Fatal(http.Serve(l, route.Router(Http.DocRoot)))
	}

	return nil
}

func StartTLS(certFile, keyFile string) error {
	session.InitStore(Http.SessionSecret)
	log.Print("listen on :" + Http.Port)
	Http.running = true
	if l, err := net.Listen("tcp4", ":"+Http.Port); err != nil {
		return err
	} else {
		log.Fatal(http.ServeTLS(l, route.Router(Http.DocRoot), certFile, keyFile))
	}

	return nil
}
