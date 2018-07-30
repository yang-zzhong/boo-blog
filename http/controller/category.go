package controller

import (
	"boo-blog/model"
	"database/sql"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-model"
	. "github.com/yang-zzhong/go-querybuilder"
)

type Category struct{ *Controller }

func (this *Category) Find(req *httprouter.Request) {
	cate := model.NewCate()
	if userId := req.FormValue("user_id"); userId != "" {
		cate.Repo().Where("user_id", userId)
	}
	this.renderCates(cate.Repo(), make(map[uint32]int))
}

func (this *Category) Create(req *httprouter.Request, p *helpers.P) {
	cate := model.NewCate().Instance()
	cate.Fill(map[string]interface{}{
		"name":    req.FormValue("name"),
		"intro":   req.FormValue("intro"),
		"tags":    req.FormSlice("tags"),
		"user_id": p.Get("visitor_id"),
	})
	cate.Repo().Where("name", cate.Name).Where("user_id", cate.UserId)
	if count, err := cate.Repo().Count(); err != nil {
		this.InternalError(err)
		return
	} else if count > 0 {
		this.String("名字已存在", 500)
		return
	}
	if err := cate.Save(); err != nil {
		this.InternalError(err)
	}
}

func (this *Category) ImageUsed(req *httprouter.Request, p *helpers.P) {
	image := model.NewUserImage()
	cate := model.NewCate()
	image.Repo().Where("user_id", p.Get("user_id")).
		Select("cate_id", E{"count(1) as quantity"}).
		GroupBy("cate_id")
	var cateIds []interface{}
	idCount := make(map[uint32]int)
	err := image.Repo().Query(func(rows *sql.Rows, _ []string) error {
		var id uint32
		var quantity int
		if err := rows.Scan(&id, &quantity); err != nil {
			return err
		}
		cateIds = append(cateIds, id)
		idCount[id] = quantity
		return nil
	})
	if err != nil {
		this.InternalError(err)
		return
	}
	if len(cateIds) == 0 {
		return
	}
	cate.Repo().WhereIn("id", cateIds)

	this.renderCates(cate.Repo(), idCount)
}

func (this *Category) ArticleUsed(req *httprouter.Request, p *helpers.P) {
	blog := model.NewBlog()
	blog.Repo().Where("user_id", p.Get("user_id")).
		Select("cate_id", E{"count(1) as quantity"}).
		GroupBy("cate_id")
	var cateIds []interface{}
	idCount := make(map[uint32]int)
	err := blog.Repo().Query(func(rows *sql.Rows, _ []string) error {
		var id uint32
		var quantity int
		if err := rows.Scan(&id, &quantity); err != nil {
			return err
		}
		cateIds = append(cateIds, id)
		idCount[id] = quantity
		return nil
	})
	if err != nil {
		this.InternalError(err)
		return
	}
	if len(cateIds) == 0 {
		return
	}
	cate := model.NewCate()
	cate.Repo().WhereIn("id", cateIds)

	this.renderCates(cate.Repo(), idCount)
}

func (this *Category) Update(req *httprouter.Request, p *helpers.P) {
	cate := model.NewCate()
	cate.Repo().Where("name", req.FormValue("name")).
		Where("user_id", p.Get("visitor_id")).
		Where("id", NEQ, p.Get("id"))
	if count, err := cate.Repo().Count(); err != nil {
		this.InternalError(err)
		return
	} else if count > 0 {
		this.String("名字已存在", 402)
		return
	}
	if m, exist, err := cate.Repo().Find(p.Get("id")); err != nil {
		this.InternalError(err)
		return
	} else if !exist {
		this.String("分类未找到", 404)
		return
	} else {
		cate = m.(*model.Cate)
	}
	if cate.UserId != p.Get("visitor_id") {
		this.String("你没有权限修改别人的分类", 405)
		return
	}
	cate.Name = req.FormValue("name")
	cate.Intro = req.FormValue("intro")
	cate.Tags = req.FormSlice("tags")
	if err := cate.Save(); err != nil {
		this.InternalError(err)
	}
}

func (this *Category) Delete(req *httprouter.Request, p *helpers.P) {
	cate := model.NewCate()
	if m, exist, err := cate.Repo().Find(p.Get("id")); err != nil {
		this.InternalError(err)
		return
	} else if !exist {
		return
	} else {
		cate = m.(*model.Cate)
	}
	if cate.UserId != p.Get("visitor_id") {
		this.String("你没有权限修改别人的分类", 405)
		return
	}
	if err := cate.Delete(); err != nil {
		this.InternalError(err)
	}
}

func (this *Category) renderCates(repo *Repo, idQuantity map[uint32]int) {
	cate := model.NewCate()
	if data, err := cate.Repo().Fetch(); err != nil {
		this.InternalError(err)
	} else {
		var result []map[string]interface{}
		for _, item := range data {
			c := item.(*model.Cate).Map()
			if quantity, ok := idQuantity[c["id"].(uint32)]; ok {
				c["quantity"] = quantity
			} else {
				c["quantity"] = 0
			}
			result = append(result, c)
		}
		this.Json(result, 200)
	}
}
