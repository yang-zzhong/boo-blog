package controller

import (
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-model"
)

type User struct{ *Controller }

func (this *User) One(req *httprouter.Request, p *helpers.P) {
	var repo *Repo
	var err error
	var user model.User
	var theme model.Theme
	if repo, err = model.NewUserRepo(); err != nil {
		this.InternalError(err)
		return
	}
	repo.Where("name", p.Get("name"))
	if m := repo.One(); m != nil {
		user = m.(model.User)
	} else {
		this.String("没有找到用户", 404)
		return
	}
	if repo, err = model.NewThemeRepo(); err != nil {
		this.InternalError(err)
		return
	}
	if m := repo.Find(user.Id); m != nil {
		theme = m.(model.Theme)
	} else {
		this.String("服务器出错咯", 500)
	}
	result := map[string]interface{}{
		"id":               user.Id,
		"name":             user.Name,
		"nickname":         "",
		"bg_image_id":      "",
		"info_bg_image_id": "",
		"bg_color":         "",
		"info_bg_color":    "",
		"name_color":       "",
		"blog_name":        theme.Name,
	}
	if user.NickName.Valid {
		result["nickname"] = user.NickName.String
	}
	if theme.BgImageId.Valid {
		result["bg_image_id"] = theme.BgImageId.String
	}
	if theme.InfoBgImageId.Valid {
		result["info_bg_image_id"] = theme.InfoBgImageId.String
	}
	if theme.BgColor.Valid {
		result["bg_color"] = theme.BgColor.String
	}
	if theme.InfoBgColor.Valid {
		result["info_bg_color"] = theme.InfoBgColor.String
	}
	if theme.NameColor.Valid {
		result["name_color"] = theme.NameColor.String
	}
	this.Json(result, 200)
}

func (this *User) SaveBlogInfo(req *httprouter.Request, p *helpers.P) {
	var repo *Repo
	var err error
	if repo, err = model.NewUserRepo(); err != nil {
		this.InternalError(err)
		return
	}
	repo.Fetch()
}
