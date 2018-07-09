package session

import (
	"boo-blog/cache"
	"boo-blog/model"
	"crypto/md5"
	"encoding/hex"
	"log"
	"strconv"
)

func Save(user *model.User) (id string, err error) {
	redis := cache.NewRedisClient(0)
	id = Id(user.Id)
	err = redis.Set(id, user.Id, 0).Err()
	return
}

func Id(userId uint32) string {
	md5Sum := md5.Sum([]byte(strconv.FormatUint(uint64(userId), 32)))

	return hex.EncodeToString(md5Sum[:])
}

func User(id string) (user *model.User, ok bool) {
	redis := cache.NewRedisClient(0)
	if userId, err := redis.Get(id).Result(); err != nil {
		ok = false
	} else {
		if m, ok, err := model.NewUser().Repo().Find(userId); err != nil {
			ok = false
		} else if ok {
			log.Print(m)
			user = m.(*model.User)
		}
	}

	return
}

func Del(id string) {
	redis := cache.NewRedisClient(0)
	redis.Del(id)
}
