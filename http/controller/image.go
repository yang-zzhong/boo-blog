package controller

import (
	"boo-blog/model"
	"database/sql"
	"github.com/nfnt/resize"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	m "github.com/yang-zzhong/go-model"
	"io"
	"log"
	. "net/http"
	"os"
	"strconv"
)

type Image struct{ *Controller }

func (this *Image) Find(req *httprouter.Request) {
	userImage := model.NewUserImage()
	if userId := req.FormValue("user_id"); userId != "" {
		userImage.Repo().Where("user_id", userId)
	}
	if cateId := req.FormValue("cate_id"); cateId != "" {
		userImage.Repo().Where("cate_id", cateId)
	} else {
		userImage.Repo().Where("cate_id", 0)
	}
	if page := req.FormInt("page"); page != 0 {
		pageSize := req.FormInt("page_size")
		if pageSize == 0 {
			pageSize = 20
		}
		userImage.Repo().Page(int(page), int(pageSize))
	}
	if userImages, err := userImage.Repo().Fetch(); err != nil {
		this.InternalError(err)
	} else {
		var result []map[string]interface{}
		for _, entity := range userImages {
			image := entity.(*model.UserImage)
			result = append(result, map[string]interface{}{
				"image_id":      image.Hash,
				"user_image_id": image.Id,
				"user_id":       image.UserId,
			})
		}
		this.Json(result, 200)
	}
}

func (this *Image) Create(req *httprouter.Request, p *helpers.P) {
	src, header, err := req.FormFile("image")
	if err != nil {
		this.InternalError(err)
		return
	}
	defer src.Close()
	image := model.NewImage().Instance()
	err = image.FillWithMultipart(src, header)
	if err != nil {
		this.InternalError(err)
		return
	}
	if !image.FileExisted() {
		dist, err := os.OpenFile(image.Pathfile(), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			this.InternalError(err)
			return
		}
		defer dist.Close()
		src.Seek(0, 0)
		io.Copy(dist, src)
	}
	var exists bool
	if exists, err = image.RecordExisted(); err != nil {
		this.InternalError(err)
		return
	}
	if !exists {
		if err = image.Create(); err != nil {
			this.InternalError(err)
			return
		}
	}
	userImage := model.NewUserImage().Instance()
	userImage.Repo().Where("user_id", p.Get("visitor_id")).Where("hash", image.Hash)
	if entity, exist, err := userImage.Repo().One(); err != nil {
		this.InternalError(err)
		return
	} else if exist {
		userImage = entity.(*model.UserImage)
	}
	userImage.UserId = p.Get("visitor_id").(uint32)
	userImage.Hash = image.Hash
	if err = userImage.Save(); err != nil {
		this.InternalError(err)
		return
	}
	this.Json(map[string]interface{}{
		"user_image_id": userImage.Id,
		"name":          image.Name,
		"image_id":      image.Hash,
	}, 200)
}

func (this *Image) Get(req *httprouter.Request, p *helpers.P) {
	image := model.NewImage()
	if m, exist, err := image.Repo().Find(p.Get("id")); err != nil {
		this.InternalError(err)
		return
	} else if !exist {
		this.String("图片不存在", 404)
		return
	} else {
		image = m.(*model.Image)
	}
	err := image.Resize(
		this.ResponseWriter(),
		(uint)(req.FormUint("w")),
		(uint)(req.FormUint("h")),
		resize.NearestNeighbor,
	)
	if err != nil {
		this.InternalError(err)
		return
	}
	this.ResponseWriter().Header().Set("Content-Type", image.MimeType())
	this.ResponseWriter().WriteHeader(StatusOK)
	return
}

func (image *Image) MoveTo(req *httprouter.Request, p *helpers.P) {
	cate := model.NewCate()
	userImage := model.NewUserImage()
	ids := req.FormSlice("image_ids")
	log.Print("image_ids", req.FormValue("image_ids"))
	log.Print("ids", ids)
	if len(ids) == 0 {
		return
	}
	image_ids := []interface{}{}
	for _, id := range ids {
		image_ids = append(image_ids, id)
	}
	if entity, ok, err := cate.Repo().Find(p.Get("cate_id")); err != nil {
		image.InternalError(err)
		return
	} else if !ok {
		image.String("分类不存在", 500)
		return
	} else {
		cate = entity.(*model.Cate)
	}
	if cate.UserId != p.Get("visitor_id") {
		image.String("你不能移动图片到别人的分类", 500)
		return
	}
	userImage.Repo().
		WhereIn("id", image_ids).
		Where("user_id", p.Get("visitor_id"))
	if r, err := userImage.Repo().Fetch(); err != nil {
		image.InternalError(err)
		return
	} else {
		log.Print(r)
		err := m.Conn.Tx(func(tx *sql.Tx) error {
			for _, entity := range r {
				userImage = entity.(*model.UserImage)
				cateId, _ := strconv.ParseUint(p.Get("cate_id").(string), 10, 32)
				userImage.CateId = uint32(cateId)
				if err := userImage.Repo().WithTx(tx).Update(userImage); err != nil {
					return err
				}
			}
			return nil
		}, nil, nil)
		if err != nil {
			image.InternalError(err)
		}
	}
}
