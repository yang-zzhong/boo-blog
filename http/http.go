package http

type Httpd struct {
	DocRoot       string
	Port          string
	SessionSecret string
	running       bool
}

var Http *Httpd

func InitHttp(docRoot, port, sessionSecret string) error {
	Http := new(Httpd)
	Http.DocRoot = docRoot
	Http.Port = port
	Http.SessionSecret = sessionSecret
	Http.running = false
}

func Start() error {
	session.InitStore(Http.SessionSecret)
	log.Print("listen on :" + Http.Port)
	Http.running = true
	go log.Fatal(http.ListenAndServe(":"+Http.Port, route.Router(Http.DocRoot)))
}
