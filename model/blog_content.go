package model

import (
	"github.com/google/uuid"
	model "github.com/yang-zzhong/go-model"
	"golang.org/x/net/html"
	"strings"
)

type BlogContent struct {
	Id      uint32 `db:"id bigint pk"`
	Content string `db:"content text"`
	*model.Base
}

func (blog *BlogContent) TableName() string {
	return "blog_contents"
}

func NewBlogContent() *BlogContent {
	content := model.NewModel(new(BlogContent)).(*BlogContent)
	content.DeclareOne("blog", new(Blog), model.Nexus{
		"content_id": "id",
	})

	return content
}

func (blog *BlogContent) PreviewImageUrl() string {
	reader := strings.NewReader(blog.Content)
	node, _ := html.Parse(reader)
	nodes := 0
	url := ""
	find(node, func(d *html.Node) bool {
		if nodes > 300 {
			return true
		}
		nodes++
		if d.Type == html.ElementNode && d.Data == "img" {
			for _, attr := range d.Attr {
				if attr.Key == "src" {
					url = attr.Val
					return true
				}
			}
		}
		return false
	})

	return url
}

func (blog *BlogContent) Preview(limit int) string {
	reader := strings.NewReader(blog.Content)
	node, _ := html.Parse(reader)
	nodes := 0
	preview := []rune{}
	find(node, func(d *html.Node) bool {
		if nodes > 300 {
			return true
		}
		if d.Type == html.ElementNode && (d.Data == "style" || d.Data == "code" || d.Data == "head") {
			return false
		}
		nodes++
		if d.Type == html.TextNode {
			preview = append(preview, []rune(d.Data)...)
		}
		if len(preview) > limit {
			return true
		}
		return false
	})
	if len(preview) > limit {
		preview = preview[0:limit]
	}

	return string(preview)
}

func (blog *BlogContent) Instance(content string) *BlogContent {
	blog.Id = uuid.New().ID()
	blog.Content = content

	return blog
}

type callback func(n *html.Node) bool

func find(n *html.Node, call callback) bool {
	if call(n) {
		return true
	}
	var found bool
	for c := n.FirstChild; c != nil; c = c.FirstChild {
		found = find(c, call)
	}
	if found {
		return found
	}
	for c := n.NextSibling; c != nil; c = c.NextSibling {
		found = find(c, call)
	}
	return found
}
