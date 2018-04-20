package controller

import (
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	_ "image/png"
	"log"
	. "net/http"
)

type Image struct{}

func (image *Image) Create(w ResponseWriter, req *Request) {
	src, header, err := req.FormFile("image")
	if err != nil {
		log.Fatal(err)
		Error(w, err.Error(), 500)
		return
	}
	mImage := model.NewImage()
	err = mImage.FillWithMultipart(src, header)
	if err != nil {
		log.Fatal(err)
		Error(w, err.Error(), 500)
		return
	}
	exists, err := mImage.RecordExisted()
	if err != nil {
		log.Fatal(err)
		Error(w, err.Error(), 500)
		return
	}
	repo, err := model.NewImageRepo()
	if err != nil {
		Error(w, err.Error(), 500)
		return
	}
	if !exists {
		if err = repo.Create(mImage); err != nil {
			log.Fatal(err)
			Error(w, err.Error(), 500)
			return
		}
	}
	exists, err = mImage.FileExisted()
	if err != nil {
		Error(w, err.Error(), 500)
		return
	}
	if !exists {
		if err = mImage.SaveFile(src); err != nil {
			Error(w, err.Error(), 500)
			return
		}
	}
}

func (image *Image) Move(w ResponseWriter, req *Request, p *helpers.P) {
	repo, err := model.NewImageRepo()
	if err != nil {
		log.Fatal(err)
		Error(w, err.Error(), 500)
		return
	}
	mImage := repo.Find(p.Get("id"))
	if mImage == nil {
		Error(w, "资源不存在", 404)
		return
	}
	m := mImage.(model.Image)
	m.GroupId = req.FormValue("group_id")

	if err = repo.Update(m); err != nil {
		Error(w, err.Error(), 500)
	}
}
