package controller

import (
	"boo-blog/model"
	"database/sql"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	m "github.com/yang-zzhong/go-model"
	. "github.com/yang-zzhong/go-querybuilder"
)

type Article struct{ *Controller }

func (this *Article) Find(req *httprouter.Request, p *helpers.P) {
	blog := model.NewBlog()
	if ownerId := req.FormValue("owner_id"); ownerId != "" {
		blog.Repo().Where("user_id", ownerId)
	}
	if cateId := req.FormValue("cate_id"); cateId != "" {
		blog.Repo().Where("cate_id", cateId)
	}
	if tag := req.FormValue("tag"); tag != "" {
		blog.Repo().WhereRaw("JSON_SEARCH(tags, 'one', " + tag + ") IS NOT NULL")
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
	if p.Get("visitor_id") != nil {
		blog.Repo().WithCustom("thumb_up", func(mo interface{}) (m.NexusValues, error) {
			return thumbedCallback(mo, p.Get("visitor_id"))
		})
		blog.Repo().WithCustom("thumb_down", func(mo interface{}) (m.NexusValues, error) {
			return thumbedCallback(mo, p.Get("visitor_id"))
		})
	}
	if req.FormValue("with-author") == "1" {
		blog.Repo().With("author")
	}
	if items, err := blog.Repo().Fetch(); err != nil {
		this.InternalError(err)
		return
	} else {
		result := []map[string]interface{}{}
		for _, i := range items {
			item := i.(*model.Blog).Map()
			if p.Get("visitor_id") != nil {
				if thumbedUp, err := i.(*model.Blog).Many("thumb_up"); err != nil {
					this.InternalError(err)
					return
				} else {
					item["thumbed_up"] = thumbedUp
				}
				if thumbedDown, err := i.(*model.Blog).Many("thumb_down"); err != nil {
					this.InternalError(err)
					return
				} else {
					item["thumbed_down"] = thumbedDown
				}
			} else {
				item["thumbed_up"] = false
				item["thumbed_down"] = false
			}
			if req.FormValue("with-author") != "1" {
				result = append(result, item)
				continue
			}
			if author, err := i.(*model.Blog).One("author"); err != nil {
				this.InternalError(err)
				return
			} else if author == nil {
				this.String("系统错误", 500)
				return
			} else {
				item["author"] = map[string]interface{}{
					"id":                author.(*model.User).Id,
					"name":              author.(*model.User).Name,
					"portrait_image_id": author.(*model.User).PortraitImageId,
				}
				result = append(result, item)
			}
		}
		this.Json(result, 200)
	}
}

func (this *Article) GetOne(req *httprouter.Request, p *helpers.P) {
	blog := model.NewBlog()
	if i, ok, err := blog.Repo().Find(p.Get("id")); err != nil {
		this.InternalError(err)
	} else if !ok {
		this.String("文章没找到", 404)
	} else {
		blog = i.(*model.Blog)
		if d, err := detail(blog, p.Get("visitor_id")); err != nil {
			this.InternalError(err)
		} else {
			d["content"] = blog.Content()
			this.Json(d, 200)
		}
	}

}

func (this *Article) FetchUserBlog(req *httprouter.Request, p *helpers.P) {
	blog := model.NewBlog()
	blog.Repo().
		Where("user_id", p.Get("user_id")).
		Where("url_id", p.Get("url_id"))
	if i, ok, err := blog.Repo().One(); err != nil {
		this.InternalError(err)
	} else if ok {
		blog = i.(*model.Blog)
		if d, err := detail(blog, p.Get("visitor_id")); err != nil {
			this.InternalError(err)
			return
		} else {
			d["content"] = blog.Content()
			this.Json(d, 200)
		}
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
	if i, ok, err := blog.Repo().Find(p.Get("id")); err != nil {
		this.InternalError(err)
		return
	} else if !ok {
		this.String("文章未找到", 404)
		return
	} else {
		blog = i.(*model.Blog)
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
	if i, ok, err := blog.Repo().Find(p.Get("id")); err != nil {
		this.InternalError(err)
		return
	} else if ok {
		this.String("文章未找到", 404)
		return
	} else {
		blog = i.(*model.Blog)
	}
	if blog.UserId != p.Get("visitor_id").(uint32) {
		this.String("你没有权限修改别人的文章", 500)
		return
	}
	if err := blog.Delete(); err != nil {
		this.InternalError(err)
	}
}

func detail(blog *model.Blog, visitorId interface{}) (result map[string]interface{}, err error) {
	result = blog.Map()
	if visitorId == nil {
		result["thumbed_up"] = false
		result["thumbed_down"] = false
		return
	}
	vote := model.NewVote()
	vote.Repo().Where("target_id", blog.Id).
		Where("user_id", visitorId)
	if i, exists, e := vote.Repo().One(); e != nil {
		err = e
		return
	} else if !exists {
		result["thumbed_up"] = false
		result["thumbed_down"] = false
	} else if i.(*model.Vote).Vote > 0 {
		result["thumbed_up"] = true
		result["thumbed_down"] = false
	} else {
		result["thumbed_up"] = false
		result["thumbed_down"] = true
	}

	return
}

type withCustomBlogIn struct {
	blogIds []uint32
}

func (wcv *withCustomBlogIn) DataOf(mo interface{}, _ m.Nexus) interface{} {
	for _, blogId := range wcv.blogIds {
		if mo.(*model.Blog).Id == blogId {
			return true
		}
	}
	return false
}

func thumbedCallback(mo interface{}, visitorId interface{}) (val m.NexusValues, err error) {
	repo := mo.(m.Model).Repo()
	repo.Where("user_id", visitorId)
	repo.GroupBy("target_id")
	repo.Select("target_id")
	var blogIds []uint32
	err = repo.Query(func(rows *sql.Rows, _ []string) error {
		var blogId uint32
		if err := rows.Scan(&blogId); err != nil {
			return err
		}
		blogIds = append(blogIds, blogId)
		return nil
	})
	if err == nil {
		val = &withCustomBlogIn{blogIds}
	}
	return
}
