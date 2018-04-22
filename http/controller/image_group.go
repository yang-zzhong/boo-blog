package controller

import (
	"boo-blog/model"
	. "net/http"
)

type ImageGroup struct{ Controller }

func (ig *ImageGroup) Create(req *Request) {
	mig := model.NewImageGroup()
	mig.Name = req.FormValue("name")
	mig.Intro = req.FormValue("intro")
	mig.TagIds = req.FormValue("tag_ids")
	mig.UserId = "36c122e0bf536f739e28a006f8b995c1"
	repo, err := model.NewImageGroupRepo()
	if err != nil {
		ig.InternalError(err)
		return
	}
	if err := repo.Create(mig); err != nil {
		ig.InternalError(err)
		return
	}
}
