package controller

import (
	"boo-blog/model"
	"github.com/nfnt/resize"
	helpers "github.com/yang-zzhong/go-helpers"
	"io"
	"log"
	. "net/http"
	"os"
	"strconv"
)

type Image struct{}

func (image *Image) Create(w ResponseWriter, req *Request) {
	src, header, err := req.FormFile("image")
	if err != nil {
		log.Fatal(err)
		Error(w, err.Error(), 500)
		return
	}
	defer src.Close()
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
	if exists = mImage.FileExisted(); !exists {
		dist, err := os.OpenFile(mImage.Pathfile(), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
			return
		}
		src.Seek(0, 0)
		defer dist.Close()
		io.Copy(dist, src)
	}
}

func (image *Image) Get(w ResponseWriter, req *Request, p *helpers.P) {
	repo, err := model.NewImageRepo()
	if err != nil {
		log.Fatal(err)
		Error(w, err.Error(), 500)
		return
	}
	mImage := repo.Find(p.Get("id")).(model.Image)
	var width, height int
	wi := req.FormValue("w")
	if wi != "" {
		width, err = strconv.Atoi(wi)
	}
	h := req.FormValue("h")
	if h != "" {
		height, err = strconv.Atoi(h)
	}
	if err := mImage.Resize(w, (uint)(width), (uint)(height), resize.NearestNeighbor); err != nil {
		log.Fatal(err)
		Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", mImage.MimeType())
	w.WriteHeader(StatusOK)
	return
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
