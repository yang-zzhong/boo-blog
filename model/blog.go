package model

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	model "github.com/yang-zzhong/go-model"
	. "github.com/yang-zzhong/go-querybuilder"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

type Blog struct {
	Id           uint32    `db:"id bigint pk"`
	Title        string    `db:"title varchar(256)"`
	Overview     string    `db:"overview text"`
	ImageHash    string    `db:"image_hash char(32) nil"`
	Image        string    `db:"image varchar(1024) nil"`
	UrlId        string    `db:"url_id varchar(256)"`
	UserId       uint32    `db:"user_id bigint"`
	CateId       uint32    `db:"cate_id bigint nil"`
	ContentId    uint32    `db:"content_id bigint"`
	Tags         []string  `db:"tags varchar(32)[] nil"`
	Comments     int       `db:"comments int"`
	ThumbUp      int       `db:"thumb_up int"`
	ThumbDown    int       `db:"thumb_down int"`
	AllowComment bool      `db:"allow_comment bool"`
	AllowThumb   bool      `db:"allow_thumb bool"`
	Privilege    string    `db:"privilege varchar(32)"`
	CreatedAt    time.Time `db:"created_at timestamp"`
	UpdatedAt    time.Time `db:"updated_at timestamp"`
	*model.Base
}

func (blog *Blog) TableName() string {
	return "blogs"
}

func (blog *Blog) DBValue(colname string, value interface{}) interface{} {
	if colname == "tags" {
		return nullArrayDBValue(value)
	}
	return value
}

func (blog *Blog) Value(colname string, value interface{}) (result reflect.Value, catch bool) {
	if colname == "tags" {
		catch = true
		result = nullArrayValue(value)
		return
	}
	catch = false
	return
}

func (blog *Blog) Pathfile() string {
	id := strconv.FormatUint(uint64(blog.UserId), 32)
	return conf.BlogDir + id + "-" + blog.Title + ".html"
}

func (blog *Blog) WithUrlId() *Blog {
	blog.UrlId = blog.GetUrlId(blog.Title)
	return blog
}

func (blog *Blog) GetUrlId(title string) string {
	reg, _ := regexp.Compile("\\s|\\?|\\&|\"|'|\\/")
	return reg.ReplaceAllString(title, "-")
}

func (blog *Blog) WithOverview() error {
	var content *BlogContent
	if m, err := blog.One("content"); err != nil {
		return err
	} else if m == nil {
		return errors.New("no content set")
	} else {
		content = m.(*BlogContent)
	}
	blog.ImageHash = content.PreviewImageHash()
	if blog.ImageHash == "" {
		blog.Image = content.PreviewImageUrl()
	} else {
		blog.Image = ""
	}
	limit := 128
	if blog.Image == "" && blog.ImageHash == "" {
		limit = 256
	}
	blog.Overview = content.Preview(limit)
	return nil
}

func (blog *Blog) Content() string {
	if m, err := blog.One("content"); err != nil || m == nil {
		return ""
	} else {
		return m.(*BlogContent).Content
	}
}

func (blog *Blog) Save() error {
	return model.Conn.Tx(func(tx *sql.Tx) error {
		blog.Repo().WithTx(tx)
		if err := blog.Base.Save(); err != nil {
			return nil
		}
		if m, err := blog.One("content"); err != nil {
			return err
		} else if m == nil {
			return errors.New("no content set")
		} else {
			m.(*BlogContent).Repo().WithTx(tx)
			m.(*BlogContent).Save()
		}
		return nil
	}, nil, nil)
}

func NewBlog() *Blog {
	blog := model.NewModel(new(Blog)).(*Blog)
	blog.DeclareOne("author", new(User), model.Nexus{
		"id": "user_id",
	})
	blog.DeclareOne("content", new(BlogContent), model.Nexus{
		"id": "content_id",
	})
	blog.DeclareOne("cate", new(Cate), model.Nexus{
		"id": "cate_id",
	})
	blog.DeclareMany("thumb_up", new(Vote), model.Nexus{
		"target_id":   "id",
		"target_type": model.NWhere{EQ, VOTE_BLOG},
		"vote":        model.NWhere{GT, 0},
	})
	blog.DeclareMany("thumb_down", new(Vote), model.Nexus{
		"target_id":   "id",
		"target_type": model.NWhere{EQ, VOTE_BLOG},
		"vote":        model.NWhere{LT, 0},
	})
	blog.OnUpdate(func(b interface{}) error {
		b.(*Blog).UpdatedAt = time.Now()
		return nil
	})

	return blog
}

func (blog *Blog) Instance() *Blog {
	blog.Id = uuid.New().ID()
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	return blog
}
