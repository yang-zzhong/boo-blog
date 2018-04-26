package route

import (
	. "boo-blog/config"
	"boo-blog/http/controller"
	"boo-blog/http/middleware"
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
	router.Get("/images/:id", func(w ResponseWriter, req *httprouter.Request, p *P) {
		image := &controller.Image{controller.NewController(w)}
		image.Get(req, p)
	})
	ms := httprouter.NewMs()
	ms.Append(middleware.AuthUser)
	router.Group("", httprouter.NewMs(), func(router *httprouter.Router) {
		router.Post("/users", func(w ResponseWriter, req *httprouter.Request, _ *P) {
			user := &controller.User{controller.NewController(w)}
			user.CreateUser(req)
		})
		router.Post("/images", func(w ResponseWriter, req *httprouter.Request, _ *P) {
			image := &controller.Image{controller.NewController(w)}
			image.Create(req)
		})
		router.Post("/image-cates", func(w ResponseWriter, req *httprouter.Request, p *P) {
			imageGroup := &controller.Category{controller.NewController(w)}
			imageGroup.Create(req, p)
		})
	})
}
