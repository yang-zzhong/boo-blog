package controller

import (
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-model"
	. "github.com/yang-zzhong/go-querybuilder"
)

type Category struct{ *Controller }

func (controller *Category) Get(req *httprouter.Request, p *helpers.P) {

}

func (controller *Category) Create(req *httprouter.Request, p *helpers.P) {
	cate := model.NewCategory()
	cate.Name = req.FormValue("name")
	cate.Intro = req.FormValue("intro")
	cate.TagIds = req.FormSlice("tag_ids")
	cate.UserId = p.Get("visitor_id").(string)
	var repo *Repo
	var err error
	if repo, err = model.NewCategoryRepo(); err != nil {
		controller.InternalError(err)
		return
	}
	repo.Where("name", cate.Name).Where("user_id", cate.UserId)
	if repo.Count() > 0 {
		controller.String("名字已存在", 402)
		return
	}
	if err = repo.Create(cate); err != nil {
		controller.InternalError(err)
	}
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
	cate.TagIds = req.FormSlice("tag_ids")
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
