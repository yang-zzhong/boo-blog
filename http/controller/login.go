package controller

import (
	"boo-blog/model"
	"database/sql"
	"github.com/gorilla/sessions"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-model"
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
	user.NickName = sql.NullString{user.Name, true}
	user.EmailAddr = sql.NullString{req.FormValue("email_addr"), true}
	user.Password = user.Encrypt(req.FormValue("password"))
	repo.Where("name", req.FormValue("name")).Or().
		Where("email_addr", req.FormValue("email_addr")).Or().
		Where("phone_number", req.FormValue("phone_number"))
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
	u := m.(*model.User)
	if u.Encrypt(req.FormValue("password")) != u.Password {
		this.String("用户名或密码不正确", 500)
		return
	}
	store := sessions.NewCookieStore([]byte("36c122e0bf536f739e28a006f8b995c1"))

	session, _ := store.Get(req.Request, "auth")
	session.Values["user_id"] = u.Id
	session.Save(req.Request, this.ResponseWriter())
}

func (this *Login) Logout(req *httprouter.Request) {
	store := sessions.NewCookieStore([]byte("36c122e0bf536f739e28a006f8b995c1"))
	session, _ := store.Get(req.Request, "auth")
	session.Values["user_id"] = nil

	session.Save(req.Request, this.ResponseWriter())
}
