package controller

import (
	"boo-blog/model"
	"github.com/nfnt/resize"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	"io"
	. "net/http"
	"os"
)

type Image struct{ *Controller }

func (this *Image) Find(req *httprouter.Request) {
	userImage := model.NewUserImage()
	if userId := req.FormValue("user_id"); userId != "" {
		userImage.Repo().Where("user_id", userId)
	}
	if groupId := req.FormValue("group_id"); groupId != "" {
		userImage.Repo().Where("group_id", groupId)
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
		for _, m := range userImages {
			image := m.(*model.UserImage)
			result = append(result, map[string]interface{}{
				"image_id": image.Hash,
				"user_id":  image.UserId,
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
	if m, exist, err := userImage.Repo().One(); err != nil {
		this.InternalError(err)
		return
	} else if exist {
		userImage = m.(*model.UserImage)

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
	// var imageRepo *Repo
	// var groupRepo *Repo
	// var image_ids []string
	// var err error
	// image_ids = req.FormSlice("image_ids")
	// if len(image_ids) == 0 {
	// 	return
	// }
	// if groupRepo, err = model.NewCategoryRepo(); err != nil {
	// 	image.InternalError(err)
	// 	return
	// }
	// var m interface{}
	// if m = groupRepo.Find(req.FormValue("group_id")); m == nil {
	// 	image.InternalError(err)
	// 	return
	// }
	// if group := m.(*model.Category); group.UserId != p.Get("visitor").(*model.User).Id {
	// 	image.String("你没有权限移动图片到其他人的分类", 405)
	// 	return
	// }
	// if imageRepo, err = model.NewUserImageRepo(); err != nil {
	// 	image.InternalError(err)
	// 	return
	// }
	// imageRepo.WhereIn("id", image_ids)
	// images, _ := imageRepo.Fetch()
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	// err = imageRepo.Tx(func(tx *sql.Tx) error {
	// 	for _, item := range images {
	// 		image := item.(model.UserImage)
	// 		if image.UserId != p.Get("visitor").(*model.User).Id {
	// 			return errors.New("你没有权限移动别人的图片")
	// 		}
	// 		image.GroupId = req.FormValue("group_id")
	// 		if err = imageRepo.WithTx(tx).Update(image); err != nil {
	// 			return err
	// 		}
	// 	}
	// 	return nil
	// }, ctx, nil)
	// if err != nil {
	// 	image.InternalError(err)
	// }
}
