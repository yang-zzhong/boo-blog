package route

import (
	"boo-blog/http/middleware"
	. "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	"io"
	. "net/http"
)

func Router(docRoot string) *httprouter.Router {
	router := httprouter.NewRouter()
	router.DocRoot = docRoot
	router.BeforeApi = func(w ResponseWriter, req *httprouter.Request, p *P) bool {
		if req.Request.Method == MethodOptions {
			ad := middleware.AcrossDomain
			(&ad).Before(w, req, p)
			return false
		}
		return true
	}
	ms := []httprouter.Middleware{
		&middleware.UsedTime,
		&middleware.AuthUser,
		// &middleware.DB,
		&middleware.AcrossDomain,
	}
	router.Group("/api", ms, func(router *httprouter.Router) {
		registerPublicRoute(router)
		router.Group("", []httprouter.Middleware{
			&middleware.MustAuthUser,
			// &middleware.MustContactAuthedUser,
		}, registerNeedAuthRoute)
	})
	return router
}

func registerNeedAuthRoute(router *httprouter.Router) {
	registerUserAuthRoutes(router)
	registerBlogAuthRoutes(router)
	registerImageAuthRoutes(router)
	registerBlogInfoRoutes(router)
	registerTagAuthRoutes(router)
	registerCateAuthRoutes(router)
	registerThemeAuthRoutes(router)
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
