package controller

import (
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-querybuilder"
)

type Theme struct{ *Controller }

func (this *Theme) Create(req *httprouter.Request, p *helpers.P) {
	theme := model.NewTheme().Instance()
	theme.Fill(map[string]interface{}{
		"name":    req.FormValue("name"),
		"content": req.FormMap("content"),
	})
	theme.Repo().Where("name", theme.Name).Where("user_id", p.Get("visitor_id"))
	if exists, err := theme.Count(); err != nil {
		this.InternalError(err)
	} else if exists {
		this.String("该主题已存在", 500)
	}
	if err := theme.Save(); err != nil {
		this.InternalError(err)
	}
}

func (this *Theme) Update(req *httprouter.Request, p *helpers.P) {
	theme := model.NewTheme()
	if m, ok, err := theme.Repo().Find(p.Get("theme_id")); err != nil {
		this.InternalError(err)
	} else if !ok {
		this.String("主题未找到", 500)
	} else {
		theme = m.(*model.Theme)
	}
	if theme.UserId != p.Get("visitor_id") {
		this.String("你没有权限修改别人的主题", 403)
		return
	}
	theme.Fill(map[string]interface{}{
		"name":    req.FormValue("name"),
		"content": req.FormMap("content"),
	})
	theme.Repo().
		Where("name", theme.Name).
		Where("user_id", p.Get("visitor_id")).
		Where("id", NEQ, theme.Id)
	if exists, err := theme.Count(); err != nil {
		this.InternalError(err)
	} else if exists {
		this.String("该主题已存在", 500)
	}
	if err := theme.Save(); err != nil {
		this.InternalError(err)
	}
}

func (this *Theme) Delete(p *helpers.P) {
	theme := model.NewTheme()
	if m, ok, err := theme.Repo().Find(p.Get("theme_id")); err != nil {
		this.InternalError(err)
	} else if !ok {
		this.String("主题未找到", 500)
	} else {
		theme = m.(*model.Theme)
	}
	if theme.UserId != p.Get("visitor_id") {
		this.String("你没有权限修改别人的主题", 403)
		return
	}
	if err := theme.Delete(); err != nil {
		this.InternalError(err)
	}
}
