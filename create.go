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
	// imageGroupRepo, err := NewCategoryRepo()
	// imageRepo, err := NewImageRepo()
	userRepo, err := NewUserRepo()
	if err != nil {
		panic(err)
	}
	repos := []*repo.Repo{userRepo}
	for _, repo := range repos {
		err := repo.CreateTable()
		if err != nil {
			fmt.Println(err)
		}
	}
}
