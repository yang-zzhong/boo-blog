package controller

import (
	"boo-blog/model"
	"database/sql"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-model"
	. "github.com/yang-zzhong/go-querybuilder"
	"log"
)

type Category struct{ *Controller }

func (this *Category) Find(req *httprouter.Request) {
	var repo *Repo
	var err error
	if repo, err = model.NewCategoryRepo(); err != nil {
		this.InternalError(err)
		return
	}
	if userId := req.FormValue("user_id"); userId != "" {
		repo.Where("user_id", userId)
	}
	this.renderCates(repo, make(map[string]int))
}

func (controller *Category) Create(req *httprouter.Request, p *helpers.P) {
	cate := model.NewCategory()
	cate.Name = req.FormValue("name")
	cate.Intro = req.FormValue("intro")
	cate.Tags = req.FormSlice("tags")
	cate.UserId = p.Get("visitor_id").(string)
	var repo *Repo
	var err error
	if repo, err = model.NewCategoryRepo(); err != nil {
		controller.InternalError(err)
		return
	}
	repo.Where("name", cate.Name).Where("user_id", cate.UserId)
	if repo.Count() > 0 {
		controller.String("名字已存在", 500)
		return
	}
	if err = repo.Create(cate); err != nil {
		controller.InternalError(err)
	}
}

func (this *Category) ImageUsed(req *httprouter.Request, p *helpers.P) {
	var repo, imageRepo *Repo
	var err error
	if imageRepo, err = model.NewImageRepo(); err != nil {
		this.InternalError(err)
		return
	}
	imageRepo.Where("user_id", p.Get("user_id")).
		Select("cate_id", E{"count(1) as quantity"}).
		GroupBy("cate_id")
	var cateIds []string
	idQuantity := make(map[string]int)
	imageRepo.QueryCallback(func(rows *sql.Rows) {
		var id string
		var quantity int
		if err = rows.Scan(&id, &quantity); err != nil {
			return
		}
		cateIds = append(cateIds, id)
		idQuantity[id] = quantity
	})
	if err != nil {
		this.InternalError(err)
		return
	}
	if repo, err = model.NewCategoryRepo(); err != nil {
		this.InternalError(err)
		return
	}
	if len(cateIds) == 0 {
		return
	}
	repo.WhereIn("id", cateIds)

	this.renderCates(repo, idQuantity)
}

func (this *Category) ArticleUsed(req *httprouter.Request, p *helpers.P) {
	var repo, articleRepo *Repo
	var err error
	if articleRepo, err = model.NewArticleRepo(); err != nil {
		this.InternalError(err)
		return
	}
	articleRepo.Where("user_id", p.Get("user_id")).
		Select("cate_id", E{"count(1) as quantity"}).
		GroupBy("cate_id")
	var cateIds []string
	idQuantity := make(map[string]int)
	articleRepo.QueryCallback(func(rows *sql.Rows) {
		var id string
		var quantity int
		if err = rows.Scan(&id, &quantity); err != nil {
			return
		}
		cateIds = append(cateIds, id)
		idQuantity[id] = quantity
	})
	if repo, err = model.NewCategoryRepo(); err != nil {
		this.InternalError(err)
		return
	}
	if len(cateIds) == 0 {
		return
	}
	log.Print("cateIds", cateIds)
	repo.WhereIn("id", cateIds)
	log.Print(repo.ForQuery())

	this.renderCates(repo, idQuantity)
}

func (controller *Category) Update(req *httprouter.Request, p *helpers.P) {
	var repo *Repo
	var err error
	if repo, err = model.NewCategoryRepo(); err != nil {
		controller.InternalError(err)
		return
	}
	repo.Where("name", req.FormValue("name")).
		Where("user_id", p.Get("visitor_id")).
		Where("id", NEQ, p.Get("id"))
	if repo.Count() > 0 {
		controller.String("名字已存在", 402)
	}
	var cate *model.Category
	if m := repo.Find(p.Get("id")); m != nil {
		c := m.(model.Category)
		cate = &c
	} else {
		controller.String("分类未找到", 404)
		return
	}
	if cate.UserId != p.Get("visitor_id") {
		controller.String("你没有权限修改别人的分类", 405)
		return
	}
	cate.Name = req.FormValue("name")
	cate.Intro = req.FormValue("intro")
	cate.Tags = req.FormSlice("tags")
	if err = repo.Update(cate); err != nil {
		controller.InternalError(err)
	}
}

func (controller *Category) Delete(req *httprouter.Request, p *helpers.P) {
	var repo *Repo
	var err error
	if repo, err = model.NewCategoryRepo(); err != nil {
		controller.InternalError(err)
		return
	}
	var cate *model.Category
	if m := repo.Find(p.Get("id")); m == nil {
		controller.String("分类未找到", 404)
		return
	} else {
		c := m.(model.Category)
		cate = &c
	}
	if cate.UserId != p.Get("visitor_id") {
		controller.String("你没有权限修改别人的分类", 405)
		return
	}
	repo.Where("id", cate.Id)
	if err = repo.Remove(); err != nil {
		controller.InternalError(err)
	}
}

func (this *Category) renderCates(repo *Repo, idQuantity map[string]int) {
	var data map[string]interface{}
	var err error
	var result []map[string]interface{}
	repo.OrderBy("created_at", DESC)
	if data, err = repo.Fetch(); err != nil {
		this.InternalError(err)
		return
	}
	for _, item := range data {
		cate := item.(model.Category)
		var quantity int
		var ok bool
		if quantity, ok = idQuantity[cate.Id]; !ok {
			quantity = 0
		}
		result = append(result, map[string]interface{}{
			"id":       cate.Id,
			"tags":     cate.Tags,
			"name":     cate.Name,
			"intro":    cate.Intro,
			"quantity": quantity,
		})
	}
	this.Json(result, 200)
}
