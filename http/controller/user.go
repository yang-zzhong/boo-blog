package controller

import (
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-querybuilder"
)

type User struct{ *Controller }

func (this *User) Find(req *httprouter.Request) {
	user := model.NewUser()
	if keyword := req.FormValue("keyword"); keyword != "" {
		user.Repo().Quote(func(repo *Builder) {
			repo.Where("name", LIKE, "%"+keyword+"%").
				Or().Where("phone_number", LIKE, "%"+keyword+"%").
				Or().Where("email_addr", LIKE, "%"+keyword+"%")
		})
	}
	user.Repo().OrderBy("created_at", ASC)
	user.Repo().Limit(10)
	user.Repo().With("current_theme")
	if ms, err := user.Repo().Fetch(); err != nil {
		this.InternalError(err)
	} else {
		result := []map[string]interface{}{}
		for _, m := range ms {
			result = append(result, m.(*model.User).Profile())
		}
		this.Json(result, 200)
	}
}

func (this *User) Profile(p *helpers.P) {
	user := model.NewUser()
	if m, exists, err := user.Repo().Find(p.Get("user_id")); err != nil {
		this.InternalError(err)
		return
	} else if !exists {
		this.String("用户未找到", 404)
		return
	} else {
		this.Json(m.(*model.User).Profile(), 200)
	}
}

func (this *User) One(req *httprouter.Request, p *helpers.P) {
	user := model.NewUser()
	user.Repo().With("current_theme").Where("name", p.Get("name"))
	if m, exist, err := user.Repo().One(); err != nil {
		this.InternalError(err)
		return
	} else if !exist {
		this.String("没有找到用户", 404)
		return
	} else {
		user = m.(*model.User)
		result := user.Map()
		if theme, err := m.(*model.User).One("current_theme"); err != nil {
			this.InternalError(err)
			return
		} else if theme == nil {
			this.Json(result, 200)
		} else {
			result["theme"] = theme.(*model.Theme).Content
			this.Json(result, 200)
		}
	}
}

func (this *User) SaveUserInfo(req *httprouter.Request, p *helpers.P) {
	user := model.NewUser()
	if m, exist, err := user.Repo().Find(p.Get("visitor_id")); err != nil {
		this.InternalError(err)
	} else if !exist {
		this.String("用户不存在", 500)
	} else {
		user = m.(*model.User)
	}
	if req.FormValue("bio") != "" {
		user.Bio = req.FormValue("bio")
	}
	if req.FormValue("portrait_image_id") != "" {
		user.PortraitImageId = req.FormValue("portrait_image_id")
	}
	if req.FormValue("blog_name") != "" {
		user.BlogName = req.FormValue("blog_name")
	}
	if req.FormValue("name") != "" {
		user.Repo().Where("name", user.Name).Where("id", NEQ, p.Get("visitor_id"))
		if exists, err := user.Repo().Count(); err != nil {
			this.InternalError(err)
			return
		} else if exists > 0 {
			this.String("名字已存在", 500)
			return
		} else {
			user.Name = req.FormValue("name")
		}
	}
	if err := user.Save(); err != nil {
		this.InternalError(err)
	}
}
