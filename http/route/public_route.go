package route

import (
	"boo-blog/http/controller"
	. "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "net/http"
)

func registerBlogPublicRoutes(router *httprouter.Router) {
	router.Get("/users/:user_id/blogs/:url_id", func(w ResponseWriter, req *httprouter.Request, p *P) {
		blog := &controller.Article{controller.NewController(w)}
		blog.FetchUserBlog(req, p)
	})
	router.Get("/blogs", func(w ResponseWriter, req *httprouter.Request, _ *P) {
		blog := &controller.Article{controller.NewController(w)}
		blog.Find(req)
	})
	router.Get("/blogs/:id", func(w ResponseWriter, req *httprouter.Request, p *P) {
		blog := &controller.Article{controller.NewController(w)}
		blog.GetOne(req, p)
	})
}

func registerCatePublicRoutes(router *httprouter.Router) {
	router.Get("/cates", func(w ResponseWriter, req *httprouter.Request, _ *P) {
		cate := &controller.Category{controller.NewController(w)}
		cate.Find(req)
	})
	router.Get("/users/:user_id/article-used-cates", func(w ResponseWriter, req *httprouter.Request, p *P) {
		cate := &controller.Category{controller.NewController(w)}
		cate.ArticleUsed(req, p)
	})
	router.Get("/users/:user_id/image-used-cates", func(w ResponseWriter, req *httprouter.Request, p *P) {
		cate := &controller.Category{controller.NewController(w)}
		cate.ImageUsed(req, p)
	})
}

func registerTagPublicRoutes(router *httprouter.Router) {
	router.Get("/tags", func(w ResponseWriter, req *httprouter.Request, _ *P) {
		tag := &controller.Tag{controller.NewController(w)}
		tag.Search(req)
	})
	router.Get("/users/:user_id/used-tags", func(w ResponseWriter, req *httprouter.Request, p *P) {
		tag := &controller.Tag{controller.NewController(w)}
		tag.ArticleUsed(req, p)
	})
}

func registerUserPublicRoutes(router *httprouter.Router) {
	router.Get("/users/:name", func(w ResponseWriter, req *httprouter.Request, p *P) {
		user := &controller.User{controller.NewController(w)}
		user.One(req, p)
	})
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
}

func registerImagePublicRoutes(router *httprouter.Router) {
	router.Get("/images", func(w ResponseWriter, req *httprouter.Request, _ *P) {
		image := &controller.Image{controller.NewController(w)}
		image.Find(req)
	})
	router.Get("/images/:id", func(w ResponseWriter, req *httprouter.Request, p *P) {
		image := &controller.Image{controller.NewController(w)}
		image.Get(req, p)
	})
}
