package main

import (
	"boo-blog/config"
	. "boo-blog/model"
	"fmt"
	"github.com/dmulholland/args"
	repo "github.com/yang-zzhong/go-model"
)

func main() {
	parser := args.NewParser()
	parser.NewString("-c", "http/http.conf")
	config.InitConfig(parser.GetString("-c"))
	InitDriver()
	imageGroupRepo, err := NewImageGroupRepo()
	// imageRepo, err := NewImageRepo()
	// userRepo, err := NewUserRepo()
	if err != nil {
		panic(err)
	}
	repos := []*repo.Repo{imageGroupRepo}
	for _, repo := range repos {
		err := repo.CreateTable()
		if err != nil {
			fmt.Println(err)
		}
	}
}
