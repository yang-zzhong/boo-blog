package main

import (
	. "boo-blog/model"
	"fmt"
	"github.com/dmulholland/args"
	"github.com/go-ini/ini"
	model "github.com/yang-zzhong/go-model"
)

func main() {
	parser := args.NewParser()
	parser.NewString("-c", "http/http.ini")
	if config, err := ini.Load(parser.GetString("-c")); err == nil {
		InitDriver(config.Section("database"))
	} else {
		panic(err)
	}
	repos := []*model.Repo{
		NewBlog().Repo(),
		NewCate().Repo(),
		NewImage().Repo(),
		NewTheme().Repo(),
		NewTag().Repo(),
		NewUserImage().Repo(),
		NewUser().Repo(),
		NewComment().Repo(),
	}
	for _, repo := range repos {
		err := repo.CreateRepo()
		if err != nil {
			fmt.Println(err)
		}
		s := "ALTER TABLE " + repo.QuotedTableName() + " CONVERT TO CHARACTER SET utf8"
		if _, err := DB.Exec(s); err != nil {
			fmt.Println(err)
		}
	}
}
