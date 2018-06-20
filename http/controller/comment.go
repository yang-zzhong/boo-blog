package controller

import (
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-querybuilder"
	"time"
)

type Comment struct{ *Controller }

func (this *Comment) Create(req *httprouter.Request, p *helpers.P) {
	comment := model.NewComment().NewInstance()
	comment.Fill(map[string]interface{}{
		"user_id":      p.Get("visitor_id"),
		"content":      req.FormValue("content"),
		"blog_id":      req.FormValue("blog_id"),
		"ats":          this.parseAts(req.FormValue("content")),
		"comment_id":   req.FormValue("comment_id"),
		"commented_at": time.Now(),
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
			result = append(result, map[string]interface{}{
				"user_id":      c.UserId,
				"content":      c.Display(),
				"reply":        c.GetOne("reply").(*model.User).Name,
				"user_name":    c.GetOne("user").(*model.User).Name,
				"commented_at": c.CommentedAt,
			})
		}
		this.Json(result, 200)
	}
}
