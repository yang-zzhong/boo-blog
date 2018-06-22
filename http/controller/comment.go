package controller

import (
	"boo-blog/model"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-querybuilder"
	"strconv"
)

type Comment struct{ *Controller }

func (this *Comment) Create(req *httprouter.Request, p *helpers.P) {
	comment := model.NewComment().Instance()
	blogId, _ := strconv.ParseUint(p.Get("blog_id").(string), 10, 32)
	if blogId == 0 {
		this.String("需要上传blog id", 400)
		return
	}
	data := map[string]interface{}{
		"user_id":    p.Get("visitor_id"),
		"content":    req.FormValue("content"),
		"blog_id":    uint32(blogId),
		"ats":        this.parseAts(req.FormValue("content")),
		"comment_id": uint32(req.FormUint("comment_id")),
	}
	comment.Fill(data)
	if err := comment.Save(); err != nil {
		this.InternalError(err)
	}
}

func (this *Comment) Articles(p *helpers.P) {
	comment := model.NewComment()
	// blogId, _ := strconv.ParseUint(p.Get("blog_id").(string), 10, 32)
	comment.Repo().Where("blog_id", p.Get("blog_id"))
	comment.Repo().With("user").With("reply")
	comment.Repo().OrderBy("commented_at", DESC)
	if data, err := comment.Repo().Fetch(); err != nil {
		this.InternalError(err)
	} else {
		var result []map[string]interface{}
		for _, item := range data {
			c := item.(*model.Comment)
			replyName := ""
			userName := ""
			portraitImageId := ""
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
				portraitImageId = user.(*model.User).PortraitImageId
			}
			result = append(result, map[string]interface{}{
				"comment_id":        c.Id,
				"user_id":           c.UserId,
				"portrait_image_id": portraitImageId,
				"content":           c.Display(),
				"reply":             replyName,
				"user_name":         userName,
				"commented_at":      c.CommentedAt,
			})
		}
		this.Json(result, 200)
	}
}

func (this *Comment) parseAts(content string) []string {
	return []string{}
}
