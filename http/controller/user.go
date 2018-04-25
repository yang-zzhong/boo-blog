package controller

import (
	"boo-blog/model"
	"database/sql"
	httprouter "github.com/yang-zzhong/go-httprouter"
)

type User struct{ *Controller }

func NewUser(responseWriter interface{}) *User {
	return &User{NewController(responseWriter)}
}

/**
 * @request
 * {
 *	 email_addr: string,
 *   name: string,
 *   phone_number: (string)
 * }
 */
func (user *User) CreateUser(req *httprouter.Request) {
	mUser := model.NewUser()
	mUser.EmailAddr = sql.NullString{req.FormValue("email_addr"), true}
	mUser.PhoneNumber = sql.NullString{req.FormValue("phone_number"), true}
	mUser.Name = req.FormValue("name")
	mUser.NickName = sql.NullString{mUser.Name, true}
	repo, err := model.NewUserRepo()
	if err != nil {
		user.InternalError(err)
		return
	}
	if err = repo.Create(mUser); err != nil {
		user.InternalError(err)
		return
	}
}

func (user *User) RenderUsers(req *httprouter.Request) {
}
