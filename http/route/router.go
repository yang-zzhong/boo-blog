package route

import (
	. "/booblog/config"
	httprouter "github.com/yang-zzhong/go-httprouter"
)

func Route() *httprouter.Router {
	router := httprouter.CreateRouter(
		Config.Server.DocumentRoot, Config.Server.Indexes,
	)
	registerRoute(router)
	return router
}

func registerRoute(router *httprouter.Router) {

}
