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
	if cateId := req.FormValue("cate_id"); cateId != "" {
		repo.Where("cate_id", cateId)
	}
	if tag := req.FormValue("tag"); tag != "" {
		repo.Where("tags", LIKE, "%"+tag+"%")
	}
	if keyword := req.FormValue("keyword"); keyword != "" {
		repo.Where("title", LIKE, "%"+keyword+"%").
			Or().Where("content", LIKE, "%"+keyword+"%")
	}
	if page := req.FormInt("page"); page != 0 {
		pageSize := req.FormInt("page_size")
		if pageSize == 0 {
			pageSize = 10
		}
		repo.Page(page, pageSize)
	}
	repo.OrderBy("created_at", DESC)
	if items, err = repo.Fetch(); err != nil {
		this.InternalError(err)
		return
	}
	for _, item := range items {
		atl := item.(model.Article)
		result = append(result, map[string]interface{}{
			"id":        atl.Id,
			"title":     atl.Title,
			"url_id":    atl.UrlId,
			"overview":  atl.Overview,
			"tags":      atl.Tags,
			"CreatedAt": article.CreatedAt,
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
		"content":   (&article).Content(),
		"Tags":      article.Tags,
		"CreatedAt": article.CreatedAt,
	}, 200)
}

func (this *Article) FetchUserBlog(req *httprouter.Request, p *helpers.P) {
	var repo *Repo
	var err error
	var article model.Article
	if repo, err = model.NewArticleRepo(); err != nil {
		this.InternalError(err)
		return
	}
	repo.Where("user_id", p.Get("user_id")).Where("url_id", p.Get("url_id"))
	if m := repo.One(); m != nil {
		article = m.(model.Article)
	} else {
		this.String("文章未找到", 404)
		return
	}
	this.Json(map[string]interface{}{
		"id":        article.Id,
		"title":     article.Title,
		"overview":  article.Overview,
		"content":   (&article).Content(),
		"Tags":      article.Tags,
		"CreatedAt": article.CreatedAt,
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
	article.UserId = p.Get("visitor_id").(string)
	article.CateId = req.FormValue("cate_id")
	article.Tags = req.FormSlice("tags")
	article.WithUrlId().WithOverview(req.FormValue("content"))
	repo.Where("user_id", article.UserId).Quote(func(repo *Builder) {
		repo.Where("title", article.Title)
		repo.Or().Where("url_id", article.UrlId)
	})
	if repo.Count() > 0 {
		this.String("该标题已使用", 500)
		return
	}
	if len(article.Tags) == 0 {
		this.String("至少选择一个标签", 500)
		return
	}
	if err = article.SaveContent(req.FormValue("content")); err != nil {
		this.InternalError(err)
		return
	}
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
	if m := repo.Find(p.Get("id")); m != nil {
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
	article.CateId = req.FormValue("cate_id")
	article.Tags = req.FormSlice("tags")
	article.WithUrlId().WithOverview(req.FormValue("content"))
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
