package controller

import (
	. "boo-blog/log"
	"boo-blog/model"
	"fmt"
	. "net/http"
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
	mUser.EmailAddr = req.FormValue("email_addr")
	mUser.PhoneNumber = req.FormValue("phone_number")
	mUser.Name = req.FormValue("name")
	mUser.NickName = mUser.Name
	repo, err := model.NewUserRepo()
	if err != nil {
		fmt.Println(err)
		Logger().Print(err)
	}
	if err = repo.Create(mUser); err != nil {
		fmt.Println(err)
		Logger().Print(err)
	}
}

func (user *User) RenderUsers(w ResponseWriter, req *Request) {
	repo, err := model.NewUserRepo()
	Logger().Print("开始执行了")
	if err != nil {
		Logger().Print(err)
		return
	}
	result, err := repo.Fetch()
	if err != nil {
		Logger().Print(err)
		return
	}
	Logger().Print("找到结果了")
	for _, item := range result {
		Logger().Print(item.(model.User))
	}
}
