package controller

import (
	"boo-blog/model"
	"database/sql"
	helpers "github.com/yang-zzhong/go-helpers"
	m "github.com/yang-zzhong/go-model"
	"log"
	"strconv"
	"time"
)

type UserFollow struct{ *Controller }

func (this *UserFollow) Follow(p *helpers.P) {
	userFollow := model.NewUserFollow()
	log.Print("user_id", p.Get("user_id"))
	id, _ := strconv.ParseUint(p.Get("user_id").(string), 10, 32)
	userId := uint32(id)
	log.Print("userId", userId)
	userFollow.Repo().
		Where("user_id", p.Get("visitor_id")).
		Where("followed", userId)
	if count, err := userFollow.Repo().Count(); err != nil {
		this.InternalError(err)
		return
	} else if count > 0 {
		this.String("你已经关注过该用户了", 500)
		return
	}
	userFollow.Instance()
	userFollow.Fill(map[string]interface{}{
		"user_id":     p.Get("visitor_id"),
		"followed":    userId,
		"followed_at": time.Now(),
	})
	var following, followed *model.User
	if i, err := userFollow.One("following"); err != nil {
		this.InternalError(err)
		return
	} else if i == nil {
		this.String("系统错误", 500)
		return
	} else {
		following = i.(*model.User)
		following.Followed = following.Followed + 1
		followed = p.Get("visitor").(*model.User)
		followed.Following = followed.Following + 1
	}
	err := m.Conn.Tx(func(tx *sql.Tx) error {
		if err := userFollow.Repo().WithTx(tx).Create(userFollow); err != nil {
			return err
		}
		log.Print("following", following)
		if err := following.Repo().WithTx(tx).Update(following); err != nil {
			return err
		}
		log.Print("followed", followed)
		if err := followed.Repo().WithTx(tx).Update(followed); err != nil {
			return err
		}
		return nil
	}, nil, nil)
	if err != nil {
		this.InternalError(err)
	}
}

func (this *UserFollow) Unfollow(p *helpers.P) {
	userFollow := model.NewUserFollow()
	id, _ := strconv.ParseUint(p.Get("user_id").(string), 10, 32)
	userId := uint32(id)
	userFollow.Repo().
		Where("user_id", p.Get("visitor_id")).
		Where("followed", userId)
	var following, followed *model.User
	if i, exists, err := userFollow.Repo().One(); err != nil {
		this.InternalError(err)
		return
	} else if !exists {
		this.String("你没有关注该用户", 500)
		return
	} else {
		userFollow = i.(*model.UserFollow)
	}
	if i, err := userFollow.One("following"); err != nil {
		this.InternalError(err)
		return
	} else if i == nil {
		this.String("系统错误", 500)
		return
	} else {
		following = i.(*model.User)
		following.Followed = following.Followed - 1
		followed = p.Get("visitor").(*model.User)
		followed.Following = followed.Following - 1
	}
	err := m.Conn.Tx(func(tx *sql.Tx) error {
		if err := userFollow.Repo().WithTx(tx).Delete(userFollow); err != nil {
			return err
		}
		if err := following.Repo().WithTx(tx).Update(following); err != nil {
			return err
		}
		if err := followed.Repo().WithTx(tx).Update(followed); err != nil {
			return err
		}
		return nil
	}, nil, nil)
	if err != nil {
		this.InternalError(err)
	}
}
