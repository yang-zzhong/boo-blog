package controller

import (
	"github.com/gorilla/sessions"
	httprouter "github.com/yang-zzhong/go-httprouter"
)

type Login struct{ *Controller }

func (this *Login) Register(req *httprouter.Request) {
	var repo *Repo
	var err error
	if repo, err = model.NewUserRepo(); err != nil {
		this.InternalError(err)
		return
	}
	user := model.NewUser()
	user.Name = req.FormValue("name")
	user.NickName = user.Name
	user.EmailAddr = req.FormValue("email_addr")
	user.Password = user.Encrypt(req.FormValue("password"))
	repo.Where("name", user.Name).Or().
		Where("email_addr", user.EmailAddr).Or().
		Where("phone_number", user.PhoneNumber)
	if repo.Count() > 0 {
		this.String("电话或者邮箱或者用户名已被使用", 500)
		return
	}
	if err = repo.Create(user); err != nil {
		this.InternalError(err)
		return
	}
}

func (this *Login) Login(req *httprouter.Request) {
	var repo *Repo
	var err error
	var account string
	var password string
	var m interface{}
	if account = req.FormValue("account"); account == "" {
		this.String("没有制定账号", 500)
		return
	}
	if repo, err = model.NewUserRepo(); err != nil {
		this.InternalError(nil)
		return
	}
	repo.Where("name", account).Or().
		Where("email_addr", account).Or().
		Where("phone_number", account)
	if m = repo.One(); m == nil {
		this.String("用户名或密码不正确", 500)
		return
	}
	user := m.(*User)
	if user.Encrypt(req.FormValue("password")) != user.Password {
		this.String("用户名或密码不正确", 500)
		return
	}
	store := sessions.NewCookieStore([]byte("36c122e0bf536f739e28a006f8b995c1"))

	session := store.Get(req, "auth")
	session.Value["user_id"] = user.Id
	session.Save()
}

func (this *Login) Logout(req *httprouter.Request) {
	var repo *Repo
	var err error

	store := sessions.NewCookieStore([]byte("36c122e0bf536f739e28a006f8b995c1"))
	session := store.Get(req, "auth")
	session.Value["user_id"] = nil

	session.Save()
}
