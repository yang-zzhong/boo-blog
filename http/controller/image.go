package controller

import (
	"boo-blog/model"
	"github.com/nfnt/resize"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-model"
	"io"
	. "net/http"
	"os"
)

type Image struct{ *Controller }

func (image *Image) Create(req *httprouter.Request) {
	src, header, err := req.FormFile("image")
	if err != nil {
		image.InternalError(err)
		return
	}
	defer src.Close()
	mImage := model.NewImage()
	err = mImage.FillWithMultipart(src, header)
	if err != nil {
		image.InternalError(err)
		return
	}
	var exists bool
	if exists, err := mImage.RecordExisted(); err != nil {
		image.InternalError(err)
		return
	}
	if !exists {
		var repo *Repo
		if repo, err = model.NewImageRepo(); err != nil {
			image.InternalError(err)
			return
		}
		if err = repo.Create(mImage); err != nil {
			image.InternalError(err)
			return
		}
	}
	if !mImage.FileExisted() {
		dist, err := os.OpenFile(mImage.Pathfile(), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			image.InternalError(err)
			return
		}
		defer dist.Close()
		src.Seek(0, 0)
		io.Copy(dist, src)
	}
	image.Json(map[string]string{
		"name": mImage.Name,
		"id":   mImage.Hash,
	}, 200)
}

func (image *Image) Get(req *httprouter.Request, p *helpers.P) {
	repo, err := model.NewImageRepo()
	if err != nil {
		image.InternalError(err)
		return
	}
	mi := repo.Find(p.Get("id"))
	if mi == nil {
		image.String("图片不存在", 404)
		return
	}
	mImage := mi.(model.Image)
	err = mImage.Resize(
		image.ResponseWriter(),
		(uint)(req.FormUint("w")),
		(uint)(req.FormUint("h")),
		resize.NearestNeighbor,
	)
	if err != nil {
		image.InternalError(err)
		return
	}
	image.ResponseWriter().Header().Set("Content-Type", mImage.MimeType())
	image.ResponseWriter().WriteHeader(StatusOK)
	return
}
