package controller

import (
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-model"
)

type User struct{ *Controller }

func (this *User) One(req *httprouter.Request, p *helpers.P) {
	var repo *Repo
	var err error
	var user model.User
	if repo, err = model.NewUserRepo(); err != nil {
		this.InternalError(err)
		return
	}
	repo.Where("name", p.Get("name"))
	if m := repo.One(); m != nil {
		user = m.(model.User)
	} else {
		this.String("没有找到用户", 404)
		return
	}
	this.Json(map[string]interface{}{
		"id":                user.Id,
		"name":              user.Name,
		"nickname":          user.NickName,
		"portrait_image_id": user.PortraitImageId,
		"small_image_id":    user.SimpleBgImageId,
		"bg_image_id":       user.BigBgImageId,
	}, 200)
}
