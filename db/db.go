package db

import (
	"sync"

	"github.com/go-pg/pg"
)

var (
	instance *pg.DB
	once     sync.Once
)

func GetDbSession() *pg.DB {
	once.Do(func() {
		db := pg.Connect(&pg.Options{
			User:     "db_user",
			Password: "db_pwd",
			Database: "im_docker",
			Addr:     "localhost:6100",
		})
		instance = db
	})
	return instance
}

type Chat struct {
	Id int64
}

type ChatMember struct {
	tableName struct{} `sql:"chat_member,alias:genre"`
	Id        int64
	ChatId    int64 `sql:"chat_id"`
	UserId    int64 `sql:"user_id"`
}
