package model

import (
	"encoding/json"
	"github.com/google/uuid"
	model "github.com/yang-zzhong/go-model"
	"reflect"
	"time"
)

const (
	MSG_CREATE_BLOG = 1
	MSG_UPDATE_BLOG = 2
	MSG_DELETE_BLOG = 3
)

const (
	MSG_STATUS_READ   = 1
	MSG_STATUS_UNREAD = 2
	MSG_STATUS_DELETE = 3
)

type Message struct {
	Id      uint32                 `db:"id bigint pk"`
	From    uint32                 `db:"from bigint"`
	To      uint32                 `db:"to bigint"`
	Content string                 `db:"content text"`
	Type    int                    `db:"type int"`
	Status  int                    `db:"status int"`
	About   map[string]interface{} `db:"about jsonb nil"`
	At      time.Time              `db:"at timestamp"`
}

func (msg *Message) TableName() string {
	return "messages"
}

func (msg *Message) DBValue(colname string, value interface{}) interface{} {
	if colname == "about" {
		res, _ := json.Marshal(value)
		return string(res)
	}

	return value
}

func (msg *Message) Value(colname string, value interface{}) (result reflect.Value, catch bool) {
	if colname == "about" {
		catch = true
		val := value.(string)
		var res map[string]string
		json.Unmarshal([]byte(val), &res)
		result = reflect.ValueOf(res)
	}

	return
}

func NewMessage() *Message {
	msg := model.NewModel(new(Message)).(*Message)
	// msg.DeclareOne("from", new(User), model.Nexus{
	// 	"id": "from",
	// })
	// msg.DeclareOne("to", new(User), model.Nexus{
	// 	"id": "to",
	// })
	return msg
}

func (msg *Message) Instance() *Message {
	msg.Id = uuid.New().ID()
	return msg
}
