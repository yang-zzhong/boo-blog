package controller

import (
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-querybuilder"
)

type Article struct{ *Controller }

func (this *Article) Find(req *httprouter.Request) {
	blog := model.NewBlog()
	if ownerId := req.FormValue("owner_id"); ownerId != "" {
		blog.Repo().Where("user_id", ownerId)
	}
	if cateId := req.FormValue("cate_id"); cateId != "" {
		blog.Repo().Where("cate_id", cateId)
	}
	if tag := req.FormValue("tag"); tag != "" {
		blog.Repo().Where("tags", LIKE, "%"+tag+"%")
	}
	if keyword := req.FormValue("keyword"); keyword != "" {
		blog.Repo().Quote(func(repo *Builder) {
			repo.Where("title", LIKE, "%"+keyword+"%").
				Or().Where("overview", LIKE, "%"+keyword+"%")
		})
	}
	if p := req.FormInt("page"); p != 0 {
		ps := req.FormInt("page_size")
		if ps == 0 {
			ps = 10
		}
		blog.Repo().Page(int(p), int(ps))
	}
	blog.Repo().OrderBy("thumb_up", DESC).
		OrderBy("created_at", DESC).
		OrderBy("thumb_down", ASC).
		OrderBy("comments", DESC)
	if req.FormValue("with-author") == "1" {
		blog.Repo().With("author")
	}
	if items, err := blog.Repo().Fetch(); err != nil {
		this.InternalError(err)
		return
	} else {
		result := []map[string]interface{}{}
		for _, m := range items {
			item := m.(*model.Blog).Map()
			if req.FormValue("with-author") != "1" {
				result = append(result, item)
				continue
			}
			if author, err := m.(*model.Blog).One("author"); err != nil {
				this.InternalError(err)
				return
			} else if author == nil {
				this.String("系统错误", 500)
				return
			} else {
				item["author"] = author.(*model.User).Map()
				result = append(result, item)
			}
		}
		this.Json(result, 200)
	}
}

func (this *Article) GetOne(req *httprouter.Request, p *helpers.P) {
	blog := model.NewBlog()
	if m, ok, err := blog.Repo().Find(p.Get("id")); err != nil {
		this.InternalError(err)
	} else if !ok {
		this.String("文章没找到", 404)
	} else {
		this.Json(m, 200)
	}

}

func (this *Article) FetchUserBlog(req *httprouter.Request, p *helpers.P) {
	blog := model.NewBlog()
	blog.Repo().
		Where("user_id", p.Get("user_id")).
		Where("url_id", p.Get("url_id"))
	if m, ok, err := blog.Repo().One(); err != nil {
		this.InternalError(err)
	} else if ok {
		blog = m.(*model.Blog)
		data := blog.Map()
		data["content"] = blog.Content()
		this.Json(data, 200)
	} else {
		this.String("文章未找到", 404)
	}
}

func (this *Article) Create(req *httprouter.Request, p *helpers.P) {
	blog := model.NewBlog().Instance()
	blog.Fill(map[string]interface{}{
		"title":   req.FormValue("title"),
		"user_id": p.Get("visitor_id"),
		"tags":    req.FormSlice("tags"),
		"cate_id": uint32(req.FormInt("cate_id")),
	})
	blog.WithUrlId().WithOverview(req.FormValue("content"))
	blog.Repo().Where("user_id", blog.UserId).Quote(func(repo *Builder) {
		repo.Where("title", blog.Title)
		repo.Or().Where("url_id", blog.UrlId)
	})
	if count, err := blog.Repo().Count(); err != nil {
		this.InternalError(err)
		return
	} else if count > 0 {
		this.String("该标题已使用", 500)
		return
	}
	if len(blog.Tags) == 0 {
		this.String("至少选择一个标签", 500)
		return
	}
	if err := blog.SaveContent(req.FormValue("content")); err != nil {
		this.InternalError(err)
		return
	}
	if err := blog.Save(); err != nil {
		this.InternalError(err)
		return
	}
	this.Json(blog, 200)
}

func (this *Article) Update(req *httprouter.Request, p *helpers.P) {
	blog := model.NewBlog()
	if m, ok, err := blog.Repo().Find(p.Get("id")); err != nil {
		this.InternalError(err)
		return
	} else if !ok {
		this.String("文章未找到", 404)
		return
	} else {
		blog = m.(*model.Blog)
	}
	if blog.UserId != p.Get("visitor_id").(uint32) {
		this.String("你没有权限修改别人的文章", 500)
		return
	}
	blog.Fill(map[string]interface{}{
		"title":   req.FormValue("title"),
		"cate_id": req.FormValue("cate_id"),
		"tags":    req.FormSlice("tags"),
	})
	blog.WithUrlId().WithOverview(req.FormValue("content"))
	if err := blog.Save(); err != nil {
		this.InternalError(err)
	}
}

func (this *Article) Remove(req *httprouter.Request, p *helpers.P) {
	blog := model.NewBlog()
	if m, ok, err := blog.Repo().Find(p.Get("id")); err != nil {
		this.InternalError(err)
		return
	} else if ok {
		this.String("文章未找到", 404)
		return
	} else {
		blog = m.(*model.Blog)
	}
	if blog.UserId != p.Get("visitor_id").(uint32) {
		this.String("你没有权限修改别人的文章", 500)
		return
	}
	if err := blog.Delete(); err != nil {
		this.InternalError(err)
	}
}
