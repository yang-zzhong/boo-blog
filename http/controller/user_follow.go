package controller

import (
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	"time"
)

type UserFollow struct{ *Controller }

func (this *UserFollow) Follow(p *helpers.P) {
	userFollow := model.NewUserFollow()
	userFollow.Repo().
		Where("user_id", p.Get("user_id")).
		Where("followed_by", p.Get("visitor_id"))
	if count, err := userFollow.Repo().Count(); err != nil {
		this.InternalError(err)
		return
	} else if count > 0 {
		this.String("你已经关注过该用户了", 500)
		return
	}
	userFollow.Instance()
	userFollow.Fill(map[string]interface{}{
		"user_id":     p.Get("user_id"),
		"followed_by": p.Get("visitor_id"),
		"followed_at": time.Now(),
	})

	if err := userFollow.Save(); err != nil {
		this.InternalError(err)
	}
}

func (this *UserFollow) Unfollow(p *helpers.P) {
	userFollow := model.NewUserFollow()
	userFollow.Repo().
		Where("user_id", p.Get("user_id")).
		Where("followed_by", p.Get("visitor_id"))

	if m, exists, err := userFollow.Repo().One(); err != nil {
		this.InternalError(err)
		return
	} else if !exists {
		this.String("你没有关注该用户", 500)
		return
	} else {
		userFollow = m.(*model.UserFollow)
	}
	if err := userFollow.Delete(); err != nil {
		this.InternalError(err)
	}
}
