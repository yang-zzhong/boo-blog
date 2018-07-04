package controller

import (
	"boo-blog/model"
	"context"
	"database/sql"
	helpers "github.com/yang-zzhong/go-helpers"
	httprouter "github.com/yang-zzhong/go-httprouter"
	. "github.com/yang-zzhong/go-querybuilder"
	"strconv"
)

type Comment struct{ *Controller }

func (this *Comment) Create(req *httprouter.Request, p *helpers.P) {
	comment := model.NewComment().Instance()
	blogId, _ := strconv.ParseUint(p.Get("blog_id").(string), 10, 32)
	var blog *model.Blog
	if blogId == 0 {
		this.String("需要上传blog id", 400)
		return
	}
	if m, ok, err := model.NewBlog().Repo().Find(blogId); err != nil {
		this.InternalError(err)
		return
	} else if !ok {
		this.String("博客不存在", 400)
	} else {
		blog = m.(*model.Blog)
	}
	data := map[string]interface{}{
		"user_id":    p.Get("visitor_id"),
		"content":    req.FormValue("content"),
		"blog_id":    uint32(blogId),
		"ats":        this.parseAts(req.FormValue("content")),
		"comment_id": uint32(req.FormUint("comment_id")),
	}
	comment.Fill(data)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := model.NewComment().Repo().Tx(func(tx *sql.Tx) error {
		if err := comment.Repo().WithTx(tx).Create(comment); err != nil {
			return err
		}
		blog.Comments = blog.Comments + 1
		if err := model.NewBlog().Repo().WithTx(tx).Update(blog); err != nil {
			return err
		}
		return nil
	}, ctx, nil)
	if err != nil {
		this.InternalError(err)
	}
}

func (this *Comment) Articles(req *httprouter.Request, p *helpers.P) {
	comment := model.NewComment()
	comment.Repo().Where("blog_id", p.Get("blog_id"))
	comment.Repo().With("user").With("reply")
	comment.Repo().Page(int(req.FormInt("page")), 10)
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

func (this *Comment) Comments(p *helpers.P) {

}

func (this *Comment) Delete(p *helpers.P) {
	comment := model.NewComment()
	var blog *model.Blog
	if m, ok, err := comment.Repo().Find(p.Get("id")); err != nil {
		this.InternalError(err)
	} else if !ok {
		this.String("评论未找到", 404)
	} else {
		comment = m.(*model.Comment)
	}
	if comment.BlogId != 0 {
		if m, err := comment.One("blog"); err != nil {
			this.InternalError(err)
			return
		} else if m == nil {
			this.String("博客不存在", 404)
			return
		} else {
			blog = m.(*model.Blog)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := blog.Repo().Tx(func(tx *sql.Tx) error {
		if err := comment.Repo().WithTx(tx).Delete(comment); err != nil {
			return err
		}
		if blog == nil || blog != nil && blog.Comments == 0 {
			return nil
		}
		blog.Comments = blog.Comments - 1
		if err := blog.Repo().WithTx(tx).Update(blog); err != nil {
			return err
		}
		return nil
	}, ctx, nil)
	if err != nil {
		this.InternalError(err)
	}
}

func (this *Comment) parseAts(content string) []string {
	return []string{}
}
