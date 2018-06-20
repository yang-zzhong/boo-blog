package controller

import (
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-querybuilder"
)

type Comment struct{ *Controller }

func (this *Comment) Create(req *httprouter.Request, p *helpers.P) {
	comment := model.NewComment().Instance()
	comment.Fill(map[string]interface{}{
		"user_id":    p.Get("visitor_id"),
		"content":    req.FormValue("content"),
		"blog_id":    req.FormValue("blog_id"),
		"ats":        this.parseAts(req.FormValue("content")),
		"comment_id": req.FormValue("comment_id"),
	})

	if err := comment.Save(); err != nil {
		this.InternalError(err)
	}
}

func (this *Comment) Articles(req *httprouter.Request) {
	comment := model.NewComment()
	comment.Repo().Where("blog_id", req.FormValue("blog_id"))
	comment.Repo().With("user_id").With("reply")
	comment.Repo().OrderBy("commented_at", DESC)
	if data, err := comment.Repo().Fetch(); err != nil {
		this.InternalError(err)
	} else {
		var result []map[string]interface{}
		for _, item := range data {
			c := item.(*model.Comment)
			replyName := ""
			userName := ""
			if reply, err := c.One("reply"); err != nil {
				this.InternalError(err)
				return
			} else if reply != nil {
				replyName = reply.(*model.User).Name
			}
			if user, err := c.One("user"); err != nil {
				this.InternalError(err)
				return
			} else if user != nil {
				userName = user.(*model.User).Name
			}
			result = append(result, map[string]interface{}{
				"user_id":      c.UserId,
				"content":      c.Display(),
				"reply":        replyName,
				"user_name":    userName,
				"commented_at": c.CommentedAt,
			})
		}
		this.Json(result, 200)
	}
}

func (this *Comment) parseAts(content string) []uint32 {
	return []uint32{}
}
