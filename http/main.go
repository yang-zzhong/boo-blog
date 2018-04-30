package main

import (
	"boo-blog/http/route"
	"boo-blog/http/session"
	"boo-blog/model"
	"github.com/dmulholland/args"
	"github.com/go-ini/ini"
	"log"
	"net/http"
)

//
// http server app
// boohttp -c /path/to/configfile
// boohttp
//
const (
	APP_NAME       = "boohttp"
	VERSION        = "0.0.1"
	DEFAULT_CONFIG = "./http.ini"
)

type server struct {
	config *ini.File
}

func NewServer(configFile string) (s *server, err error) {
	s = new(server)
	if s.config, err = ini.Load(configFile); err != nil {
		return
	}
	return
}

func (s *server) Config() *ini.File {
	return s.config
}

func (s *server) StartServer() {
	root := s.Config().Section("server").Key("doc_root").String()
	sessionSecret := s.Config().Section("server").Key("session_secret").String()
	session.InitStore(sessionSecret)
	router := route.Router(root)
	port := s.Config().Section("server").Key("port").String()
	serverAddr := ":" + port
	log.Print("listen on " + serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}

func main() {
	parser := args.NewParser()
	parser.Helptext = help(APP_NAME, VERSION)
	parser.Version = VERSION
	parser.NewString("-c", DEFAULT_CONFIG)
	parser.Parse()
	s, err := NewServer(parser.GetString("-c"))
	if err != nil {
		panic(err)
	}
	model.InitDriver(s.Config().Section("database"))
	s.StartServer()
}

func help(appName string, version string) string {
	return appName + version
}
