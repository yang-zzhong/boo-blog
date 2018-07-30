package route

import (
	"boo-blog/http/controller"
	. "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "net/http"
)

func registerThemeAuthRoutes(router *httprouter.Router) {
	router.Post("/applied-themes", func(w ResponseWriter, req *httprouter.Request, p *P) {
		theme := &controller.Theme{controller.NewController(w)}
		theme.Apply(req, p)
	})
	router.Post("/themes", func(w ResponseWriter, req *httprouter.Request, p *P) {
		theme := &controller.Theme{controller.NewController(w)}
		theme.Create(req, p)
	})
	router.Get("/themes", func(w ResponseWriter, req *httprouter.Request, p *P) {
		theme := &controller.Theme{controller.NewController(w)}
		theme.Find(p)
	})
	router.Put("/themes/:theme_id", func(w ResponseWriter, req *httprouter.Request, p *P) {
		theme := &controller.Theme{controller.NewController(w)}
		theme.Update(req, p)
	})
	router.Delete("/themes/:theme_id", func(w ResponseWriter, _ *httprouter.Request, p *P) {
		theme := &controller.Theme{controller.NewController(w)}
		theme.Delete(p)
	})
}

func registerBlogAuthRoutes(router *httprouter.Router) {
	router.Post("/blogs/:blog_id/thumb-up", func(w ResponseWriter, _ *httprouter.Request, p *P) {
		vote := &controller.Vote{controller.NewController(w)}
		p.Set("vote", 1)
		vote.Create(p)
	})
	router.Post("/blogs/:blog_id/thumb-down", func(w ResponseWriter, _ *httprouter.Request, p *P) {
		vote := &controller.Vote{controller.NewController(w)}
		p.Set("vote", -1)
		vote.Create(p)
	})
	router.Post("/blogs/:blog_id/unthumb", func(w ResponseWriter, _ *httprouter.Request, p *P) {
		vote := &controller.Vote{controller.NewController(w)}
		vote.Delete(p)
	})
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
	router.Post("/blogs/:blog_id/comments", func(w ResponseWriter, req *httprouter.Request, p *P) {
		comment := &controller.Comment{controller.NewController(w)}
		comment.Create(req, p)
	})
}

func registerImageAuthRoutes(router *httprouter.Router) {
	router.Post("/images", func(w ResponseWriter, req *httprouter.Request, p *P) {
		image := &controller.Image{controller.NewController(w)}
		image.Create(req, p)
	})
	router.Put("/images/to/:cate_id", func(w ResponseWriter, req *httprouter.Request, p *P) {
		image := &controller.Image{controller.NewController(w)}
		image.MoveTo(req, p)
	})
}

func registerBlogInfoRoutes(router *httprouter.Router) {
	router.Post("/user-info", func(w ResponseWriter, req *httprouter.Request, p *P) {
		user := &controller.User{controller.NewController(w)}
		user.SaveUserInfo(req, p)
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
