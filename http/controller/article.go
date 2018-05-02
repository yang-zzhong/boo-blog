package controller

import (
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-model"
	. "github.com/yang-zzhong/go-querybuilder"
)

type Article struct{ *Controller }

func (this *Article) Find(req *httprouter.Request) {
	var repo *Repo
	var err error
	var items map[string]interface{}
	var result []map[string]interface{}
	if repo, err = model.NewArticleRepo(); err != nil {
		this.InternalError(err)
		return
	}
	if ownerId := req.FormValue("owner_id"); ownerId != "" {
		repo.Where("user_id", ownerId)
	}
	if groupId := req.FormValue("group_id"); groupId != "" {
		repo.Where("group_id", groupId)
	}
	if tag := req.FormValue("tag"); tag != "" {
		repo.Where("tag_ids", LIKE, "%"+tag+"%")
	}
	if keyword := req.FormValue("keyword"); keyword != "" {
		repo.Where("title", LIKE, "%"+keyword+"%").Or().Where("content", LIKE, "%"+keyword+"%")
	}
	if page := req.FormInt("page"); page != 0 {
		pageSize := req.FormInt("page_size")
		if pageSize == 0 {
			pageSize = 10
		}
		repo.Offset((int)((page - 1) * pageSize)).Limit((int)(pageSize))
	}
	if items, err = repo.Fetch(); err != nil {
		this.InternalError(err)
		return
	}
	for _, item := range items {
		atl := item.(model.Article)
		result = append(result, map[string]interface{}{
			"id":         atl.Id,
			"title":      atl.Title,
			"group_id":   atl.GroupId,
			"user_id":    atl.UserId,
			"content":    atl.Content,
			"created_at": atl.CreatedAt,
			"updated_at": atl.UpdatedAt,
		})
	}
	this.Json(result, 200)
}

func (this *Article) GetOne(req *httprouter.Request, p *helpers.P) {
	var repo *Repo
	var err error
	var article model.Article
	if repo, err = model.NewArticleRepo(); err != nil {
		this.InternalError(err)
		return
	}
	if m := repo.Find(p.Get("id").(string)); m != nil {
		article = m.(model.Article)
	} else {
		this.String("文章未找到", 404)
		return
	}

	this.Json(map[string]interface{}{
		"id":        article.Id,
		"title":     article.Title,
		"content":   article.Content,
		"userId":    article.UserId,
		"groupId":   article.GroupId,
		"Tags":      article.TagIds,
		"CreatedAt": article.CreatedAt,
		"UpdatedAt": article.UpdatedAt,
	}, 200)
}

func (this *Article) Create(req *httprouter.Request, p *helpers.P) {
	var repo *Repo
	var err error
	if repo, err = model.NewArticleRepo(); err != nil {
		this.InternalError(err)
		return
	}
	article := model.NewArticle()
	article.Title = req.FormValue("title")
	article.Content = req.FormValue("content")
	article.UserId = p.Get("visitor_id").(string)
	article.GroupId = req.FormValue("group_id")
	article.TagIds = req.FormSlice("tag_ids")
	if err = repo.Create(article); err != nil {
		this.InternalError(err)
		return
	}
	this.Json(map[string]string{
		"id": article.Id,
	}, 200)
}

func (this *Article) Update(req *httprouter.Request, p *helpers.P) {
	var repo *Repo
	var err error
	var article model.Article
	if repo, err = model.NewArticleRepo(); err != nil {
		this.InternalError(err)
		return
	}
	if m := repo.Find(p.Get("id").(string)); m != nil {
		article = m.(model.Article)
	} else {
		this.String("文章未找到", 404)
		return
	}
	if article.UserId != p.Get("visitor_id").(string) {
		this.String("你没有权限修改别人的文章", 500)
		return
	}
	article.Title = req.FormValue("title")
	article.Content = req.FormValue("content")
	article.GroupId = req.FormValue("group_id")
	article.TagIds = req.FormSlice("tag_ids")
	if err = repo.Update(&article); err != nil {
		this.InternalError(err)
	}
}

func (this *Article) Remove(req *httprouter.Request, p *helpers.P) {
	var repo *Repo
	var err error
	var article model.Article
	if repo, err = model.NewArticleRepo(); err != nil {
		this.InternalError(err)
		return
	}
	if m := repo.Find(p.Get("id").(string)); m != nil {
		article = m.(model.Article)
	} else {
		this.String("文章未找到", 404)
		return
	}
	if article.UserId != p.Get("visitor_id").(string) {
		this.String("你没有权限修改别人的文章", 500)
		return
	}
	repo.Where("id", article.Id)
	if err = repo.Remove(); err != nil {
		this.InternalError(err)
	}
}
