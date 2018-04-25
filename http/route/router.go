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
	router.Get("/hello-world", func(w ResponseWriter, req *httprouter.Request, _ *P) {
		io.WriteString(w, "hello world")
	})
	router.Get("/users", func(w ResponseWriter, req *httprouter.Request, _ *P) {
		user := controller.NewUser(w)
		user.RenderUsers(req)
	})
	router.Post("/users", func(w ResponseWriter, req *httprouter.Request, _ *P) {
		user := &controller.User{controller.NewController(w)}
		user.CreateUser(req)
	})
	router.Post("/images", func(w ResponseWriter, req *httprouter.Request, _ *P) {
		image := &controller.Image{controller.NewController(w)}
		image.Create(req)
	})
	router.Get("/images/:id", func(w ResponseWriter, req *httprouter.Request, p *P) {
		image := &controller.Image{controller.NewController(w)}
		image.Get(req, p)
	})
	router.Post("/image-groups", func(w ResponseWriter, req *httprouter.Request, p *P) {
		imageGroup := &controller.ImageGroup{controller.NewController(w)}
		imageGroup.Create(req)
	})
}
