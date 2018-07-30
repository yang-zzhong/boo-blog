package controller

import (
	"boo-blog/model"
	"context"
	"database/sql"
	helpers "github.com/yang-zzhong/go-helpers"
	m "github.com/yang-zzhong/go-model"
	"strconv"
	"time"
)

type Vote struct{ *Controller }

func (this *Vote) Create(p *helpers.P) {
	blog := model.NewBlog()
	if m, ok, err := blog.Repo().Find(p.Get("blog_id")); err != nil {
		this.InternalError(err)
		return
	} else if !ok {
		this.String("文章不存在", 404)
		return
	} else {
		blog = m.(*model.Blog)
	}
	vote := model.NewVote().Instance()
	id, _ := strconv.ParseUint(p.Get("blog_id").(string), 10, 32)
	blogId := uint32(id)
	vote.Repo().Where("target_type", model.VOTE_BLOG).
		Where("target_id", blogId).
		Where("user_id", p.Get("visitor_id"))
	voted, err := vote.Repo().Count()
	if err != nil {
		this.InternalError(err)
		return
	}
	if voted > 0 {
		this.String("你已经投过票咯", 500)
		return
	}
	vote.Fill(map[string]interface{}{
		"target_type": model.VOTE_BLOG,
		"target_id":   blogId,
		"user_id":     p.Get("visitor_id"),
		"voted_at":    time.Now(),
	})
	if p.Get("vote").(int) > 0 {
		blog.ThumbUp = blog.ThumbUp + 1
	} else {
		blog.ThumbDown = blog.ThumbDown + 1
	}
	vote.Vote = int8(p.Get("vote").(int))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = m.Conn.Tx(func(tx *sql.Tx) error {
		if err := blog.Repo().WithTx(tx).Update(blog); err != nil {
			return err
		}
		if err := vote.Repo().WithTx(tx).Create(vote); err != nil {
			return err
		}
		return nil
	}, ctx, nil)
	if err != nil {
		this.InternalError(err)
	}
}

func (this *Vote) Delete(p *helpers.P) {
	id, _ := strconv.ParseUint(p.Get("blog_id").(string), 10, 32)
	blogId := uint32(id)
	vote := model.NewVote()
	vote.Repo().Where("target_type", model.VOTE_BLOG).
		Where("target_id", blogId).
		Where("user_id", p.Get("visitor_id"))
	if i, exists, err := vote.Repo().One(); err != nil {
		this.InternalError(err)
		return
	} else if !exists {
		this.String("你没有对该文章投过票", 500)
	} else {
		vote = i.(*model.Vote)
	}
	if i, err := vote.One("blog"); err != nil {
		this.InternalError(err)
		return
	} else if i == nil {
		if err := vote.Delete(); err != nil {
			this.InternalError(err)
			return
		}
	} else {
		blog := i.(*model.Blog)
		if vote.Vote > 0 {
			blog.ThumbUp = blog.ThumbUp - 1
		} else {
			blog.ThumbDown = blog.ThumbDown - 1
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err = m.Conn.Tx(func(tx *sql.Tx) error {
			if err := blog.Repo().WithTx(tx).Update(blog); err != nil {
				return err
			}
			if err := vote.Repo().WithTx(tx).Delete(vote); err != nil {
				return err
			}
			return nil
		}, ctx, nil)
		if err != nil {
			this.InternalError(err)
		}
	}
}
