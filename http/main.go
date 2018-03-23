package main

import (
	. "/boo-blog/config"
	"/boo-blog/http/route"
	"github.com/dmulholland/args"
	"net/http"
)

//
// http server app
// boohttp -c /path/to/configfile
// boohttp
//
const (
	APP_NAME = "boohttp"
	VERSION  = "0.0.1"
)

func main() {
	parser := args.NewParser()
	parser.Helptext = help(APP_NAME, VERSION)
	parser.Version = VERSION
	parser.NewString("--config", DEFAULT_CONFIG)
	parser.Parse()

	InitConfig(parser.GetString("--config"))

	startHttpServer()
}

func startHttpServer() {
	router := route.Route()
}

func help(appName string, version string) string {
	return appName + version
}
