package controller

import (
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
)

type User struct{ *Controller }

func (this *User) One(req *httprouter.Request, p *helpers.P) {
	user := model.NewUser()
	user.Repo().With("theme").Where("name", p.Get("name"))
	if m, exist, err := user.Repo().One(); err != nil {
		this.InternalError(err)
		return
	} else if !exist {
		this.String("没有找到用户", 404)
		return
	} else {
		user = m.(*model.User)
		result := user.Map()
		if theme, err := m.(*model.User).One("theme"); err != nil {
			this.InternalError(err)
			return
		} else {
			for f, v := range theme.(*model.Theme).Map() {
				if f == "name" {
					f = "blog_name"
				}
				result[f] = v
			}
			this.Json(result, 200)
		}
	}
}

func (this *User) SaveBlogInfo(req *httprouter.Request, p *helpers.P) {
	theme := model.NewTheme()
	if m, exist, err := theme.Repo().Find(p.Get("visitor_id")); err != nil {
		this.InternalError(err)
	} else if !exist {
		this.String("没有找到主题", 404)
	} else {
		theme = m.(*model.Theme)
		data := map[string]interface{}{
			"bg_image_id":   req.FormValue("bg_image_id"),
			"name":          req.FormValue("blog_name"),
			"bg_color":      req.FormValue("bg_color"),
			"info_bg_color": req.FormValue("info_bg_color"),
			"name_color":    req.FormValue("name_color"),
		}
		theme.Fill(data)
		if err := theme.Save(); err != nil {
			this.InternalError(err)
		}
	}
}
