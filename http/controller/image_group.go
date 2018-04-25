package controller

import (
	"boo-blog/model"
	httprouter "github.com/yang-zzhong/go-httprouter"
	"log"
)

type ImageGroup struct{ *Controller }

func (ig *ImageGroup) Create(req *httprouter.Request) {
	mig := model.NewImageGroup()
	mig.Name = req.FormValue("name")
	mig.Intro = req.FormValue("intro")
	log.Println(req.FormValue("tag_ids"))
	// mig.TagIds = req.FormValue("tag_ids")
	// mig.UserId = "36c122e0bf536f739e28a006f8b995c1"
	// repo, err := model.NewImageGroupRepo()
	// if err != nil {
	// 	ig.InternalError(err)
	// 	return
	// }
	// if err := repo.Create(mig); err != nil {
	// 	ig.InternalError(err)
	// 	return
	// }
}
