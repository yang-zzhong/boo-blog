package route

import (
	. "boo-blog/config"
	"boo-blog/http/controller"
	. "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	"io"
	. "net/http"
)

func Router() *httprouter.Router {
	router := httprouter.NewRouter()
	router.DocRoot = Config.Server.DocumentRoot
	registerRoute(router)
	return router
}

func registerRoute(router *httprouter.Router) {
	router.Get("/hello-world", func(w ResponseWriter, req *Request, _ *P) {
		io.WriteString(w, "hello world")
	})

	router.Get("/users", func(w ResponseWriter, req *Request, _ *P) {
		user := &controller.User{}
		user.RenderUsers(w, req)
	})

	router.Post("/users", func(w ResponseWriter, req *Request, _ *P) {
		user := &controller.User{}
		user.CreateUser(w, req)
	})
}
