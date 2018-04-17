package main

import (
	. "boo-blog/config"
	"boo-blog/http/route"
	. "boo-blog/log"
	"boo-blog/model"
	"github.com/dmulholland/args"
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
	DEFAULT_CONFIG = "./http.conf"
)

func main() {
	parser := args.NewParser()
	parser.Helptext = help(APP_NAME, VERSION)
	parser.Version = VERSION
	parser.NewString("-c", DEFAULT_CONFIG)
	parser.Parse()
	Logger().Print("hello world")
	InitConfig(parser.GetString("-c"))
	model.InitDriver()

	router := route.Router()
	// serverAddr := Config.Server.Domain + ":" + Config.Server.Port
	serverAddr := ":" + Config.Server.Port
	Logger().Print("listen on " + serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}

func help(appName string, version string) string {
	return appName + version
}
