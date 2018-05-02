package main

import (
	. "boo-blog/model"
	"fmt"
	"github.com/dmulholland/args"
	"github.com/go-ini/ini"
	repo "github.com/yang-zzhong/go-model"
)

func main() {
	parser := args.NewParser()
	parser.NewString("-c", "http/http.ini")
	if config, err := ini.Load(parser.GetString("-c")); err == nil {
		InitDriver(config.Section("database"))
	} else {
		panic(err)
	}
	cateRepo, err := NewCategoryRepo()
	if err != nil {
		panic(err)
	}
	imageRepo, err := NewImageRepo()
	if err != nil {
		panic(err)
	}
	userImageRepo, err := NewUserImageRepo()
	if err != nil {
		panic(err)
	}
	userRepo, err := NewUserRepo()
	if err != nil {
		panic(err)
	}
	tagRepo, err := NewTagRepo()
	repos := []*repo.Repo{userRepo, cateRepo, imageRepo, userImageRepo, tagRepo}
	for _, repo := range repos {
		err := repo.CreateTable()
		if err != nil {
			fmt.Println(err)
		}
	}
}
