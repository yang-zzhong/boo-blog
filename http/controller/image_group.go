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
