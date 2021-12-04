package service

import (
	"fmt"
	"time"

	"github.com/CodFrm/qqbot-official/internal/db"
	"github.com/CodFrm/qqbot-official/internal/utils/errs"
	"github.com/CodFrm/qqbot/utils"
	"github.com/tencent-connect/botgo/openapi"
)

type Punish interface {
	PunishUser(guild, operator, userid string) (int64, error)
	PunishLevel(guild, userid string) (int64, error)
}

type punish struct {
	openapi openapi.OpenAPI
}

func NewPunish() Punish {
	return &punish{}
}

func (p *punish) PunishUser(guild, operator, userid string) (int64, error) {
	num, err := db.Incr("guild:punish:user:num:"+time.Now().Format("2006:01:02")+":"+operator, 1, 3600*24)
	if err != nil {
		return 0, err
	}
	if num > 40 {
		return 0, errs.NewReplyError("今天已经处理够多人了")
	}
	return db.Incr(fmt.Sprintf("guild:punish:%v:%v", guild, userid), 1, 604800)
}

func (p *punish) PunishLevel(guild, userid string) (int64, error) {
	val, err := db.Get(fmt.Sprintf("guild:punish:%v:%v", guild, userid))
	if err != nil {
		return 0, err
	}
	return utils.StringToInt64(string(val)), nil
}
