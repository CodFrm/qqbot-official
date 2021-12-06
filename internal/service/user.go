package service

import (
	"github.com/CodFrm/qqbot-official/internal/db"
	"github.com/CodFrm/qqbot/utils"
)

type User interface {
	AddIntegral(guild, user string, num int64) (int64, error)
	Integral(guild, user string) (int64, error)
}

type user struct {
}

func NewUser() User {
	return &user{}
}

func (u *user) AddIntegral(guild, user string, num int64) (int64, error) {
	return db.Incr("user:integral:"+guild+":"+user, num, 0)
}

func (u *user) Integral(guild, user string) (int64, error) {
	val, err := db.Get("user:integral:" + guild + ":" + user)
	if err != nil {
		return 0, err
	}
	return utils.StringToInt64(string(val)), nil
}
