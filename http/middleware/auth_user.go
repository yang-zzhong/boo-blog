package middleware

import (
	"boo-blog/http/session"
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	"net/http"
)

// 获取登录用户信息的中间件
type authUser struct{}

func (au *authUser) Before(w http.ResponseWriter, req *httprouter.Request, p *helpers.P) bool {
	user, logged := session.User(req.Header.Get("id"))
	if !logged {
		return true
	}
	p.Set("visitor_id", user.Id)
	p.Set("visitor", user)

	return true
}

func (au *authUser) After(_ http.ResponseWriter, _ *httprouter.Request, _ *helpers.P) bool {
	return true
}

// 判断用户是否登录的中间件
type mustAuthUser struct{}

func (au *mustAuthUser) Before(w http.ResponseWriter, req *httprouter.Request, p *helpers.P) bool {
	if p.Get("visitor_id") == nil {
		w.WriteHeader(401)
		return false
	}
	return true
}

func (au *mustAuthUser) After(_ http.ResponseWriter, _ *httprouter.Request, _ *helpers.P) bool {
	return true
}

// 判断用户是否验证电话或者邮箱的中间件
type mustContactAuthedUser struct{}

func (au *mustContactAuthedUser) Before(w http.ResponseWriter, req *httprouter.Request, p *helpers.P) bool {
	visitor := p.Get("visitor")
	if visitor == nil || visitor != nil && !visitor.(*model.User).EmailAddrAuthed {
		// 需要验证联系方式, 电话或者邮件
		w.WriteHeader(402)
		return false
	}
	return true
}

func (au *mustContactAuthedUser) After(_ http.ResponseWriter, _ *httprouter.Request, _ *helpers.P) bool {
	return true
}

var (
	AuthUser              authUser
	MustAuthUser          mustAuthUser
	MustContactAuthedUser mustContactAuthedUser
)
