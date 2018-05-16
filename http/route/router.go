package route

import (
	"boo-blog/http/controller"
	"boo-blog/http/middleware"
	. "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	"io"
	. "net/http"
)

func Router(docRoot string) *httprouter.Router {
	router := httprouter.NewRouter()
	router.DocRoot = docRoot
	router.Group("/api", httprouter.NewMs(), func(router *httprouter.Router) {
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
	registerTagAuthRoutes(router)
	registerCateAuthRoutes(router)
}

func registerPublicRoute(router *httprouter.Router) {
	router.Get("/hello-world", func(w ResponseWriter, req *httprouter.Request, _ *P) {
		io.WriteString(w, "hello world")
	})
	registerBlogPublicRoutes(router)
	registerCatePublicRoutes(router)
	registerTagPublicRoutes(router)
	registerUserPublicRoutes(router)
	registerImagePublicRoutes(router)
}
