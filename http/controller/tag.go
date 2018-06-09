package controller

import (
	"boo-blog/model"
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
	tag := model.NewTag()
	tag.Repo().Where("title", name)
	if m, exist, err := tag.Repo().One(); err != nil {
		this.InternalError(err)
	} else if exist {
		this.Json(m, 200)
		return
	}
	tag.Name = name
	tag.Intro = req.FormValue("intro")
	tag.IntroUrl = req.FormValue("intro_url")
	tag.UserId = p.Get("visitor_id").(string)
	if err := tag.Save(); err != nil {
		this.InternalError(err)
		return
	}
	this.Json(tag, 200)
}

func (this *Tag) ArticleUsed(req *httprouter.Request, p *helpers.P) {
	blog := model.NewBlog()
	blog.Repo().Where("user_id", p.Get("user_id"))
	blog.Repo().Select("id", "tags")
	if models, err := blog.Repo().Fetch(); err != nil {
		this.InternalError(err)
		return
	} else {
		tags := make(map[string]string)
		var result []map[string]string
		for _, m := range models {
			for _, tag := range m.(*model.Blog).Tags {
				tags[tag] = tag
			}
		}
		for _, tag := range tags {
			result = append(result, map[string]string{"name": tag})
		}

		this.Json(result, 200)
	}
}

func (this *Tag) Search(req *httprouter.Request) {
	tag := model.NewTag()
	if keyword := req.FormValue("keyword"); keyword == "" {
		this.Json([]string{}, 200)
		return
	} else {
		tag.Repo().Where("title", LIKE, keyword+"%").Limit(10)
	}
	this.renderRepo(tag.Repo())
}

func (this *Tag) renderRepo(repo *Repo) {
	if models, err := repo.Fetch(); err != nil {
		this.InternalError(err)
	} else {
		this.Json(models, 200)
	}
}
