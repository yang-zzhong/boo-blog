package route

import (
	"boo-blog/http/controller"
	. "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "net/http"
)

func registerBlogAuthRoutes(router *httprouter.Router) {
	router.Post("/blogs", func(w ResponseWriter, req *httprouter.Request, p *P) {
		blog := &controller.Article{controller.NewController(w)}
		blog.Create(req, p)
	})
	router.Put("/blogs/:id", func(w ResponseWriter, req *httprouter.Request, p *P) {
		blog := &controller.Article{controller.NewController(w)}
		blog.Update(req, p)
	})
	router.Delete("/blogs/:id", func(w ResponseWriter, req *httprouter.Request, p *P) {
		blog := &controller.Article{controller.NewController(w)}
		blog.Remove(req, p)
	})
}

func registerImageAuthRoutes(router *httprouter.Router) {
	router.Post("/images", func(w ResponseWriter, req *httprouter.Request, p *P) {
		image := &controller.Image{controller.NewController(w)}
		image.Create(req, p)
	})
}

func registerTagAuthRoutes(router *httprouter.Router) {
	router.Post("/tags", func(w ResponseWriter, req *httprouter.Request, p *P) {
		tag := &controller.Tag{controller.NewController(w)}
		tag.Create(req, p)
	})
}

func registerCateAuthRoutes(router *httprouter.Router) {
	router.Post("/cates", func(w ResponseWriter, req *httprouter.Request, p *P) {
		cate := &controller.Category{controller.NewController(w)}
		cate.Create(req, p)
	})
	router.Put("/cates/:id", func(w ResponseWriter, req *httprouter.Request, p *P) {
		cate := &controller.Category{controller.NewController(w)}
		cate.Update(req, p)
	})
	router.Delete("/cates/:id", func(w ResponseWriter, req *httprouter.Request, p *P) {
		cate := &controller.Category{controller.NewController(w)}
		cate.Delete(req, p)
	})
}
