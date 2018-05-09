package model

import (
	helpers "github.com/yang-zzhong/go-helpers"
	. "github.com/yang-zzhong/go-model"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type Article struct {
	Id        string    `db:"id char(32) pk"`
	Title     string    `db:"title varchar(64)"`
	Overview  string    `db:"overview text"`
	UrlId     string    `db:"url_id varchar(256)"`
	UserId    string    `db:"user_id char(32)"`
	CateId    string    `db:"cate_id char(32) nil"`
	Tags      []string  `db:"tags varchar(256) nil"`
	CreatedAt time.Time `db:"created_at datetime"`
	UpdatedAt time.Time `db:"updated_at datetime"`
}

func (atl *Article) PK() string {
	return "id"
}

func (atl *Article) SaveContent(content string) error {
	log.Print(atl.Pathfile())
	return ioutil.WriteFile(atl.Pathfile(), []byte(content), 0755)
}

func (atl *Article) Pathfile() string {
	return conf.blog_dir + atl.UserId + "-" + atl.Title + ".html"
}

func (atl *Article) WithUrlId() *Article {
	atl.UrlId = atl.GetUrlId(atl.Title)
	return atl
}

func (atl *Article) GetUrlId(title string) string {
	reg, _ := regexp.Compile("\\s|\\?|\\&")
	return reg.ReplaceAllString(title, "_")
}

func (atl *Article) WithOverview(content string) {
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
				if len(overview) > 256 {
					return true
				}
				return false
			})
			return true
		}
		return false
	})
	if len(overview) > 256 {
		overview = overview[0:255]
	}
	log.Print(overview)
	atl.Overview = overview
}

func (atl *Article) Content() string {
	content, _ := ioutil.ReadFile(atl.Pathfile())
	return string(content)
}

func (atl *Article) TableName() string {
	return "article"
}

func (atl *Article) NewId() interface{} {
	return helpers.RandString(32)
}

func (atl *Article) DBValue(fieldName string, value interface{}) interface{} {
	if fieldName == "tags" {
		result := strings.Join(value.([]string), ",")
		return result
	}
	return value
}

func (atl *Article) Value(fieldName string, value interface{}) (result reflect.Value, catched bool) {
	if fieldName == "tags" {
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

func NewArticle() *Article {
	return CreateModel(new(Article)).(*Article)
}

func NewArticleRepo() (*Repo, error) {
	return CreateRepo(new(Article))
}
