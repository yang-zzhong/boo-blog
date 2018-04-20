package controller

import (
	"boo-blog/model"
	"database/sql"
	"fmt"
	"log"
	. "net/http"
	"reflect"
)

type User struct {
}

/**
 * @request
 * {
 *	 email_addr: string,
 *   name: string,
 *   phone_number: (string)
 * }
 */
func (user *User) CreateUser(w ResponseWriter, req *Request) {
	mUser := model.NewUser()
	mUser.EmailAddr = sql.NullString{req.FormValue("email_addr"), true}
	mUser.PhoneNumber = sql.NullString{req.FormValue("phone_number"), true}
	mUser.Name = req.FormValue("name")
	mUser.NickName = sql.NullString{mUser.Name, true}
	repo, err := model.NewUserRepo()
	if err != nil {
		fmt.Println(err)
		log.Print(err)
	}
	if err = repo.Create(mUser); err != nil {
		fmt.Println(err)
		log.Print(err)
	}
}

func (user *User) RenderUsers(w ResponseWriter, req *Request) {
	fmt.Println(reflect.TypeOf(w).Name)
	repo, err := model.NewUserRepo()
	log.Print("开始执行了")
	if err != nil {
		log.Print(err)
		return
	}
	result, err := repo.Fetch()
	if err != nil {
		log.Print(err)
		return
	}
	log.Print("找到结果了")
	for _, item := range result {
		log.Print(item.(model.User))
	}
}
