package controller

import (
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-model"
)

type Category struct{ *Controller }

func (cate *Category) Create(req *httprouter.Request, p *helpers.P) {
	mcate := model.NewImageGroup()
	mcate.Name = req.FormValue("name")
	mcate.Intro = req.FormValue("intro")
	mcate.TagIds = req.FormSlice("tag_ids")
	mcate.UserId = p.Get("visitor").(*model.User).Id
	var repo *Repo
	var err error
	if repo, err = model.NewCategoryRepo(); err != nil {
		cate.InternalError(err)
		return
	}
	if err = repo.Create(mcate); err != nil {
		cate.InternalError(err)
	}
}

func (controller *Category) Update(req *httprouter.Request, p *helpers.P) {
	var repo *Repo
	var err error
	if repo, err = model.NewCategoryRepo(); err != nil {
		controller.InternalError(err)
		return
	}
	m := repo.Find(p.Get("id"))
	if m == nil {
		controller.String("分类未找到", 404)
		return
	}
	cate := m.(*model.Category)
	if cate.UserId != p.Get("visitor").(*model.User).Id {
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
	if repo, err = model.NewCategoryRepo; err != nil {
		controller.InternalError(err)
		return
	}
	m := repo.Find(p.Get("id"))
	if m == nil {
		controller.String("分类未找到", 404)
		return
	}
	cate := m.(*model.Catetory)
	if cate.UserId != p.Get("visitor").(*model.User).Id {
		controller.String("你没有权限修改别人的分类", 405)
		return
	}
	if err = repo.Remove(cate); err != nil {
		controller.InternalError(err)
	}
}
