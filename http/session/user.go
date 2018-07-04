package session

import (
	"boo-blog/model"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"github.com/go-redis/redis"
	"strconv"
)

func Save(user *model.User) string {
	id := Id(user.Id)
	if users == nil {
		users = make(map[string]*model.User)
	}
	users[id] = togob64(user)

	return id
}

func Id(userId uint32) string {
	md5Sum := md5.Sum([]byte(strconv.FormatUint(uint64(userId), 32)))

	return hex.EncodeToString(md5Sum[:])
}

func User(id string) (user *model.User, ok bool) {
	var str string
	redis := cache.NewRedisClient(0)
	if !ok {
		return
	}
	user = fromgob64(str)

	return
}

func Del(id string) {
	delete(users, id)
}

func togob64(user *model.User) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(*user)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func fromgob64(str string) *model.User {
	user := model.User{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&user)
	if err != nil {
		panic(err)
	}
	return &user
}
