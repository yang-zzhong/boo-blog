package main

import blog "boo-blog"

func main() {
	blogger := blog.NewBlogger("./http.ini")
	blogger.StartHttp()
}
