package route

import (
	"boo-blog/http/controller"
	qrcode "github.com/skip2/go-qrcode"
	. "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
)

func registerBlogPublicRoutes(router *httprouter.Router) {
	router.Get("/users/:user_id/blogs/:url_id", func(w *httprouter.ResponseWriter, req *httprouter.Request, p *P) {
		blog := &controller.Article{controller.NewController(w)}
		blog.FetchUserBlog(req, p)
	})
	router.Get("/blogs", func(w *httprouter.ResponseWriter, req *httprouter.Request, p *P) {
		blog := &controller.Article{controller.NewController(w)}
		blog.Find(req, p)
	})
	router.Get("/blogs/:id", func(w *httprouter.ResponseWriter, req *httprouter.Request, p *P) {
		blog := &controller.Article{controller.NewController(w)}
		blog.GetOne(req, p)
	})
	router.Get("/blogs/:blog_id/comments", func(w *httprouter.ResponseWriter, req *httprouter.Request, p *P) {
		comment := &controller.Comment{controller.NewController(w)}
		comment.Articles(req, p)
	})
}

func registerCatePublicRoutes(router *httprouter.Router) {
	router.Get("/cates", func(w *httprouter.ResponseWriter, req *httprouter.Request, _ *P) {
		cate := &controller.Category{controller.NewController(w)}
		cate.Find(req)
	})
	router.Get("/users/:user_id/article-used-cates", func(w *httprouter.ResponseWriter, req *httprouter.Request, p *P) {
		cate := &controller.Category{controller.NewController(w)}
		cate.ArticleUsed(req, p)
	})
	router.Get("/users/:user_id/image-used-cates", func(w *httprouter.ResponseWriter, req *httprouter.Request, p *P) {
		cate := &controller.Category{controller.NewController(w)}
		cate.ImageUsed(req, p)
	})
}

func registerTagPublicRoutes(router *httprouter.Router) {
	router.Get("/tags", func(w *httprouter.ResponseWriter, req *httprouter.Request, _ *P) {
		tag := &controller.Tag{controller.NewController(w)}
		tag.Search(req)
	})
	router.Get("/users/:user_id/used-tags", func(w *httprouter.ResponseWriter, req *httprouter.Request, p *P) {
		tag := &controller.Tag{controller.NewController(w)}
		tag.ArticleUsed(req, p)
	})
}

func registerUserPublicRoutes(router *httprouter.Router) {
	router.Get("/users/:user_id/profile", func(w *httprouter.ResponseWriter, _ *httprouter.Request, p *P) {
		user := &controller.User{controller.NewController(w)}
		user.Profile(p)
	})
	router.Get("/users/:name", func(w *httprouter.ResponseWriter, req *httprouter.Request, p *P) {
		user := &controller.User{controller.NewController(w)}
		user.One(req, p)
	})
	router.Get("/users", func(w *httprouter.ResponseWriter, req *httprouter.Request, p *P) {
		user := &controller.User{controller.NewController(w)}
		user.Find(req, p)
	})
	router.Post("/register", func(w *httprouter.ResponseWriter, req *httprouter.Request, _ *P) {
		login := &controller.Login{controller.NewController(w)}
		login.Register(req)
	})
	router.Post("/login", func(w *httprouter.ResponseWriter, req *httprouter.Request, _ *P) {
		login := &controller.Login{controller.NewController(w)}
		login.Login(req)
	})
	router.Delete("/logout", func(w *httprouter.ResponseWriter, req *httprouter.Request, _ *P) {
		login := &controller.Login{controller.NewController(w)}
		login.Logout(req)
	})
}

func registerImagePublicRoutes(router *httprouter.Router) {
	router.Get("/images", func(w *httprouter.ResponseWriter, req *httprouter.Request, _ *P) {
		image := &controller.Image{controller.NewController(w)}
		image.Find(req)
	})
	router.Get("/images/:id", func(w *httprouter.ResponseWriter, req *httprouter.Request, p *P) {
		image := &controller.Image{controller.NewController(w)}
		image.Get(req, p)
	})
}

func registerQrCodeRoute(router *httprouter.Router) {
	router.Get("/qr-code", func(w *httprouter.ResponseWriter, req *httprouter.Request, _ *P) {
		if png, err := qrcode.Encode(req.FormValue("url"), qrcode.Medium, 256); err != nil {
			w.InternalError(err)
			return
		} else {
			w.WithHeader("Content-Type", "image/png")
			if _, err := w.Write(png); err != nil {
				w.InternalError(err)
			}
		}
	})
}
