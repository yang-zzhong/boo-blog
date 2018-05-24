package session

import (
	"boo-blog/model"
	"crypto/md5"
	"encoding/hex"
)

var users map[string]*model.User

func Save(user *model.User) string {
	id := Id(user.Id)
	if users == nil {
		users = make(map[string]*model.User)
	}
	users[id] = user

	return id
}

func Id(userId string) string {
	md5Sum := md5.Sum([]byte(userId))

	return hex.EncodeToString(md5Sum[:])
}

func User(id string) (user *model.User, ok bool) {
	user, ok = users[id]
	return
}

func Del(id string) {
	delete(users, id)
}
