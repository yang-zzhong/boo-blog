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
	router.Post("/register", func(w ResponseWriter, req *httprouter.Request, _ *P) {
		login := &controller.Login{controller.NewController(w)}
		login.Register(req)
	})
	router.Post("/login", func(w ResponseWriter, req *httprouter.Request, _ *P) {
		login := &controller.Login{controller.NewController(w)}
		login.Login(req)
	})
	router.Delete("/logout", func(w ResponseWriter, req *httprouter.Request, _ *P) {
		login := &controller.Login{controller.NewController(w)}
		login.Logout(req)
	})
	router.Get("/hello-world", func(w ResponseWriter, req *httprouter.Request, _ *P) {
		io.WriteString(w, "hello world")
	})
	router.Get("/imgs/:id", func(w ResponseWriter, req *httprouter.Request, p *P) {
		image := &controller.Image{controller.NewController(w)}
		image.Get(req, p)
	})
	router.Get("/tags", func(w ResponseWriter, req *httprouter.Request, _ *P) {
		tag := &controller.Tag{controller.NewController(w)}
		tag.Search(req)
	})
	ms := httprouter.NewMs()
	ms.Append(middleware.AuthUser)
	router.Group("", ms, func(router *httprouter.Router) {
		router.Post("/imgs", func(w ResponseWriter, req *httprouter.Request, p *P) {
			image := &controller.Image{controller.NewController(w)}
			image.Create(req, p)
		})
		router.Post("/tags", func(w ResponseWriter, req *httprouter.Request, p *P) {
			tag := &controller.Tag{controller.NewController(w)}
			tag.Create(req, p)
		})
		router.Post("/cates", func(w ResponseWriter, req *httprouter.Request, p *P) {
			cate := &controller.Category{controller.NewController(w)}
			cate.Create(req, p)
		})
		router.Put("/cates", func(w ResponseWriter, req *httprouter.Request, p *P) {
			cate := &controller.Category{controller.NewController(w)}
			cate.Update(req, p)
		})
		router.Get("/cates", func(w ResponseWriter, req *httprouter.Request, p *P) {
			cate := &controller.Category{controller.NewController(w)}
			cate.Get(req, p)
		})
	})
}
