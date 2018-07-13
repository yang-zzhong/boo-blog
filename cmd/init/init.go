package main

import (
	blog "boo-blog"
	"log"
)

func main() {
	blogger := blog.NewBlogger("./http.ini")
	if err := blogger.CreateTable(); err != nil {
		log.Fatal(err)
	}
}
