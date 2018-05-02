package controller

import (
	"boo-blog/model"
	"context"
	"database/sql"
	"errors"
	"github.com/nfnt/resize"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-model"
	"io"
	. "net/http"
	"os"
)

type Image struct{ *Controller }

func (this *Image) Find(req *httprouter.Request) {
	var userImageRepo *Repo
	var err error
	var userImages map[string]interface{}
	if userImageRepo, err = model.NewUserImageRepo(); err != nil {
		this.InternalError(err)
		return
	}
	if userId := req.FormValue("user_id"); userId != "" {
		userImageRepo.Where("user_id", userId)
	}
	if groupId := req.FormValue("group_id"); groupId != "" {
		userImageRepo.Where("group_id", groupId)
	}
	if page := req.FormInt("page"); page != 0 {
		pageSize := req.FormInt("page_size")
		if pageSize == 0 {
			pageSize = 20
		}
		userImageRepo.Offset((int)((page - 1) * pageSize)).Limit((int)(pageSize))
	}
	if userImages, err = userImageRepo.Fetch(); err != nil {
		this.InternalError(err)
		return
	}
	var result []map[string]interface{}
	for _, m := range userImages {
		image := m.(model.UserImage)
		result = append(result, map[string]interface{}{
			"image_id": image.Hash,
			"user_id":  image.UserId,
		})
	}
	this.Json(result, 200)
}

func (this *Image) Create(req *httprouter.Request, p *helpers.P) {
	src, header, err := req.FormFile("image")
	if err != nil {
		this.InternalError(err)
		return
	}
	defer src.Close()
	image := model.NewImage()
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
		var imageRepo *Repo
		if imageRepo, err = model.NewImageRepo(); err != nil {
			this.InternalError(err)
			return
		}
		if err = imageRepo.Create(image); err != nil {
			this.InternalError(err)
			return
		}
	}
	var userImageRepo *Repo
	if userImageRepo, err = model.NewUserImageRepo(); err != nil {
		this.InternalError(err)
		return
	}
	userImageRepo.Where("user_id", p.Get("visitor_id")).Where("hash", image.Hash)
	var userImage *model.UserImage
	if m := userImageRepo.One(); m != nil {
		ui := m.(model.UserImage)
		userImage = &ui
	} else {
		userImage = model.NewUserImage()
		userImage.UserId = p.Get("visitor_id").(string)
		userImage.Hash = image.Hash
		if err = userImageRepo.Create(userImage); err != nil {
			this.InternalError(err)
			return
		}
	}
	this.Json(map[string]string{
		"user_image_id": userImage.Id,
		"name":          image.Name,
		"image_id":      image.Hash,
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

func (image *Image) MoveTo(req *httprouter.Request, p *helpers.P) {
	var imageRepo *Repo
	var groupRepo *Repo
	var image_ids []string
	var err error
	image_ids = req.FormSlice("image_ids")
	if len(image_ids) == 0 {
		return
	}
	if groupRepo, err = model.NewCategoryRepo(); err != nil {
		image.InternalError(err)
		return
	}
	var m interface{}
	if m = groupRepo.Find(req.FormValue("group_id")); m == nil {
		image.InternalError(err)
		return
	}
	if group := m.(*model.Category); group.UserId != p.Get("visitor").(*model.User).Id {
		image.String("你没有权限移动图片到其他人的分类", 405)
		return
	}
	if imageRepo, err = model.NewUserImageRepo(); err != nil {
		image.InternalError(err)
		return
	}
	imageRepo.WhereIn("id", image_ids)
	images, _ := imageRepo.Fetch()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = imageRepo.Tx(func(tx *sql.Tx) error {
		for _, item := range images {
			image := item.(model.UserImage)
			if image.UserId != p.Get("visitor").(*model.User).Id {
				return errors.New("你没有权限移动别人的图片")
			}
			image.GroupId = req.FormValue("group_id")
			if err = imageRepo.WithTx(tx).Update(image); err != nil {
				return err
			}
		}
		return nil
	}, ctx, nil)
	if err != nil {
		image.InternalError(err)
	}
}
