package controller

import (
	"boo-blog/model"
	"context"
	"database/sql"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	m "github.com/yang-zzhong/go-model"
	"strconv"
	"time"
)

type Vote struct{ *Controller }

func (this *Vote) Create(req *httprouter.Request, p *helpers.P) {
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
	blogId, _ := strconv.ParseUint(p.Get("blog_id").(string), 10, 32)
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
	if req.FormInt("vote") > 0 {
		blog.ThumbUp = blog.ThumbUp + 1
		vote.Vote = 1
	} else {
		blog.ThumbDown = blog.ThumbDown + 1
		vote.Vote = -1
	}
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
