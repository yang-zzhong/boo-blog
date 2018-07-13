package main

import (
	blog "boo-blog"
	"log"
)

func main() {
	blogger := blog.NewBlogger("./http.ini")
	blogger.HandleQuitEvent()
	if err := blogger.StartHttp(); err != nil {
		log.Fatal(err)
	}
}
