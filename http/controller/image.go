package controller

import (
	"boo-blog/model"
	"github.com/nfnt/resize"
	helpers "github.com/yang-zzhong/go-helpers"
	"io"
	. "net/http"
	"os"
	"strconv"
)

type Image struct{ Controller }

func (image *Image) Create(req *Request) {
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
	exists, err := mImage.RecordExisted()
	if err != nil {
		image.InternalError(err)
		return
	}
	if !exists {
		repo, err := model.NewImageRepo()
		if err != nil {
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

func (image *Image) Get(req *Request, p *helpers.P) {
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
	var width, height int
	w := req.FormValue("w")
	if w != "" {
		width, err = strconv.Atoi(w)
	}
	h := req.FormValue("h")
	if h != "" {
		height, err = strconv.Atoi(h)
	}
	if err := mImage.Resize(image.Writer(), (uint)(width), (uint)(height), resize.NearestNeighbor); err != nil {
		image.InternalError(err)
		return
	}
	image.Writer().Header().Set("Content-Type", mImage.MimeType())
	image.Writer().WriteHeader(StatusOK)
	return
}
