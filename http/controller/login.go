package controller

import (
	"boo-blog/http/session"
	"boo-blog/model"
	httprouter "github.com/yang-zzhong/go-httprouter"
)

type Login struct{ *Controller }

func (this *Login) Register(req *httprouter.Request) {
	if req.FormValue("name") == "" {
		this.String("没有制定用户名", 500)
		return
	}
	if len(req.FormValue("password")) < 12 {
		this.String("密码长度不够", 500)
		return
	}
	user := model.NewUser().Instance()
	user.Fill(map[string]interface{}{
		"name":       req.FormValue("name"),
		"nickname":   req.FormValue("name"),
		"blog_name":  req.FormValue("name") + "的博客",
		"email_addr": req.FormValue("email_addr"),
	})
	user.Password = user.Encrypt(req.FormValue("password"))
	user.Repo().Where("name", user.Name)
	if user.EmailAddr != "" {
		user.Repo().Or().Where("email_addr", req.FormValue("email_addr"))
	}
	if user.PhoneNumber != "" {
		user.Repo().Or().Where("phone_number", req.FormValue("phone_number"))
	}
	if count, err := user.Repo().Count(); err != nil {
		this.InternalError(err)
		return
	} else if count > 0 {
		this.String("电话或者邮箱或者用户名已被使用", 500)
		return
	}
	if user.EmailAddr != "" {
		if err := this.SendVerifyEmail(user); err != nil {
			this.InternalError(err)
			return
		}
	}
	if err := user.Save(); err != nil {
		this.InternalError(err)
	}
}

func (this *Login) Login(req *httprouter.Request) {
	var account string
	if account = req.FormValue("account"); account == "" {
		this.String("没有制定账号", 500)
		return
	}
	user := model.NewUser()
	user.Repo().Where("name", account).Or().
		WhereRaw("email_addr is not null").Where("email_addr", account).Or().
		WhereRaw("phone_number is not null").Where("phone_number", account)
	if m, exist, err := user.Repo().One(); err != nil {
		this.InternalError(err)
		return
	} else if !exist {
		this.String("用户名或密码不正确", 500)
		return
	} else {
		user = m.(*model.User)
	}
	if user.Encrypt(req.FormValue("password")) != user.Password {
		this.String("用户名或密码不正确", 500)
		return
	}
	if id, err := session.Save(user); err != nil {
		this.InternalError(err)
	} else {
		profile := user.Profile()
		profile["sId"] = id
		this.Json(profile, 200)
	}
}

func (this *Login) Logout(req *httprouter.Request) {
	session.Del(req.Header.Get("id"))
}

func (this *Login) SendVerifyEmail(user *model.User) error {
	return nil
}
