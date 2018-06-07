package model

import (
	model "github.com/yang-zzhong/go-model"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type Blog struct {
	Id        string    `db:"id char(32) pk"`
	Title     string    `db:"title varchar(256)"`
	Overview  string    `db:"overview text"`
	UrlId     string    `db:"url_id varchar(256)"`
	UserId    string    `db:"user_id char(32)"`
	CateId    string    `db:"cate_id char(32) nil"`
	Tags      []string  `db:"tags varchar(256) nil"`
	CreatedAt time.Time `db:"created_at datetime"`
	UpdatedAt time.Time `db:"updated_at datetime"`
	*model.Base
}

func (blog *Blog) TableName() string {
	return "article"
}

func (blog *Blog) SaveContent(content string) error {
	log.Print(blog.Pathfile())
	return ioutil.WriteFile(blog.Pathfile(), []byte(content), 0755)
}

func (blog *Blog) Pathfile() string {
	return conf.blog_dir + blog.UserId + "-" + blog.Title + ".html"
}

func (blog *Blog) WithUrlId() *Blog {
	blog.UrlId = blog.GetUrlId(blog.Title)
	return blog
}

func (blog *Blog) GetUrlId(title string) string {
	reg, _ := regexp.Compile("\\s|\\?|\\&|\"|'")
	return reg.ReplaceAllString(title, "_")
}

func (blog *Blog) WithOverview(content string) {
	reader := strings.NewReader(content)
	node, _ := html.Parse(reader)
	type callback func(n *html.Node) bool
	var find func(n *html.Node, call callback) bool
	find = func(n *html.Node, call callback) bool {
		if call(n) {
			return true
		}
		for c := n.FirstChild; c != nil; c = c.FirstChild {
			if find(c, call) {
				return true
			}
		}
		for c := n.NextSibling; c != nil; c = c.NextSibling {
			if find(c, call) {
				return true
			}
		}
		return false
	}
	overview := ""
	find(node, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "div" {
			find(n, func(d *html.Node) bool {
				if d.Type == html.TextNode {
					overview += d.Data
				}
				if len(overview) > 512 {
					return true
				}
				return false
			})
			return true
		}
		return false
	})
	blog.Overview = overview
}

func (blog *Blog) Content() string {
	content, _ := ioutil.ReadFile(blog.Pathfile())
	return string(content)
}

func (blog *Blog) DBValue(colname string, value interface{}) interface{} {
	if colname == "tags" {
		result := strings.Join(value.([]string), ",")
		return result
	}
	return value
}

func (blog *Blog) Value(colname string, value interface{}) (result reflect.Value, catched bool) {
	if colname == "tags" {
		catched = true
		val, _ := value.(string)
		if val != "" {
			result = reflect.ValueOf(strings.Split(val, ","))
			return
		}
		result = reflect.ValueOf([]string{})
		return
	}
	catched = false
	return
}

func NewBlog() *Blog {
	blog := model.NewModel(new(Blog)).(*Blog)
	blog.DeclareOne("author", new(User), map[string]string{
		"user_id": "id",
	})
	blog.DeclareOne("cate", new(Cate), map[string]string{
		"cate_id": "id",
	})

	return blog
}

func (blog *Blog) Instance() *Blog {
	return Instance(blog).(*Blog)
}
