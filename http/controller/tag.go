package controller

import (
	"boo-blog/model"
	"database/sql"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-model"
	. "github.com/yang-zzhong/go-querybuilder"
)

type Tag struct{ *Controller }

func (this *Tag) Create(req *httprouter.Request, p *helpers.P) {
	name := req.FormValue("name")
	if name == "" {
		return
	}
	var repo *Repo
	var err error
	if repo, err = model.NewTagRepo(); err != nil {
		this.InternalError(err)
		return
	}
	repo.Where("name", name)
	if m := repo.One(); m != nil {
		tag := m.(model.Tag)
		intro, _ := tag.Intro.Value()
		introUrl, _ := tag.IntroUrl.Value()
		this.Json(map[string]interface{}{
			"name":      tag.Name,
			"intro":     intro,
			"intro_url": introUrl,
		}, 200)
		return
	}
	tag := model.NewTag()
	tag.Name = name
	tag.Intro = sql.NullString{req.FormValue("intro"), false}
	tag.IntroUrl = sql.NullString{req.FormValue("intro_url"), false}
	tag.UserId = p.Get("visitor_id").(string)
	if err = repo.Create(tag); err != nil {
		this.InternalError(err)
		return
	}
	this.Json(map[string]string{
		"name":      name,
		"intro":     req.FormValue("intro"),
		"intro_url": req.FormValue("intro_url"),
	}, 200)
}

func (this *Tag) ArticleUsed(req *httprouter.Request, p *helpers.P) {
	var repo *Repo
	var err error
	var models map[string]interface{}
	if repo, err = model.NewArticleRepo(); err != nil {
		this.InternalError(err)
		return
	}
	repo.Where("user_id", p.Get("user_id"))
	repo.Select("id", "tags")
	if models, err = repo.Fetch(); err != nil {
		this.InternalError(err)
		return
	}
	tags := make(map[string]string)
	var result []map[string]string
	for _, m := range models {
		atl := m.(model.Article)
		for _, tag := range atl.Tags {
			tags[tag] = tag
		}
	}
	for _, tag := range tags {
		result = append(result, map[string]string{"name": tag})
	}

	this.Json(result, 200)
}

func (this *Tag) Search(req *httprouter.Request) {
	var repo *Repo
	var err error
	var keyword string
	if repo, err = model.NewTagRepo(); err != nil {
		this.InternalError(err)
		return
	}
	if keyword = req.FormValue("keyword"); keyword == "" {
		this.Json([]string{}, 200)
		return
	}
	repo.Where("title", LIKE, keyword+"%").Limit(10)
	this.renderRepo(repo)
}

func (this *Tag) Get(req *httprouter.Request, _ *helpers.P) {
	var repo *Repo
	var err error
	if repo, err = model.NewTagRepo(); err != nil {
		this.InternalError(err)
		return
	}
	repo.OrderBy("created_at", DESC)
	this.renderRepo(repo)
}

func (this *Tag) renderRepo(repo *Repo) {
	var models map[string]interface{}
	var err error
	if models, err = repo.Fetch(); err != nil {
		this.InternalError(err)
		return
	}
	result := []map[string]interface{}{}
	for _, item := range models {
		tag := item.(model.Tag)
		intro, _ := tag.Intro.Value()
		introUrl, _ := tag.IntroUrl.Value()
		result = append(result, map[string]interface{}{
			"name":      tag.Name,
			"intro":     intro,
			"intro_url": introUrl,
		})
	}
	this.Json(result, 200)
}
