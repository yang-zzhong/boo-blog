package model

import (
	"github.com/google/uuid"
	model "github.com/yang-zzhong/go-model"
	. "github.com/yang-zzhong/go-querybuilder"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Blog struct {
	Id           uint32    `db:"id bigint pk"`
	Title        string    `db:"title varchar(256)"`
	Overview     string    `db:"overview text"`
	Image        string    `db:"image varchar(1024)"`
	UrlId        string    `db:"url_id varchar(256)"`
	UserId       uint32    `db:"user_id bigint"`
	CateId       uint32    `db:"cate_id bigint nil"`
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

func (blog *Blog) SaveContent(content string) error {
	return ioutil.WriteFile(blog.Pathfile(), []byte(content), 0755)
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
	reg, _ := regexp.Compile("\\s|\\?|\\&|\"|'")
	return reg.ReplaceAllString(title, "-")
}

func (blog *Blog) WithOverview(content string) {
	reader := strings.NewReader(content)
	node, _ := html.Parse(reader)
	nodes := 0
	find(node, func(d *html.Node) bool {
		if nodes > 300 {
			return true
		}
		nodes++
		if d.Type == html.ElementNode && d.Data == "img" {
			for _, attr := range d.Attr {
				if attr.Key == "src" {
					blog.Image = attr.Val
					return true
				}
			}
		}
		return false
	})
	limit := 256
	if blog.Image != "" {
		limit = 128
	}
	nodes = 0
	find(node, func(d *html.Node) bool {
		if nodes > 300 {
			return true
		}
		nodes++
		if d.Type == html.TextNode && inblock(d) {
			blog.Overview += strings.Trim(d.Data, " ")
		}
		if len(blog.Overview) > limit {
			return true
		}
		return false
	})
	if len(blog.Overview) > limit {
		blog.Overview = blog.Overview[0:limit]
	}
	log.Printf("%v, %d", blog.Overview, len(blog.Overview))
}

func (blog *Blog) Content() string {
	content, _ := ioutil.ReadFile(blog.Pathfile())
	return string(content)
}

func NewBlog() *Blog {
	blog := model.NewModel(new(Blog)).(*Blog)
	blog.DeclareOne("author", new(User), model.Nexus{
		"id": "user_id",
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

type callback func(n *html.Node) bool

func find(n *html.Node, call callback) bool {
	if call(n) {
		return true
	}
	var w, d bool
	for c := n.FirstChild; c != nil; c = c.FirstChild {
		d = find(c, call)
	}
	for c := n.NextSibling; c != nil; c = c.NextSibling {
		w = find(c, call)
	}
	return w || d
}

func inblock(n *html.Node) bool {
	node := n
	for node != nil {
		if node.Type == html.ElementNode && (node.Data == "style" || node.Data == "head" || node.Data == "code") {
			return false
		}
		node = node.Parent
	}

	return true
}
