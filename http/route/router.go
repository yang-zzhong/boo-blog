package route

import (
	"boo-blog/http/middleware"
	. "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	"io"
	"log"
	. "net/http"
)

func Router(docRoot string) *httprouter.Router {
	defer func() {
		if r := recover(); r != nil {
			log.Print(r)
		}
	}()
	router := httprouter.NewRouter()
	router.DocRoot = docRoot
	router.Before = func(w ResponseWriter, req *httprouter.Request, p *P) bool {
		if req.Request.Method == MethodOptions {
			middleware.AcrossDomain(w, req, p)
			return false
		}
		return true
	}
	ms := httprouter.NewMs()
	ms.Append(middleware.AcrossDomain)
	router.Group("/api", ms, func(router *httprouter.Router) {
		registerPublicRoute(router)
		ms := httprouter.NewMs()
		ms.Append(middleware.AuthUser)
		router.Group("", ms, registerNeedAuthRoute)
	})
	return router
}

func registerNeedAuthRoute(router *httprouter.Router) {
	registerBlogAuthRoutes(router)
	registerImageAuthRoutes(router)
	registerBlogInfoRoutes(router)
	registerTagAuthRoutes(router)
	registerCateAuthRoutes(router)
}

func registerPublicRoute(router *httprouter.Router) {
	router.Get("/hello-world", func(w ResponseWriter, req *httprouter.Request, _ *P) {
		io.WriteString(w, "hello world")
	})
	registerQrCodeRoute(router)
	registerBlogPublicRoutes(router)
	registerCatePublicRoutes(router)
	registerTagPublicRoutes(router)
	registerUserPublicRoutes(router)
	registerImagePublicRoutes(router)
}
