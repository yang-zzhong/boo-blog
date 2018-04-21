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
	repo, err := model.NewImageRepo()
	if err != nil {
		image.InternalError(err)
		return
	}
	if !exists {
		if err = repo.Create(mImage); err != nil {
			image.InternalError(err)
			return
		}
	}
	if mImage.FileExisted() {
		return
	}
	dist, err := os.OpenFile(mImage.Pathfile(), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		image.InternalError(err)
		return
	}
	src.Seek(0, 0)
	defer dist.Close()
	io.Copy(dist, src)
}

func (image *Image) Get(req *Request, p *helpers.P) {
	repo, err := model.NewImageRepo()
	if err != nil {
		image.InternalError(err)
		return
	}
	mImage := repo.Find(p.Get("id")).(model.Image)
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
